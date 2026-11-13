// Package wallet holds the money rules: ownership, amounts, balances and history.
// Every call takes the caller's user id and only works on wallets that user owns
// (transfers may send to any wallet).
package wallet

import (
	"context"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	"github.com/Template7/common/t7Id"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
)

type Service interface {
	GetUserWallets(ctx context.Context, userId string) ([]*v1.Wallet, error)
	GetWallet(ctx context.Context, userId string, walletId string) (*v1.Wallet, error)

	Deposit(ctx context.Context, userId string, walletId string, currency v1.Currency, amount uint32, note string) error
	Withdraw(ctx context.Context, userId string, walletId string, currency v1.Currency, amount uint32, note string) error
	// Transfer moves money from one of the caller's wallets to any other wallet.
	Transfer(ctx context.Context, userId string, fromWalletId string, toWalletId string, currency v1.Currency, amount uint32, note string) error

	GetBalanceHistory(ctx context.Context, userId string, walletId string) ([]*v1.BalanceHistory, error)
	GetBalanceHistoryByCurrency(ctx context.Context, userId string, walletId string, currency v1.Currency) ([]*v1.CurrencyBalanceHistory, error)
}

type service struct {
	db  db.Client
	log *logger.Logger
}

func New(db db.Client, log *logger.Logger) Service {
	return &service{
		db:  db,
		log: log.WithService("wallet"),
	}
}

func (s *service) GetUserWallets(ctx context.Context, userId string) ([]*v1.Wallet, error) {
	rows, err := s.db.GetUserWallets(ctx, userId)
	if err != nil {
		return nil, dbError(err)
	}

	// rows are ordered by wallet id, so balances of one wallet are adjacent
	var data []*v1.Wallet
	for _, r := range rows {
		if len(data) == 0 || data[len(data)-1].Id != r.WalletId {
			data = append(data, &v1.Wallet{Id: r.WalletId})
		}
		w := data[len(data)-1]
		w.Balances = append(w.Balances, toBalance(r.Currency, r.Amount))
	}
	return data, nil
}

func (s *service) GetWallet(ctx context.Context, userId string, walletId string) (*v1.Wallet, error) {
	if err := s.checkOwner(ctx, userId, walletId); err != nil {
		return nil, err
	}
	rows, err := s.db.GetWalletBalances(ctx, walletId)
	if err != nil {
		return nil, dbError(err)
	}

	w := &v1.Wallet{Id: walletId}
	for _, r := range rows {
		w.Balances = append(w.Balances, toBalance(r.Currency, r.Amount))
	}
	return w, nil
}

func (s *service) Deposit(ctx context.Context, userId string, walletId string, currency v1.Currency, amount uint32, note string) error {
	m, err := newMoney(currency, amount)
	if err != nil {
		return err
	}
	if err := s.checkOwner(ctx, userId, walletId); err != nil {
		return err
	}

	return s.inTx(ctx, func(tx db.Client) error {
		before, err := lockOne(ctx, tx, walletId, m.Currency)
		if err != nil {
			return err
		}
		if err := tx.AddBalance(ctx, walletId, m.Currency, m.Amount); err != nil {
			return err
		}
		return tx.CreateDepositHistory(ctx, models.DepositHistory{
			Id:            newId(),
			WalletId:      walletId,
			Currency:      m.Currency,
			Amount:        m.Amount,
			BalanceBefore: before,
			BalanceAfter:  before.Add(m.Amount),
			Note:          note,
		})
	})
}

func (s *service) Withdraw(ctx context.Context, userId string, walletId string, currency v1.Currency, amount uint32, note string) error {
	m, err := newMoney(currency, amount)
	if err != nil {
		return err
	}
	if err := s.checkOwner(ctx, userId, walletId); err != nil {
		return err
	}

	return s.inTx(ctx, func(tx db.Client) error {
		before, err := lockOne(ctx, tx, walletId, m.Currency)
		if err != nil {
			return err
		}
		if before.LessThan(m.Amount) {
			return t7Error.InsufficientBalance
		}
		if err := tx.AddBalance(ctx, walletId, m.Currency, m.Amount.Neg()); err != nil {
			return err
		}
		return tx.CreateWithdrawHistory(ctx, models.WithdrawHistory{
			Id:            newId(),
			WalletId:      walletId,
			Currency:      m.Currency,
			Amount:        m.Amount,
			BalanceBefore: before,
			BalanceAfter:  before.Sub(m.Amount),
			Note:          note,
		})
	})
}

func (s *service) Transfer(ctx context.Context, userId string, fromWalletId string, toWalletId string, currency v1.Currency, amount uint32, note string) error {
	if fromWalletId == toWalletId {
		return t7Error.InvalidBody.WithDetail("fromWalletId and toWalletId must differ")
	}
	m, err := newMoney(currency, amount)
	if err != nil {
		return err
	}
	if err := s.checkOwner(ctx, userId, fromWalletId); err != nil {
		return err
	}

	return s.inTx(ctx, func(tx db.Client) error {
		bls, err := tx.LockBalances(ctx, []string{fromWalletId, toWalletId}, m.Currency)
		if err != nil {
			return err
		}
		senderBefore, ok := bls[fromWalletId]
		if !ok {
			return t7Error.WalletNotFound
		}
		receiverBefore, ok := bls[toWalletId]
		if !ok {
			return t7Error.WalletNotFound.WithDetail("receiver wallet")
		}
		if senderBefore.LessThan(m.Amount) {
			return t7Error.InsufficientBalance
		}

		if err := tx.AddBalance(ctx, fromWalletId, m.Currency, m.Amount.Neg()); err != nil {
			return err
		}
		if err := tx.AddBalance(ctx, toWalletId, m.Currency, m.Amount); err != nil {
			return err
		}
		return tx.CreateTransferHistory(ctx, models.TransferHistory{
			Id:                    newId(),
			FromWalletId:          fromWalletId,
			ToWalletId:            toWalletId,
			Currency:              m.Currency,
			Amount:                m.Amount,
			SenderBalanceBefore:   senderBefore,
			SenderBalanceAfter:    senderBefore.Sub(m.Amount),
			ReceiverBalanceBefore: receiverBefore,
			ReceiverBalanceAfter:  receiverBefore.Add(m.Amount),
			Note:                  note,
		})
	})
}

// checkOwner returns WalletNotFound unless walletId exists and belongs to userId,
// so a wallet someone else owns looks the same as one that doesn't exist.
func (s *service) checkOwner(ctx context.Context, userId string, walletId string) error {
	if userId == "" || walletId == "" {
		return t7Error.WalletNotFound
	}
	w, err := s.db.GetWallet(ctx, walletId)
	if db.IsNotFound(err) {
		return t7Error.WalletNotFound
	}
	if err != nil {
		return dbError(err)
	}
	if w.UserId != userId {
		return t7Error.WalletNotFound
	}
	return nil
}

// inTx runs fn in a transaction and turns storage errors into t7Errors.
func (s *service) inTx(ctx context.Context, fn func(tx db.Client) error) error {
	if err := s.db.Transaction(ctx, fn); err != nil {
		return dbError(err)
	}
	return nil
}

func lockOne(ctx context.Context, tx db.Client, walletId string, currency string) (decimal.Decimal, error) {
	bls, err := tx.LockBalances(ctx, []string{walletId}, currency)
	if err != nil {
		return decimal.Zero, err
	}
	b, ok := bls[walletId]
	if !ok {
		return decimal.Zero, t7Error.WalletNotFound.WithDetail("no " + currency + " balance")
	}
	return b, nil
}

// newMoney validates the amount and converts it without overflow
// (uint32 -> int32 would turn amounts above math.MaxInt32 negative).
func newMoney(currency v1.Currency, amount uint32) (models.Money, error) {
	if amount == 0 {
		return models.Money{}, t7Error.InvalidBody.WithDetail("amount must be greater than 0")
	}
	name, ok := v1.Currency_name[int32(currency)]
	if !ok {
		return models.Money{}, t7Error.InvalidBody.WithDetail("unsupported currency")
	}
	return models.Money{
		Currency: name,
		Amount:   decimal.NewFromInt(int64(amount)),
	}, nil
}

// dbError keeps t7Errors (business errors returned inside a transaction) and
// wraps anything else as a DB failure.
func dbError(err error) error {
	if _, ok := t7Error.ToT7Error(err); ok {
		return err
	}
	if db.IsNotFound(err) {
		return t7Error.WalletNotFound
	}
	return t7Error.DbOperationFail.WithDetail(err.Error())
}

func newId() int64 {
	return t7Id.New().Generate().Int64()
}

func toBalance(currency string, amount decimal.Decimal) *v1.Balance {
	return &v1.Balance{
		Currency: v1.Currency(v1.Currency_value[currency]),
		Amount:   amount.String(),
	}
}
