package handler

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/sms"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/user"
	"github.com/Template7/backend/internal/pkg/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// SendVerifyCode
// @Summary Send verify code to the user mobile
// @Tags Sms,SignUp
// @version 1.0
// @Param smsRequest body sms.Request true "Sms request"
// @produce json
// @Success 204
// @failure 400 {object} t7Error.Error
// @Router /api/v1/sign-up/verification [post]
func SendVerifyCode(c *gin.Context) {
	log.Debug("handle send verify code")

	// check request body
	var r sms.Request
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// send verify code
	if err := sms.SendVerifyCode(user.SignUpPrefix, r.Mobile, util.GenVerifyCode()); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// ConfirmVerifyCode
// @Summary Confirm verify code
// @Tags Sms,SignUp
// @version 1.0
// @Param smsConfirm body sms.Confirm true "Sms confirm"
// @produce json
// @Success 200 {object} collection.Token "Token object"
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /api/v1/sign-up/confirmation [post]
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
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
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
