package wallet

import (
	"context"
	"errors"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
	"net/http"
)

type Service struct {
	db  db.Client
	log *logger.Logger
}

func New(db db.Client, log *logger.Logger) *Service {
	s := &Service{
		db:  db,
		log: log.WithService("wallet"),
	}
	s.log.Info("wallet service initialized")
	return s
}

func (s *Service) GetWallet(ctx context.Context, walletId string) (*v1.Wallet, error) {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet")

	wbs, err := s.db.GetWalletBalances(ctx, walletId)
	if err != nil {
		log.WithError(err).Error("fail to get wallet balances")
		return nil, t7Error.DbOperationFail.WithDetail(err.Error())
	}
	if len(wbs) == 0 {
		log.Info("no balances")
		return nil, t7Error.WalletNotFound.WithStatus(http.StatusNotFound)
	}

	w := v1.Wallet{
		Id: wbs[0].WalletId,
	}
	for _, wb := range wbs {
		w.Balances = append(w.Balances, &v1.Balance{
			Currency: v1.Currency(v1.Currency_value[wb.Currency]),
			Amount:   wb.Amount.String(),
		})
	}

	return &w, nil
}

func (s *Service) Deposit(ctx context.Context, walletId string, currency v1.Currency, amount uint32, note string) error {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("deposit")

	m, err := newMoney(currency, amount)
	if err != nil {
		log.WithError(err).Warn("invalid money")
		return err
	}
	if err := s.db.Deposit(ctx, walletId, m, note); err != nil {
		log.WithError(err).Error("fail to deposit")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("deposit success")
	return nil
}

func (s *Service) Withdraw(ctx context.Context, walletId string, currency v1.Currency, amount uint32, note string) error {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("withdraw")

	m, err := newMoney(currency, amount)
	if err != nil {
		log.WithError(err).Warn("invalid money")
		return err
	}
	if err := s.db.Withdraw(ctx, walletId, m, note); err != nil {
		if errors.Is(err, db.ErrInsufficientBalance) {
			return t7Error.InsufficientBalance
		}
		log.WithError(err).Error("fail to withdraw")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("withdraw success")
	return nil
}

// Transfer moves money from one of userId's wallets to another wallet.
// The caller must own fromWalletId.
func (s *Service) Transfer(ctx context.Context, userId string, fromWalletId string, toWalletId string, currency v1.Currency, amount uint32, note string) error {
	log := s.log.WithContext(ctx).With("userId", userId).With("fromWalletId", fromWalletId).With("toWalletId", toWalletId)
	log.Debug("transfer")

	if fromWalletId == toWalletId {
		log.Warn("transfer to the same wallet")
		return t7Error.InvalidBody.WithDetail("fromWalletId and toWalletId must differ")
	}
	if !s.isUserWallet(ctx, userId, fromWalletId) {
		log.Warn("user does not own the source wallet")
		return t7Error.WalletNotFound.WithStatus(http.StatusForbidden)
	}

	m, err := newMoney(currency, amount)
	if err != nil {
		log.WithError(err).Warn("invalid money")
		return err
	}
	if err := s.db.Transfer(ctx, fromWalletId, toWalletId, m, note); err != nil {
		if errors.Is(err, db.ErrInsufficientBalance) {
			return t7Error.InsufficientBalance
		}
		log.WithError(err).Error("fail to transfer")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("transfer success")
	return nil
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

func (s *Service) isUserWallet(ctx context.Context, userId string, walletId string) bool {
	if userId == "" || walletId == "" {
		return false
	}
	for _, uw := range s.db.GetUserWallets(ctx, userId) {
		if uw.WalletId == walletId {
			return true
		}
	}
	return false
}
