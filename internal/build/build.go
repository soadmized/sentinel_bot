package build

import (
	"github.com/pkg/errors"
	"sentinel_bot/internal/api"
	"sentinel_bot/internal/config"
)

type Builder struct {
	conf config.Config
}

func New(conf config.Config) *Builder {
	return &Builder{
		conf: conf,
	}
}

func (b *Builder) Bot() (*api.Sentinel, error) {
	bot, err := api.New(b.conf)
	if err != nil {
		return nil, errors.Wrap(err, "build sentinel bot")
	}

	return bot, nil
}
