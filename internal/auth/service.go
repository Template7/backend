package auth

import (
	"fmt"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/config"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
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
	UserId string `json:"userId"`
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
		cs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Db.Sql.Username, cfg.Db.Sql.Password, cfg.Db.Sql.Host, cfg.Db.Sql.Port, cfg.Db.Sql.Db)
		adapter, err := gormadapter.NewAdapter("mysql", cs, true)
		if err != nil {
			log.WithError(err).Panic("fail to new mysql adapter")
			panic(err)
		}

		e, err := casbin.NewEnforcer()
		if err != nil {
			log.WithError(err).Panic("fail to new enforcer")
			panic(err)
		}

		ms := "[request_definition]\nr = sub, obj, act\n\n[policy_definition]\np = sub, obj, act\n\n[role_definition]\ng = _, _\n\n[policy_effect]\ne = some(where (p.eft == allow))\n\n[matchers]\nm = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || checkAdmin(r.sub)"
		md, err := model.NewModelFromString(ms)
		if err != nil {
			log.WithError(err).Panic("fail to new model")
			panic(err)
		}
		if err := e.InitWithModelAndAdapter(md, adapter); err != nil {
			log.WithError(err).Panic("fail to init enforcer")
			panic(err)
		}

		//e, err := casbin.NewEnforcer("./config/rbac_model.conf", adapter)
		//if err != nil {
		//	log.WithError(err).Error("fail to new enforcer")
		//	panic(err)
		//}

		err = e.LoadPolicy()
		if err != nil {
			log.WithError(err).Panic("fail to load policy")
			panic(err)
		}
		e.AddFunction("checkAdmin", func(args ...interface{}) (interface{}, error) {
			username := args[0].(string)
			return e.HasRoleForUser(username, v1.Role_name[int32(v1.Role_admin)])
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

func (s *service) loadDefaultPolicies() {
	pPolicy := [][]string{
		{authV1.Role_user.String(), "/api/v1/users/:userId/info", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/users/:userId/info", http.MethodPut},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId", http.MethodGet},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/deposit", http.MethodPost},
		{authV1.Role_user.String(), "/api/v1/wallets/:walletId/withdraw", http.MethodPost},
		{authV1.Role_user.String(), "/api/v1/transaction", http.MethodPost},
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
