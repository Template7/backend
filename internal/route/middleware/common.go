package middleware

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

// keys set on the gin context by AuthToken
const (
	UserId = "userId"
	Role   = "role"
	Status = "status"

	// HeaderRequestId is kept for existing callers; prefer response.HeaderRequestId.
	HeaderRequestId = response.HeaderRequestId
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
		log:     log.WithService("middleware"),
	}
}

// CurrentUserId returns the id of the authenticated user set by AuthToken.
func CurrentUserId(c *gin.Context) (string, bool) {
	id := c.GetString(UserId)
	return id, id != ""
}

func currentRole(c *gin.Context) (string, bool) {
	r := c.GetString(Role)
	return r, r != ""
}

func currentStatus(c *gin.Context) (int, bool) {
	v, ok := c.Get(Status)
	if !ok {
		return 0, false
	}
	s, ok := v.(int)
	return s, ok
}
