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
		Format  string // json(default) | console
		Level   string // debug(default) | info | warn | error
		Version string // commit id
	}
	Auth struct {
		RbacModelPath string
	}
	Db struct {
		Sql struct {
			Db         string
			Host       string
			Port       int
			Username   string
			Password   string
			Connection struct {
				Min int
				Max int
			}
		}
		NoSql struct {
			Db       string
			Host     string
			Port     int
			Username string
			Password string
		}
	}
	Cache struct {
		Host         string
		Password     string
		ReadTimeout  int
		WriteTimeout int
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
