package internal

import (
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/common/cache"
	"github.com/Template7/common/db"
	"github.com/Template7/common/logger"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var CommitHash string

func ProvideSqlCore(cfg *config.Config, log *logger.Logger) *gorm.DB {
	return db.NewSql(
		cfg.Db.Sql.Host,
		cfg.Db.Sql.Port,
		cfg.Db.Sql.Username,
		cfg.Db.Sql.Password,
		cfg.Db.Sql.Db,
		cfg.Db.Sql.Connection.Min,
		cfg.Db.Sql.Connection.Max,
		log)
}

func ProvideNoSqlCore(cfg *config.Config, log *logger.Logger) *mongo.Client {
	return db.NewNoSql(cfg.Db.NoSql.Host, cfg.Db.NoSql.Port, cfg.Db.NoSql.Username, cfg.Db.NoSql.Password, log)
}

func ProvideCacheCore(cfg *config.Config, log *logger.Logger) *redis.Client {
	return cache.New(cfg.Cache.Host, cfg.Cache.Password, cfg.Cache.ReadTimeout, cfg.Cache.WriteTimeout, log)
}

func ProvideLogger(cfg *config.Config) *logger.Logger {
	return logger.New(cfg.Log.Level, cfg.Log.Format, CommitHash)
}
