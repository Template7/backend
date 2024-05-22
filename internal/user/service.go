package user

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
)

type Service struct {
	authSvc auth.Auth
	db      db.Client
	log     *logger.Logger
}

func New(authSvc auth.Auth, db db.Client, log *logger.Logger) *Service {
	return &Service{
		db:      db,
		authSvc: authSvc,
		log:     log.WithService("user"),
	}
}
