package bot

import (
	"context"
	"fmt"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/pkg/errors"
	"github.com/soadmized/sentinel/pkg/dataset"

	"sentinel_bot/internal/config"
)

const (
	lastValuesCommand = "/last_values"
	triggerCommand    = "/trigger"
	queryPrefix       = "sensor_"
)

type Provider interface { // provide data from sentinel server
	LastValues(ctx context.Context, id string) (*dataset.Dataset, error)
	SensorIDs(ctx context.Context) ([]string, error)
}

type SentinelBot struct {
	bot      *tgbot.Bot
	Provider Provider
}

func New(conf config.Config, provider Provider) (*SentinelBot, error) {
	mw := Middleware{AllowedUsers: conf.AllowedUsers}

	opts := []tgbot.Option{
		tgbot.WithMiddlewares(mw.checkAuth),
	}

	bot, err := tgbot.New(conf.Token, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create bot")
	}

	return &SentinelBot{
		bot:      bot,
		Provider: provider,
	}, nil
}

func (s *SentinelBot) Run(ctx context.Context) {
	s.bot.RegisterHandler(tgbot.HandlerTypeMessageText, lastValuesCommand, tgbot.MatchTypeExact, s.lastValues)
	s.bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, queryPrefix, tgbot.MatchTypePrefix, s.lastValuesCallback)

	s.bot.RegisterHandler(tgbot.HandlerTypeMessageText, triggerCommand, tgbot.MatchTypeExact, s.triggerSentinel)
	s.bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, queryPrefix, tgbot.MatchTypePrefix, s.triggerSentinelCallback)

	s.bot.Start(ctx)
}

func (s *SentinelBot) errorHandler(ctx context.Context, chatID int64, err error) {
	s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:errcheck,exhaustruct
		ChatID: chatID,
		Text:   fmt.Sprintf("error occurred %s, sentinel is offline", err.Error()),
	})

	log.Fatal(err)
}
