package middleware

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
)

const (
	UserId = "userId"
	Role   = "role"
	Status = "status"

	HeaderRequestId = "request_id"
)

type Controller struct {
	userSvc *user.Service
	authSvc auth.Auth
	log     *logger.Logger
}

func New(userSvc *user.Service, authSvc auth.Auth, log *logger.Logger) *Controller {
	return &Controller{
		userSvc: userSvc,
		authSvc: authSvc,
		log:     log.WithService("middleController"),
	}
}
