package cache

import (
	"github.com/Template7/common/logger"
	"github.com/redis/go-redis/v9"
)

type client struct {
	core *redis.Client
	log  *logger.Logger
}

func New(core *redis.Client, log *logger.Logger) Interface {
	return &client{
		core: core,
		log:  log.WithService("redis"),
	}
}
