package auth

import (
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

const (
	minJwtSecretLen = 32
	bcryptCost      = bcrypt.DefaultCost
)

// dummyPasswordHash is compared against when the user doesn't exist,
// so a login for an unknown user takes as long as a wrong password.
var dummyPasswordHash, _ = hashedPassword("dummy-password-for-timing")

type UserTokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Status int    `json:"status"`
}

type service struct {
	core      *casbin.Enforcer
	jwtSecret []byte
	db        db.Client
	dbCore    *gorm.DB
	cache     cache.Interface
	log       *logger.Logger
}

func New(db db.Client, dbCore *gorm.DB, cache cache.Interface, log *logger.Logger, cfg *config.Config) Auth {
	if len(cfg.Auth.JwtSecret) < minJwtSecretLen {
		log.With("minLength", minJwtSecretLen).Panic("jwt secret missing or too short, set auth.jwtSecret / T7_AUTH_JWTSECRET")
	}

	adapter, err := gormadapter.NewAdapterByDB(dbCore)
	if err != nil {
		log.WithError(err).Panic("fail to new mysql adapter")
		panic(err)
	}

	e, err := casbin.NewEnforcer(cfg.Auth.RbacModelPath, adapter)
	if err != nil {
		log.WithError(err).Error("fail to new enforcer")
		panic(err)
	}

	err = e.LoadPolicy()
	if err != nil {
		log.WithError(err).Panic("fail to load policy")
		panic(err)
	}
	e.AddFunction("checkAdmin", func(args ...interface{}) (interface{}, error) {
		log.With("args", args).Debug("check admin")

		role := args[0].(string)
		return role == authV1.Role_admin.String(), nil
	})

	s := service{
		core:      e,
		jwtSecret: []byte(cfg.Auth.JwtSecret),
		db:        db,
		cache:     cache,
		log:       log.WithService("auth"),
	}
	s.loadDefaultPolicies()

	s.log.Debug("auth service initialized")
	return &s
}

func (s *service) loadDefaultPolicies() {
	pPolicy := [][]string{
		{authV1.Role_user.String(), "/api/v1/user/info", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/user/info", http.MethodPut},
		{authV1.Role_user.String(), "/api/v1/user/wallets", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/deposit", http.MethodPost},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/withdraw", http.MethodPost},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/history", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/currencies/:currency/history", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/transfer", http.MethodPost},
	}

	if _, err := s.core.RemovePolicy(authV1.Role_user.String()); err != nil {
		s.log.WithError(err).Warn("fail to clear policies from db")
	}

	for _, p := range pPolicy {
		ok, err := s.core.AddPolicy(p)
		if err != nil {
			s.log.WithError(err).Panic("fail to add policy")
			panic(err)
		}
		if !ok {
			s.log.With("policy", p).Info("policy already exists")
		}
	}
}

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
