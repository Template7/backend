package handler

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/sms"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ConfirmVerifyCode
// @Summary Confirm verify code
// @Tags Sms,SignUp
// @version 1.0
// @Param smsConfirm body sms.Confirm true "Sms confirm"
// @produce json
// @Success 200 {object} structs.Token "Token object"
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
	confirm, err := sms.ConfirmVerifyCode(sms.VerifyCodePrefix, r.Mobile, r.Code)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	if !confirm {
		c.JSON(http.StatusUnauthorized, t7Error.UnAuthorized)
		return
	}

	// create native user
	userId, err := user.CreateNativeUser(r.Mobile)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	// gen token
	log.Debug("gen user token")
	token, err := auth.GenUserToken(userId)
	if err != nil {
		log.Error("fail to gen user token: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, token)
	return
}
