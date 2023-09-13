package middleware

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Permission(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("check user permission")

	uId, ok := c.Get(UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		c.Abort()
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		c.Abort()
		return
	}
	_, ok = c.Get(Role)
	if !ok {
		log.Warn("no user role from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		c.Abort()
		return
	}
	_, ok = c.Get(Status)
	if !ok {
		log.Warn("no user status from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		c.Abort()
		return
	}

	if !auth.New().CheckPermission(c, userId, c.Request.URL.Path, c.Request.Method) {
		c.JSON(http.StatusOK, t7Error.UnAuthorized)
		c.Abort()
		return
	}
	c.Next()
}

func AuthToken(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("auth user token")

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		c.Abort()
		return
	}
	claims, err := auth.New().ParseToken(c, token)
	if err != nil {
		defer c.Abort()
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(http.StatusBadRequest, t7Err)
		return
	}

	c.Set(UserId, claims.UserId)
	c.Set(Role, claims.Role)
	c.Set(Status, claims.Status)

	log.With("userId", claims.UserId).Debug("user token authorized")
	c.Next()
}
