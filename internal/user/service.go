package user

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	"sync"
)

var (
	once     sync.Once
	instance *Service
)

type Service struct {
	authSvc auth.Auth
	db      db.Client
	log     *logger.Logger
}

func New() *Service {
	once.Do(func() {
		log := logger.New().WithService("user")
		instance = &Service{
			authSvc: auth.New(),
			db:      db.New(),
			log:     log,
		}
		log.Info("user service initialized")
	})
	return instance
}
