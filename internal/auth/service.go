package auth

import (
	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const (
	minJwtSecretLen = 32
	bcryptCost      = bcrypt.DefaultCost

	defaultMaxFailuresPerUser = 5
	defaultMaxFailuresPerIp   = 20
	defaultFailureWindow      = 15 * time.Minute
)

// dummyPasswordHash is compared against when the user doesn't exist,
// so a login for an unknown user takes as long as a wrong password.
var dummyPasswordHash, _ = HashPassword("dummy-password-for-timing")

type UserTokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Status int    `json:"status"`
}

type service struct {
	core      *casbin.Enforcer
	jwtSecret []byte
	db        db.UserRepo
	cache     cache.Interface
	limits    loginLimits
	log       *logger.Logger
}

type loginLimits struct {
	perUser int64
	perIp   int64
	window  time.Duration
}

func New(db db.Client, dbCore *gorm.DB, cache cache.Interface, log *logger.Logger, cfg *config.Config) Auth {
	if len(cfg.Auth.JwtSecret) < minJwtSecretLen {
		log.With("minLength", minJwtSecretLen).Panic("jwt secret missing or too short, set auth.jwtSecret / T7_AUTH_JWTSECRET")
	}

	adapter, err := gormadapter.NewAdapterByDB(dbCore)
	if err != nil {
		log.WithError(err).Panic("fail to new mysql adapter")
	}

	e, err := casbin.NewEnforcer(cfg.Auth.RbacModelPath, adapter)
	if err != nil {
		log.WithError(err).Panic("fail to new enforcer")
	}

	err = e.LoadPolicy()
	if err != nil {
		log.WithError(err).Panic("fail to load policy")
	}
	e.AddFunction("checkAdmin", func(args ...interface{}) (interface{}, error) {
		role, _ := args[0].(string)
		return role == authV1.Role_admin.String(), nil
	})

	s := service{
		core:      e,
		jwtSecret: []byte(cfg.Auth.JwtSecret),
		db:        db,
		cache:     cache,
		limits:    newLoginLimits(cfg),
		log:       log.WithService("auth"),
	}

	s.log.Debug("auth service initialized")
	return &s
}

// HashPassword returns the bcrypt hash to store for a new password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func newLoginLimits(cfg *config.Config) loginLimits {
	l := loginLimits{
		perUser: int64(cfg.Auth.Login.MaxFailuresPerUser),
		perIp:   int64(cfg.Auth.Login.MaxFailuresPerIp),
		window:  time.Duration(cfg.Auth.Login.WindowMinutes) * time.Minute,
	}
	if l.perUser <= 0 {
		l.perUser = defaultMaxFailuresPerUser
	}
	if l.perIp <= 0 {
		l.perIp = defaultMaxFailuresPerIp
	}
	if l.window <= 0 {
		l.window = defaultFailureWindow
	}
	return l
}
