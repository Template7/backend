package middleware

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Permission(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("check user permission")

	//uId, ok := c.Get(UserId)
	//if !ok {
	//	log.Warn("no user id from the previous middleware")
	//	c.JSON(http.StatusUnauthorized, types.HttpRespBase{
	//		RequestId: c.GetHeader(HeaderRequestId),
	//		Code:      int(t7Error.InvalidToken.Code),
	//		Message:   t7Error.InvalidToken.Message,
	//	})
	//	c.Abort()
	//	return
	//}
	//userId, ok := uId.(string)
	//if !ok {
	//	log.With("uId", uId).Error("type assertion fail")
	//	c.JSON(http.StatusUnauthorized, types.HttpRespBase{
	//		RequestId: c.GetHeader(HeaderRequestId),
	//		Code:      int(t7Error.InvalidToken.Code),
	//		Message:   t7Error.InvalidToken.Message,
	//	})
	//	c.Abort()
	//	return
	//}

	role, ok := c.Get(Role)
	if !ok {
		log.Warn("no user role from the previous middleware")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		c.Abort()
		return
	}
	userRole, ok := role.(string)
	if !ok {
		log.With("role", role).Error("type assertion fail")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		c.Abort()
		return
	}

	if !auth.New().CheckPermission(c, userRole, c.Request.URL.Path, c.Request.Method) {
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
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
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		c.Abort()
		return
	}
	claims, err := auth.New().ParseToken(c, token)
	if err != nil {
		defer c.Abort()
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.Set(UserId, claims.UserId)
	c.Set(Role, claims.Role)

	log.With("userId", claims.UserId).Debug("user token authorized")
	c.Next()
}

func AuthUserWallet(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("check user wallet")

	walletId := c.Param("walletId")
	if walletId == "" {
		log.Warn("empty wallet id")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		c.Abort()
		return
	}
	log = log.With("walletId", walletId)

	uId, ok := c.Get(UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		c.Abort()
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		c.Abort()
		return
	}
	log = log.With("userId", userId)

	for _, uw := range user.New().GetUserWallets(c, userId) {
		if uw.Id == walletId {
			log.Debug("user wallet check ok")
			c.Next()
			return
		}
	}

	log.Warn("user has no permission to the wallet")
	c.JSON(http.StatusUnauthorized, types.HttpRespBase{
		RequestId: c.GetHeader(HeaderRequestId),
		Code:      int(t7Error.InvalidToken.Code),
		Message:   t7Error.InvalidToken.Message,
	})
	c.Abort()
	return
}
