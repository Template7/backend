package facebook

import (
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/db/collection"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/util"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"net/http"
	"time"
)

const (
	birthdayLayout = "01/02/2006"
)

var (
	OauthConf = &oauth2.Config{
		ClientID: config.New().Facebook.AppId,
		//ClientSecret: config.New().Facebook.Secret,
		Scopes:      []string{"public_profile"},
		RedirectURL: config.New().Facebook.Callback,
		Endpoint:    facebook.Endpoint,
	}
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type tokenData struct {
	Data struct {
		AppID               string   `json:"app_id"`
		Type                string   `json:"type"`
		Application         string   `json:"application"`
		DataAccessExpiresAt int      `json:"data_access_expires_at"`
		ExpiresAt           int      `json:"expires_at"`
		IsValid             bool     `json:"is_valid"`
		IssuedAt            int      `json:"issued_at"`
		Scopes              []string `json:"scopes"`
		UserID              string   `json:"user_id"`
	} `json:"data"`
}

type basicUserData struct {
	Id       string `json:"id"` // facebook user id
	Name     string `json:"name"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
}

func (b basicUserData) String() string {
	return fmt.Sprintf("id: %s, name: %s, email: %s, birthday: %s, gender: %s", b.Id, b.Name, b.Email, b.Birthday, b.Gender)
}

func (b basicUserData) GetBirthday() (birthday int64) {
	t, err := time.Parse(birthdayLayout, b.Birthday)
	if err != nil {
		log.Warn("fail to parse birthday: ", b.Birthday, ", ", err.Error())
		return
	}
	birthday = t.Unix()
	return
}

func (b basicUserData) GetGender() (gender collection.Gender) {
	switch b.Gender {
	case "male":
		gender = collection.GenderMale
	case "female":
		gender = collection.GenderFemale
	default:
		gender = collection.GenderUnknown
	}
	return
}

func (c client) SignIn(code string) (userData basicUserData, err *t7Error.Error) {
	log.Debug("sign in facebook user")

	token, err := c.getAccessToken(code)
	if err != nil {
		return
	}

	tokenData, err := c.getTokenInfo(token)
	if err != nil {
		return
	}

	userData, err = c.getUserData(token, tokenData)
	if err != nil {
		return
	}
	return
}

func (c client) getAccessToken(code string) (token accessToken, err *t7Error.Error) {
	log.Debug("get facebook access token")

	// construct request
	req, _ := http.NewRequest(http.MethodGet, baseUrl+accessTokenUri, nil)
	q := req.URL.Query()
	q.Set("client_id", c.AppId)
	q.Set("redirect_uri", redirectUri)
	q.Set("client_secret", c.Secret)
	q.Set("code", code)
	req.URL.RawQuery = q.Encode()

	resp, err := util.SendHttpRequest(req)
	if err != nil {
		return
	}

	if mErr := json.Unmarshal(resp, &token); mErr != nil {
		log.Error("fail to unmarshal response: ", mErr.Error())
		err = t7Error.DecodeFail.WithDetailAndStatus(mErr.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("get facebook access token successfully")
	return
}

func (c client) getTokenInfo(token accessToken) (tokenData tokenData, err *t7Error.Error) {
	log.Debug("get facebook user data")

	// construct request
	req, _ := http.NewRequest(http.MethodGet, baseUrl+parseTokenUri, nil)
	q := req.URL.Query()
	q.Set("input_token", token.AccessToken)
	q.Set("access_token", token.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := util.SendHttpRequest(req)
	if err != nil {
		return
	}

	if mErr := json.Unmarshal(resp, &tokenData); mErr != nil {
		log.Error("fail to unmarshal response: ", mErr.Error())
		err = t7Error.DecodeFail.WithDetailAndStatus(mErr.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("get facebook token data successfully")
	return
}

func (c client) getUserData(token accessToken, data tokenData) (userDara basicUserData, err *t7Error.Error) {
	log.Debug("get facebook user data")

	// construct request
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrl, data.Data.UserID), nil)
	q := req.URL.Query()
	q.Set("fields", "id,name,email,birthday,gender")
	q.Set("input_token", token.AccessToken)
	q.Set("access_token", token.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := util.SendHttpRequest(req)
	if err != nil {
		return
	}

	if mErr := json.Unmarshal(resp, &userDara); mErr != nil {
		log.Error("fail to unmarshal response: ", mErr.Error())
		err = t7Error.DecodeFail.WithDetailAndStatus(mErr.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("get user data successfully: ", userDara.String())
	return
}
