package facebook

import (
	"backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	baseUrl = "https://graph.facebook.com"
	parseTokenUri = "/debug_token"
	accessTokenUri = "/v9.0/oauth/access_token"

	// TODO: update after communicate with app
	redirectUri = "http://localhost:8080/auth/facebook"
)

var (
	once     sync.Once
	instance *client
)

type client struct {
	AppId  string
	Secret string
}

func New() *client {
	once.Do(func() {
		instance = &client{
			AppId:  config.New().Facebook.AppId,
			Secret: config.New().Facebook.Secret,
		}
		log.Debug("facebook client initialized")
	})
	return instance
}
