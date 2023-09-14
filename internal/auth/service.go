package auth

import (
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
)

const (
	jwtSign = "45519f46c06c8340a34f9a32982860c1a8d6bb57eaeb338b7f0119062b8a3b67"
)

var (
	once     sync.Once
	instance *service
)

type UserTokenClaims struct {
	jwt.StandardClaims
	v1.TokenClaims
}

type service struct {
	core *casbin.Enforcer
	db   db.Client
	log  *logger.Logger
}

func New() Auth {
	once.Do(func() {
		log := logger.New().WithService("auth")

		cfg := config.New()
		adapter, err := gormadapter.NewAdapter("mysql", cfg.Sql.ConnectionString, true)
		if err != nil {
			log.WithError(err).Error("fail to new mysql adapter")
			panic(err)
		}
		e, err := casbin.NewEnforcer("./config/rbac_model.conf", adapter)
		if err != nil {
			log.WithError(err).Error("fail to new enforcer")
			panic(err)
		}
		err = e.LoadPolicy()
		if err != nil {
			log.WithError(err).Error("fail to load policy")
			panic(err)
		}
		e.AddFunction("checkAdmin", func(args ...interface{}) (interface{}, error) {
			username := args[0].(string)
			return e.HasRoleForUser(username, "admin")
		})

		instance = &service{
			core: e,
			db:   db.New(),
			log:  logger.New().WithService("auth"),
		}
		instance.log.Debug("auth service initialized")
	})

	return instance
}

// TODO: reset db policies
func (s *service) loadDefaultPolicies() {
	pPolicy := [][]string{
		{v1.Role_User.String(), "/api/v1/users/:userId/info", http.MethodGet},
		{v1.Role_User.String(), "/api/v1/users/:userId/info", http.MethodPut},
		{v1.Role_User.String(), "/api/v1/wallets/:walletId", http.MethodGet},
		{v1.Role_User.String(), "/api/v1/wallets/:walletId/deposit", http.MethodPost},
		{v1.Role_User.String(), "/api/v1/wallets/:walletId/withdraw", http.MethodPost},
		{v1.Role_User.String(), "/api/v1/transaction", http.MethodPost},
	}

	ok, err := s.core.AddPolicies(pPolicy)
	if err != nil {
		s.log.WithError(err).Error("fail to add policies")
	}
	if !ok {
		s.log.With("policy", pPolicy).Warn("no policies added")
	}
}

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
