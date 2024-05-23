package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func init() {
	viper.AddConfigPath("../../config")
}

type testDbClient struct {
	core *gorm.DB
}

// user methods
func (c *testDbClient) CreateUser(ctx context.Context, data entity.User) error {
	return c.core.Create(&data).Error
}

func (c *testDbClient) GetUser(ctx context.Context, username string) (entity.User, error) {
	return entity.User{}, errors.New("not implemented")
}

func (c *testDbClient) GetUserById(ctx context.Context, userId string) (entity.User, error) {
	return entity.User{}, errors.New("not implemented")
}

func (c *testDbClient) UpdateUserInfo(ctx context.Context, userId string, info entity.UserInfo) error {
	return errors.New("not implemented")
}

func (c *testDbClient) GetUserWallets(ctx context.Context, userId string) []entity.UserWalletBalance {
	return nil
}

func (c *testDbClient) DeleteUser(ctx context.Context, userId string) error {
	return errors.New("not implemented")
}

// wallet methods
func (c *testDbClient) GetWalletBalances(ctx context.Context, walletId string) ([]entity.WalletBalance, error) {
	return nil, errors.New("not implemented")
}

func (c *testDbClient) Deposit(ctx context.Context, walletId string, money entity.Money, note string) error {
	return errors.New("not implemented")
}

func (c *testDbClient) Withdraw(ctx context.Context, walletId string, money entity.Money, note string) error {
	return errors.New("not implemented")
}

func (c *testDbClient) Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money, note string) error {
	return errors.New("not implemented")
}

func (c *testDbClient) GetBalance(ctx context.Context, walletId string, currency string) (decimal.Decimal, error) {
	return decimal.Decimal{}, errors.New("not implemented")
}

func (c *testDbClient) getWalletsBalance(ctx context.Context, tx *gorm.DB, walletId []string, currency string) ([]entity.Balance, error) {
	return nil, errors.New("not implemented")
}

func (c *testDbClient) GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]entity.WalletBalanceHistory, error) {
	return nil, errors.New("not implemented")
}

type testCache struct {
	actCode map[string]string
}

func (c *testCache) SetUserActivationCode(ctx context.Context, userId string, code string) error {
	c.actCode[userId] = code
	return nil
}

func (c *testCache) GetUserActivationCode(ctx context.Context, userId string) (string, error) {
	code, ok := c.actCode[userId]
	if !ok {
		return "", fmt.Errorf("user not fount")
	}
	return code, nil
}

func newTestDbClient() db.Client {
	return &testDbClient{
		core: newTestDbCore(),
	}
}

func newTestCache() cache.Interface {
	return &testCache{
		actCode: map[string]string{},
	}
}

var (
	once       sync.Once
	testDbCore *gorm.DB
)

func newTestDbCore() *gorm.DB {
	once.Do(func() {
		tdc, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		if err := tdc.AutoMigrate(
			&entity.User{},
		); err != nil {
			panic(err)
		}
		testDbCore = tdc
	})
	return testDbCore
}

func newTestConfig() *config.Config {
	cfg := config.Config{}
	cfg.Auth.RbacModelPath = "../../config/rbac_model.conf"
	return &cfg
}

func Test_service_CreateUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), "traceId", "Test_service_CreateUser")

	req := userV1.CreateUserRequest{
		Username: "admin",
		Password: "password",
		Role:     authV1.Role_admin,
	}

	authAvc := New(newTestDbClient(), newTestDbCore(), newTestCache(), logger.New(), newTestConfig())

	userId, err := authAvc.CreateUser(ctx, &req)
	if err != nil {
		t.Error(err)
	}

	if userId == "" {
		t.Error("empty user id")
	}
}

func Test_service_GenActivationCode(t *testing.T) {
	ctx := context.WithValue(context.Background(), "traceId", "Test_service_GenActivationCode")
	authAvc := New(newTestDbClient(), newTestDbCore(), newTestCache(), logger.New(), newTestConfig())

	code, err := authAvc.GenActivationCode(ctx, "testUserId")
	if err != nil {
		t.Error(err)
	}

	if code == "" {
		t.Error("empty act code")
	}
}
