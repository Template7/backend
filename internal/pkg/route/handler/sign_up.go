package handler

import (
	"backend/internal/pkg/auth"
	"backend/internal/pkg/sms"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/user"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendVerifyCode(c *gin.Context) {
	log.Debug("handle send verify code")

	// check request body
	var r sms.Request
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// send verify code
	if err := sms.SendVerifyCode(user.SignUpPrefix, r.Mobile, sms.GenVerifyCode()); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// return token string if confirmed
func ConfirmVerifyCode(c *gin.Context) {
	log.Debug("handle confirm verify code")

	// check request cody
	var r sms.Confirm
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// confirm verify code
	confirm, err := sms.ConfirmVerifyCode(user.SignUpPrefix, r.Mobile, r.Code)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	if !confirm {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// create native user
	userData, err := user.CreateNativeUser(r.Mobile)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	// gen token
	log.Debug("gen user token")
	token, err := auth.GenUserToken(userData.Id.Hex())
	if err != nil {
		log.Error("fail to gen user token: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, token)
	return
}
