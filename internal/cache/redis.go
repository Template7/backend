package cache

import (
	"github.com/Template7/common/cache"
	"github.com/Template7/common/logger"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	once     sync.Once
	instance *client
)

type client struct {
	core *redis.Client
	log  *logger.Logger
}

func New() Interface {
	once.Do(func() {
		log := logger.New().WithService("redis")
		instance = &client{
			core: cache.New(),
			log:  log,
		}

		log.Info("redis client initialized")
	})
	return instance
}
