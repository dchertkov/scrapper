package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "SCRAPPER"

type (
	Config struct {
		Server  Server
		Checker Checker

		SourceFile string `envconfig:"SOURCE_FILE" default:"sites.txt"`
	}

	Server struct {
		Host string `envconfig:"SERVER_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"SERVER_PORT" default:"3000"`
	}

	Checker struct {
		Workers  int           `envconfig:"CHECKER_WORKERS" default:"20"`
		Timeout  time.Duration `envconfig:"CHECKER_TIMEOUT" default:"5s"`
		Interval time.Duration `envconfig:"CHECKER_INTERVAL" default:"1m"`
	}
)

func (s Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func FromEnv() (*Config, error) {
	conf := &Config{}
	if err := envconfig.Process(envPrefix, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
