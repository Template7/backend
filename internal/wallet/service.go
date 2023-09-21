package wallet

import (
	"context"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
	"sync"
)

var (
	once     sync.Once
	instance *Service
)

type Service struct {
	db  db.Client
	log *logger.Logger
}

func New() *Service {
	once.Do(func() {
		log := logger.New().WithService("wallet")
		instance = &Service{
			db:  db.New(),
			log: log,
		}
		log.Info("wallet service initialized")
	})
	return instance
}

func (s *Service) GetWallet(ctx context.Context, walletId string) (v1.Wallet, error) {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet")

	wallet, err := s.db.GetWallet(ctx, walletId)
	if err != nil {
		log.WithError(err).Error("fail to get wallet")
		return v1.Wallet{}, t7Error.DbOperationFail.WithDetail(err.Error())
	}

	return *wallet.ToProto(ctx), nil
}

func (s *Service) Deposit(ctx context.Context, walletId string, currency v1.Currency, amount uint32) error {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("deposit")

	m := entity.Money{
		Currency: v1.Currency_name[int32(currency)],
		Amount:   decimal.NewFromInt32(int32(amount)),
	}
	if err := s.db.Deposit(ctx, walletId, m); err != nil {
		log.WithError(err).Error("fail to deposit")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("deposit success")
	return nil
}

func (s *Service) Withdraw(ctx context.Context, walletId string, currency v1.Currency, amount uint32) error {
	log := s.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("withdraw")

	m := entity.Money{
		Currency: v1.Currency_name[int32(currency)],
		Amount:   decimal.NewFromInt32(int32(amount)),
	}
	if err := s.db.Withdraw(ctx, walletId, m); err != nil {
		log.WithError(err).Error("fail to withdraw")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("withdraw success")
	return nil
}

func (s *Service) Transfer(ctx context.Context, fromWalletId string, toWalletId string, currency v1.Currency, amount uint32) error {
	log := s.log.WithContext(ctx).With("fromWalletId", fromWalletId).With("toWalletId", toWalletId)
	log.Debug("transfer")

	m := entity.Money{
		Currency: v1.Currency_name[int32(currency)],
		Amount:   decimal.NewFromInt32(int32(amount)),
	}
	if err := s.db.Transfer(ctx, fromWalletId, toWalletId, m); err != nil {
		log.WithError(err).Error("fail to transfer")
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}

	log.Debug("transfer success")
	return nil
}
