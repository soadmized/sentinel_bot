package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Token  string `envconfig:"TOKEN"`
	UserID int64  `envconfig:"USER_ID"`
	Debug  bool   `envconfig:"DEBUG"`
}

func Load() Config {
	conf := Config{}

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		panic(err)
	}

	return conf
}
