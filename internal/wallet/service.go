package wallet

import (
	"context"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
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
