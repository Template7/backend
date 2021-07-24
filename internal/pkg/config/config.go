package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"strings"
	"sync"
)

const (
	configPath = "configs"
)

type config struct {
	Mode    string `yaml:"mode"`
	Port    int    `yaml:"port"`
	JwtSign string `yaml:"jwtSign"`
	Log     struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Mongo struct {
		Db               string `yaml:"db"`
		Host             string `yaml:"host"`
		Port             int    `yaml:"port"`
		Username         string `yaml:"username"`
		Password         string `yaml:"password"`
		ConnectionString string
	} `yaml:"mongo"`
	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		PollSize int    `yaml:"poll_size"`
		//ReadTimeout int    `yaml:"read_timeout"`
	} `yaml:"redis"`
	Facebook struct {
		AppId  string `yaml:"app_id"`
		Secret string `yaml:"secret"`
	} `yaml:"facebook"`
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
		instance.initLog()

		if instance.Mongo.Username != "" && instance.Mongo.Password != "" {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d", instance.Mongo.Username, instance.Mongo.Password, instance.Mongo.Host, instance.Mongo.Port)
		} else {
			instance.Mongo.ConnectionString = fmt.Sprintf("mongodb://%s:%d", instance.Mongo.Host, instance.Mongo.Port)
		}

		log.Debug("config initialized")
	})
	return instance
}

func (c *config) initLog() {
	logLevel := map[string]log.Level{
		"DEBUG": log.DebugLevel,
		"INFO":  log.InfoLevel,
		"WARN":  log.WarnLevel,
		"ERROR": log.ErrorLevel,
		"FATAL": log.FatalLevel,
	}

	callerFormatter := func(path string) string {
		arr := strings.Split(path, "/")
		return arr[len(arr)-1]
	}
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", callerFormatter(f.File), f.Line)
		},
		//PrettyPrint: true,
	}

	log.SetLevel(logLevel[c.Log.Level])
	log.SetFormatter(customFormatter)
	log.SetReportCaller(true)
	log.Debug("logger initialized")
}
