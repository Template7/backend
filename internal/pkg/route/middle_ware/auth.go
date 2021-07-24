package middle_ware

import (
	"backend/internal/pkg/config"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserTokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

func AuthUser(c *gin.Context) {
	log.Debug("auth user")

	idToken := util.ParseBearerToken(c.GetHeader("Authorization"))
	token, err := jwt.ParseWithClaims(idToken, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.New().JwtSign), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
		return
	}

	_, ok := token.Claims.(*UserTokenClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
		return
	}

	//c.Set("claims", claims)
	//c.Keys = map[string]interface{}{
	//	"userId": claims.UserId,
	//	"claims": claims,
	//}
}

func AuthAdmin(c *gin.Context) {
	log.Debug("auth admin")

	// TODO: implementation
}
