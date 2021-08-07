package facebook

import (
	"github.com/Template7/backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	baseUrl        = "https://graph.facebook.com"
	parseTokenUri  = "/debug_token"
	accessTokenUri = "/v9.0/oauth/access_token"
)

var (
	once     sync.Once
	instance *client
	
	// TODO: update after communicate with app
	redirectUri = config.New().Facebook.Callback
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
