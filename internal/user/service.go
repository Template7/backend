// Package user manages accounts: creation, activation, profile and deletion.
package user

import (
	"context"

	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
)

type Service interface {
	// CreateUser creates an initialized account with a default wallet and returns
	// its id and a one-time activation code.
	CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (userId string, activationCode string, err error)
	// ActivateUser activates the account if code is its current activation code.
	ActivateUser(ctx context.Context, userId string, code string) error
	DeleteUser(ctx context.Context, userId string) error

	GetInfo(ctx context.Context, userId string) (*userV1.UserInfoResponse, error)
	UpdateInfo(ctx context.Context, userId string, info models.UserInfo) error
	GetStatus(ctx context.Context, userId string) (authV1.AccountStatus, error)
}

// defaultCurrencies get a zero balance in every new wallet.
var defaultCurrencies = []string{
	walletV1.Currency_usd.String(),
	walletV1.Currency_ntd.String(),
	walletV1.Currency_cny.String(),
	walletV1.Currency_jpy.String(),
}

type service struct {
	auth  auth.Auth
	db    db.Client
	cache cache.Interface
	log   *logger.Logger
}

func New(authSvc auth.Auth, db db.Client, cache cache.Interface, log *logger.Logger) Service {
	return &service{
		auth:  authSvc,
		db:    db,
		cache: cache,
		log:   log.WithService("user"),
	}
}
