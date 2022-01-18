package facebook

import (
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/common/logger"
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
	log = logger.GetLogger()
	
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
