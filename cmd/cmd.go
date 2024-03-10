package cmd

import (
	"sentinel_bot/internal/build"
	"sentinel_bot/internal/config"
)

func Run(conf config.Config) error {
	builder := build.New(conf)

	bot, err := builder.Bot()
	if err != nil {
		return err
	}

	bot.Run()

	return nil
}
