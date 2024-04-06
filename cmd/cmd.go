package cmd

import (
	"context"

	"github.com/pkg/errors"

	"sentinel_bot/internal/build"
	"sentinel_bot/internal/config"
)

func Run(conf config.Config) error {
	builder := build.New(conf)

	bot, err := builder.Bot()
	if err != nil {
		return errors.Wrap(err, "create bot")
	}

	ctx := context.Background()
	bot.Run(ctx)

	return nil
}
