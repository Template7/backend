package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

const (
	configPath = "config"

	// envPrefix: any config key can be overridden by an env var,
	// e.g. db.sql.password -> T7_DB_SQL_PASSWORD, auth.jwtSecret -> T7_AUTH_JWTSECRET
	envPrefix = "T7"
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
		JwtSecret     string
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
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

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
