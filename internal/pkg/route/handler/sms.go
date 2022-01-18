package handler

import (
	"github.com/Template7/backend/internal/pkg/sms"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendSmsVerifyCode
// @Summary Send verify code to the user mobile
// @Tags Sms
// @version 1.0
// @Param smsRequest body sms.Request true "Sms request"
// @produce json
// @Success 204
// @failure 400 {object} t7Error.Error
// @Router /api/v1/verify-code/sms [post]
func SendSmsVerifyCode(c *gin.Context) {
	log.Debug("handle send sms verify code")

	// check request body
	var r sms.Request
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// send verify code
	if err := sms.SendVerifyCode(sms.VerifyCodePrefix, r.Mobile, util.GenVerifyCode()); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

func SendEmailVerifyCode(c *gin.Context) {
	log.Debug("handle send email verify code")

	return
}
