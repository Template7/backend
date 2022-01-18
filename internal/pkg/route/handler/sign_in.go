package handler

import (
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/sms"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/thirdParty/facebook"
	"github.com/Template7/backend/internal/pkg/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

// MobileSignIn
// @Summary Mobile sign in
// @Tags v1,SignIn,Sms
// @version 1.0
// @Param smsConfirm body sms.Confirm true "Sms confirm"
// @produce json
// @Success 200 {object} structs.Token "Token object"
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /api/v1/sign-in/mobile [post]
func MobileSignIn(c *gin.Context) {
	log.Debug("handle mobile sign in confirm")

	// check request cody
	var r sms.Confirm
	if err := c.ShouldBindJSON(&r); err != nil {
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
		c.JSON(http.StatusUnauthorized, t7Error.IncorrectVerifyCode)
		return
	}

	// sign in with mobile
	userToken, err := user.MobileSignIn(r.Mobile)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, userToken)
}

func FacebookSignInHome(c *gin.Context) {
	log.Debug("handle facebook sign in home")

	c.HTML(http.StatusOK, "facebook_login.html", nil)
	return
}

func FacebookSignIn(c *gin.Context) {
	log.Debug("handle facebook sign in")

	u, _ := url.Parse(facebook.OauthConf.Endpoint.AuthURL)
	parameters := url.Values{}
	parameters.Add("client_id", config.New().Facebook.AppId)
	parameters.Add("scope", strings.Join(facebook.OauthConf.Scopes, " "))
	parameters.Add("redirect_uri", facebook.OauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	//parameters.Add("state", oauthStateString)
	u.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, u.String())
	return
}

func FacebookSignInCallback(c *gin.Context) {
	log.Debug("handle facebook sign in callback")

	code := c.Query("code")

	// sign in from facebook
	userToken, err := facebook.New().SignIn(code)
	if err != nil {
		log.Error("fail to sign facebook user: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, userToken)
}
