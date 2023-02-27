package db

import (
	"context"
	"github.com/Template7/common/structs"
	"github.com/Template7/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (c *impl) GetWallet(userId string) (data structs.WalletData, err error) {
	var wallet structs.Wallet
	if err = c.mysql.db.Model(&structs.Wallet{}).Where("userId = ?", userId).Take(&wallet).Error; err != nil {
		log.Error("fail to get wallet for user: ", userId, ". ", err.Error())
		return
	}

	data.Id = wallet.Id
	data.UserId = wallet.UserId
	err = c.mysql.db.Model(&structs.Balance{}).Where("walletId = ?", wallet.Id).Find(&data.Balance).Error
	return
}

func (c *impl) Deposit(data DepositData) (err error) {
	data.CreatedAt = time.Now()
	if _, err = c.mongo.depositHistory.InsertOne(context.Background(), data); err != nil {
		log.Error("fail to insert document: ", err.Error())
		return
	}

	err = c.mysql.db.Model(&structs.Balance{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "walletId"}, {Name: "currency"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", util.ToPico(data.Money))}),
		}).Create(&structs.Balance{WalletId: data.WalletId, Money: structs.Money{Currency: data.Money.Currency, Amount: util.ToPico(data.Money), Unit: structs.UnitPico}}).Error
	return
}

func (c *impl) Withdraw(data WithdrawData) (err error) {
	data.CreatedAt = time.Now()
	if _, err = c.mongo.withdrawHistory.InsertOne(context.Background(), data); err != nil {
		log.Error("fail to insert document: ", err.Error())
		return
	}

	amount := util.ToPico(data.Money)
	err = c.mysql.db.Model(&structs.Balance{}).
		Take(&structs.Balance{}, "walletId = ? AND currency = ? AND amount >= ?", data.WalletId, data.Money.Currency, amount).
		Update("amount", gorm.Expr("amount - ?", amount)).Error
	return
}

func (c *impl) Transfer(data TransactionData) (err error) {
	data.CreatedAt = time.Now()
	// insert to mongodb
	if _, err := c.mongo.transactionHistory.InsertOne(context.Background(), data); err != nil {
		log.Error("fail to insert transaction data: ", err.Error())
		return err
	}

	return c.mysql.db.Model(&structs.Balance{}).Transaction(func(tx *gorm.DB) error {
		// reduce from the source wallet
		var blc structs.Balance
		amount := util.ToPico(data.Money)

		if err := tx.Take(&blc, "walletId = ? AND currency = ? AND amount >= ?", data.FromWalletId, data.Currency, amount).
			Update("amount", gorm.Expr("amount - ?", amount)).Error; err != nil {
			log.Error("fail to take balance: ", err.Error())
			return err
		}

		// increment to the target wallet
		if err := tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "walletId"}, {Name: "currency"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", amount)}),
			}).Create(&structs.Balance{WalletId: data.ToWalletId, Money: structs.Money{Currency: data.Currency, Amount: amount, Unit: structs.UnitPico}}).Error; err != nil {
			log.Error("fail to add money: ", err.Error())
			return err
		}

		return nil
	})
}

// TODO: add paging and some query filter
func (c *impl) GetTransactions(userId string) (data []TransactionData, err error) {
	var wallet structs.Wallet
	if err = c.mysql.db.Model(&structs.Wallet{}).Where("userId = ?", userId).Take(&wallet).Error; err != nil {
		log.Error("fail to get wallet for user: ", userId, ". ", err.Error())
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
