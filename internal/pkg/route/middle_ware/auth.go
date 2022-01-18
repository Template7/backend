package middle_ware

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	log = logger.GetLogger()
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

	c.Set("userId", utc.UserId)

	//if utc.UserId != c.Param("userId") {
	//	log.Warn("user id not match with claim")
	//	c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
	//	c.Abort()
	//	return
	//}

	log.Debug("user authorized")
}

func AuthActiveUser(c *gin.Context) {
	log.Debug("auth active user")

	userId, exist := c.Get("userId")
	if !exist || userId == "" {
		log.Warn("userId not found")
		c.Abort()
		return
	}

	user, dbErr := db.New().GetUserById(userId.(string))
	if dbErr != nil {
		log.Warn("fail to get user by id: ", dbErr.Error())
		c.JSON(http.StatusInternalServerError, t7Error.DbOperationFail.WithDetail(dbErr.Error()))
		return
	}

	if user.Status != structs.UserStatusActivate {
		log.Info("non-active user status: ", user.Status, ". abort further process")
		c.JSON(http.StatusForbidden, t7Error.UnAuthorized)
		return
	}

	log.Debug("active user authorized")
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
