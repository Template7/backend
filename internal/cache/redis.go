package cache

import (
	"github.com/Template7/common/cache"
	"github.com/Template7/common/logger"
	"github.com/redis/go-redis/v9"
)

type client struct {
	core *redis.Client
	log  *logger.Logger
}

func New(log *logger.Logger) Interface {
	return &client{
		core: cache.New(),
		log:  log.WithService("redis"),
	}
}
