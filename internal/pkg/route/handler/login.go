package handler

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/sms"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/thirdParty/facebook"
	"github.com/Template7/backend/internal/pkg/user"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// NativeLogin
// @Summary Native login
// @Tags v1,login
// @version 1.0
// @Param request body v1.LoginRequest true "Request"
// @produce json
// @Success 200 {object} v1.LoginResponse "Response"
// @failure 400 {object} t7Error.Error
// @Router /api/v1/login/native [post]
func NativeLogin(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle native login")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req v1.LoginRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	token, err := auth.New().Login(c, req.Username, req.Password)
	if err != nil {
		defer c.Abort()
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}

	c.JSON(http.StatusOK, v1.LoginResponse{
		Token: token,
	})
}

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
	log := logger.New().WithContext(c)
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
	log := logger.New().WithContext(c)
	log.Debug("handle facebook sign in home")

	c.HTML(http.StatusOK, "facebook_login.html", nil)
	return
}

func FacebookSignIn(c *gin.Context) {
	log := logger.New().WithContext(c)
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
	log := logger.New().WithContext(c)
	log.Debug("handle facebook sign in callback")

	code := c.Query("code")

	// sign in from facebook
	userToken, err := facebook.New().SignIn(code)
	if err != nil {
		log.WithError(err).Error("fail to sign facebook user")
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, userToken)
}
