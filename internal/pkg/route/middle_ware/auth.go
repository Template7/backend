package middle_ware

import (
	"backend/internal/pkg/config"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/user"
	"backend/internal/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthUserToken(c *gin.Context) {
	log.Debug("auth user")

	userToken := util.ParseBearerToken(c.GetHeader("Authorization"))
	token, err := jwt.ParseWithClaims(userToken, &user.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.New().JwtSign, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized.WithDetail(err.Error()))
		c.Abort()
		return
	}

	utc, ok := token.Claims.(*user.TokenClaims)
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
}

func AuthAdmin(c *gin.Context) {
	log.Debug("auth admin")

	// TODO: implementation
}
