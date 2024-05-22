package config

import (
	"github.com/spf13/viper"
	"log"
)

const (
	configPath = "config"
)

type Config struct {
	Env     string
	Service struct {
		Port int
	}
	Log struct {
		Level string
	}
}

func New() *Config {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")

	cfg := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	log.Println("config initialized")
	return cfg
}
