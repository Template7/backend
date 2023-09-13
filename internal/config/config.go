package config

import (
	"fmt"
	"github.com/Template7/common/logger"
	"github.com/spf13/viper"
	"log"
	"sync"
)

const (
	configPath = "configs"
)

type Config struct {
	JwtSign []byte
	Log     struct {
		Formatter string
		Level     string
	}
	Gin struct {
		ListenPort      int
		Mode            string
		ShutdownTimeout int
	}
	Mongo struct {
		Db               string
		Host             string
		Port             int
		Username         string
		Password         string
		ConnectionString string
	}
	Sql struct {
		Db               string
		Host             string
		Port             int
		Username         string
		Password         string
		ConnectionString string
	}
	Redis struct {
		Host     string
		Password string
		//PollSize int
		//ReadTimeout int
	}
	Facebook struct {
		AppId    string
		Secret   string
		Callback string
	}
}

var (
	once     sync.Once
	instance *Config
)

func New() *Config {
	once.Do(func() {
		viper.SetConfigType("yaml")
		instance = &Config{}
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("fail to load config file: ", err.Error())
		}
		if err := viper.Unmarshal(&instance); err != nil {
			log.Fatal(err)
		}

		if instance.Mongo.Username != "" && instance.Mongo.Password != "" {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d", instance.Mongo.Username, instance.Mongo.Password, instance.Mongo.Host, instance.Mongo.Port)
		} else {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%d", instance.Mongo.Host, instance.Mongo.Port)
		}
		instance.Sql.ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", instance.Sql.Username, instance.Sql.Password, instance.Sql.Host, instance.Sql.Port, instance.Sql.Db)

		logger.New().Info("config initialized")
	})
	return instance
}
