package middleware

import (
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
	log     *logger.Logger
}

func New(userSvc *user.Service, log *logger.Logger) *Controller {
	return &Controller{
		userSvc: userSvc,
		log:     log.WithService("middleController"),
	}
}
