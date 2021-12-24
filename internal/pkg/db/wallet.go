package db

import (
	"context"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/structs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (c client) GetWallet(userId string) (data structs.WalletData, err error) {
	var wallet structs.Wallet
	c.mysql.db.Model(&structs.Wallet{}).Where("userId = ?", userId).Take(&wallet)
	if wallet.Id == "" {
		log.Warn("no related wallet for the user: ", userId)
		err = t7Error.WalletNotFound
		return
	}

	data.Id = wallet.Id
	data.UserId = wallet.UserId
	err = c.mysql.db.Model(&structs.Balance{}).Where("walletId = ?", wallet.Id).Find(&data.Balance).Error
	return
}

func (c client) Deposit(walletId string, money structs.Money) (err error) {
	err = c.mysql.db.Model(&structs.Balance{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "walletId"}, {Name: "currency"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", money.Amount)}),
		}).Create(&structs.Balance{WalletId: walletId, Money: structs.Money{Currency: money.Currency, Amount: money.Amount, Unit: structs.UnitPico}}).Error
	return
}

func (c client) Withdraw(walletId string, money structs.Money) (err error) {
	err = c.mysql.db.
		Take(&structs.Balance{}, "walletId = ? AND currency = ? AND amount >= ?", walletId, money.Currency, money.Amount).
		Update("amount", gorm.Expr("amount - ?", money.Amount)).Error
	return
}

func (c client) Transfer(data TransactionData) (err error) {
	data.CreatedAt = time.Now()
	// insert to mongodb
	if _, err := c.mongo.transactionHistory.InsertOne(context.Background(), data); err != nil {
		log.Error("fail to insert transaction data: ", err.Error())
		return err
	}

	return c.mysql.db.Transaction(func(tx *gorm.DB) error {
		// reduce from the source wallet
		tx.Model(&structs.Balance{}).
			Where("walletId = ? AND currency = ? AND amount >= ?", data.FromWalletId, data.Currency, data.Amount).
			Update("amount", gorm.Expr("amount - ?", data.Amount))

		// increment to the target wallet
		tx.Model(&structs.Balance{}).Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "walletId"}, {Name: "currency"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", data.Amount)}),
			}).Create(&structs.Balance{WalletId: data.ToWalletId, Money: structs.Money{Currency: data.Currency, Amount: data.Amount, Unit: structs.UnitPico}})

		return nil
	})
}

// TODO: add paging and some query filter
func (c client) GetTransactions(userId string) (data []TransactionData, err error) {
	var wallet structs.Wallet
	c.mysql.db.Model(&structs.Wallet{}).Where("userId = ?", userId).Take(&wallet)
	if wallet.Id == "" {
		log.Warn("no related wallet for the user: ", userId)
		err = t7Error.WalletNotFound
		return
	}
	filter := bson.M{
		"$or": []bson.M{
			{
				"from_wallet_id": wallet.Id,
			}, {
				"to_wallet_id": wallet.Id,
			},
		},
	}
	cursor, err := c.mongo.transactionHistory.Find(context.Background(), filter)
	if err != nil {
		log.Error("fail to find document: ", err.Error())
		return nil, err
	}
	if err = cursor.All(context.Background(), &data); err != nil {
		log.Error("fail tp decode data: ", err.Error())
	}
	return
}
