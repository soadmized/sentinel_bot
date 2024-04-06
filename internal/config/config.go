//nolint:stylecheck
package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Token        string  `envconfig:"TOKEN"`
	AllowedUsers []int64 `envconfig:"ALLOWED_USERS"`
	Debug        bool    `envconfig:"DEBUG"`

	SentinelServerUrl string `envconfig:"SENTINEL_SERVER_URL"`
	SentinelUser      string `envconfig:"SENTINEL_USER"`
	SentinelPass      string `envconfig:"SENTINEL_PASS"`
}

func Load() Config {
	conf := Config{} //nolint:exhaustruct

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		panic(err)
	}

	return conf
}
