package config

import (
	"fmt"
	"github.com/Template7/common/logger"
	"github.com/spf13/viper"
	"sync"
)

const (
	configPath = "configs"
	jwtSign    = "45519f46c06c8340a34f9a32982860c1a8d6bb57eaeb338b7f0119062b8a3b67"
)

var log = logger.GetLogger()

type config struct {
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
	MySql struct {
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
	instance *config
)

func New() *config {
	once.Do(func() {
		viper.SetConfigType("yaml")
		instance = &config{}
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("fail to load config file: ", err.Error())
		}
		if err := viper.Unmarshal(&instance); err != nil {
			log.Fatal(err)
		}
		instance.JwtSign = []byte(jwtSign)
		instance.initLog()

		if instance.Mongo.Username != "" && instance.Mongo.Password != "" {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d", instance.Mongo.Username, instance.Mongo.Password, instance.Mongo.Host, instance.Mongo.Port)
		} else {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%d", instance.Mongo.Host, instance.Mongo.Port)
		}
		instance.MySql.ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", instance.MySql.Username, instance.MySql.Password, instance.MySql.Host, instance.MySql.Port, instance.MySql.Db)

		log.Debug("config initialized")
	})
	return instance
}

func (c *config) initLog() {
	logger.SetLevel(c.Log.Level)
	logger.SetFormatter(c.Log.Formatter)
	log.Debug("logger initialized")
}
