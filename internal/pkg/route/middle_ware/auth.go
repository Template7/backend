package middle_ware

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthUserToken(c *gin.Context) {
	log.Debug("auth user")

	userToken := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(userToken, &auth.UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.New().JwtSign, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized.WithDetail(err.Error()))
		c.Abort()
		return
	}

	utc, ok := token.Claims.(*auth.UserTokenClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
		c.Abort()
		return
	}

	if utc.UserId != c.Param("user-id") {
		log.Warn("user id not match with claim")
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
		c.Abort()
		return
	}

	log.Debug("user authorized")
}

func AuthAdmin(c *gin.Context) {
	log.Debug("auth admin")

	adminToken := c.GetHeader("Authorization")
	_, err := jwt.Parse(adminToken, func(token *jwt.Token) (interface{}, error) {
		return config.New().JwtSign, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized.WithDetail(err.Error()))
		c.Abort()
		return
	}

	log.Debug("admin authorized")
}
