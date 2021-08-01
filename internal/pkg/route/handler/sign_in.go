package handler

import (
	"backend/internal/pkg/auth"
	"backend/internal/pkg/config"
	"backend/internal/pkg/db"
	"backend/internal/pkg/sms"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/thirdParty/facebook"
	"backend/internal/pkg/user"
	"backend/internal/pkg/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

// MobileSignIn
// @Summary Mobile sign in
// @Tags v1,SignIn,Sms
// @version 1.0
// @Param smsRequest body sms.Request true "Sms request"
// @produce json
// @Success 204
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /api/v1/sign-up/confirmation [post]
func MobileSignIn(c *gin.Context) {
	log.Debug("handle user sign in")

	// check request body
	var r sms.Request
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// check user exist
	filter := db.GetUserFilter{
		Mobile: r.Mobile,
	}
	exist, err := user.Exist(filter)
	if err != nil {
		log.Error("fail to check user existence: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}
	if !exist {
		log.Warn("user not exist")
		c.JSON(http.StatusBadRequest, t7Error.UserNotfound)
		return
	}

	// send verify code
	if err := sms.SendVerifyCode(user.SignInPrefix, r.Mobile, util.GenVerifyCode()); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

// MobileSignInConfirm
// @Summary Mobile sign in confirm
// @Tags v1,SignIn,Sms
// @version 1.0
// @Param smsConfirm body sms.Confirm true "Sms confirm"
// @produce json
// @Success 200 {object} collection.Token "Token object"
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /api/v1/sign-up/confirmation [post]
func MobileSignInConfirm(c *gin.Context) {
	log.Debug("handle mobile sign in confirm")

	// check request cody
	var r sms.Confirm
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	// confirm verify code
	confirm, err := sms.ConfirmVerifyCode(user.SignInPrefix, r.Mobile, r.Code)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	if !confirm {
		c.JSON(http.StatusUnauthorized, t7Error.IncorrectVerifyCode)
		return
	}

	// get user data
	data, err := user.GetByMobile(r.Mobile)
	if err != nil {
		log.Error("fail to get user: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}

	// gen token
	log.Debug("gen user token")
	token, err := auth.GenUserToken(data.Id.Hex())
	if err != nil {
		log.Error("fail")
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, token)
	return
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
	userData, err := facebook.New().SignIn(code)
	if err != nil {
		log.Error("fail to sign facebook user: ", err.Error())
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, userData.String())

	//data, err := user.GetByChannel(collection.LoginChannelFacebook, userData.Id)
	//
	//// return user token if user exist
	//if err == nil {
	//	token, err := auth.GenUserToken(data.Id.Hex())
	//	if err != nil {
	//		log.Error("fail to get user token: ", err.Error())
	//		c.JSON(err.GetStatus(), err)
	//		return
	//	}
	//	c.JSON(http.StatusOK, token)
	//	return
	//}
	//
	//// sign up new user if user not found
	//if err.Code == t7Error.UserNotfound.Code {
	//	log.Debug("new user sign up from facebook")
	//
	//	u := collection.User{
	//		BasicInfo: collection.UserInfo{
	//			NickName: userData.Name,
	//			Gender:   userData.GetGender(),
	//			Birthday: userData.GetBirthday(),
	//		},
	//		Email: userData.Email,
	//		LoginClient: collection.LoginInfo{
	//			Channel:       collection.LoginChannelFacebook,
	//			ChannelUserId: userData.Id,
	//		},
	//	}
	//
	//	userId, err := user.CreateUser(u)
	//	if err != nil {
	//		log.Error("fail to create user: ", err.Error())
	//		c.JSON(err.GetStatus(), err)
	//		return
	//	}
	//	token, err := auth.GenUserToken(userId.Hex())
	//	if err != nil {
	//		log.Error("fail to get user token: ", err.Error())
	//		c.JSON(err.GetStatus(), err)
	//		return
	//	}
	//	c.JSON(http.StatusOK, token)
	//	return
	//}
	//
	//// other error
	//log.Error("fail to get user id: ", err.Error())
	//c.JSON(err.GetStatus(), err)
	//return
}
