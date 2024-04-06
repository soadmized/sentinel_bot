package build

import (
	"github.com/pkg/errors"

	"sentinel_bot/internal/bot"
	"sentinel_bot/internal/config"
	"sentinel_bot/internal/provider"
)

type Builder struct {
	conf config.Config
}

func New(conf config.Config) *Builder {
	return &Builder{
		conf: conf,
	}
}

func (b *Builder) Bot() (*bot.SentinelBot, error) {
	p := b.provider()

	sentinelBot, err := bot.New(b.conf, p)
	if err != nil {
		return nil, errors.Wrap(err, "build sentinel sentinelBot")
	}

	return sentinelBot, nil
}

func (b *Builder) provider() *provider.Provider {
	return provider.New(b.conf)
}
