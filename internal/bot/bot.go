package bot

import (
	"context"
	"fmt"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/soadmized/sentinel/pkg/dataset"

	"sentinel_bot/internal/config"
)

type Provider interface { // provide data from sentinel server
	LastValues(ctx context.Context, id string) (*dataset.Dataset, error)
	SensorIDs(ctx context.Context) ([]string, error)
}

type SentinelBot struct {
	bot          *tgbot.Bot
	allowedUsers []int64
	Provider     Provider
}

func New(conf config.Config, provider Provider) (*SentinelBot, error) {
	mw := Middleware{AllowedUsers: conf.AllowedUsers}

	opts := []tgbot.Option{
		tgbot.WithMiddlewares(mw.checkAuth),
		tgbot.WithCallbackQueryDataHandler("button", tgbot.MatchTypePrefix, callbackHandler),
		tgbot.WithCallbackQueryDataHandler("button", tgbot.MatchTypePrefix, callbackHandler),
	}

	bot, err := tgbot.New(conf.Token, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create bot")
	}

	return &SentinelBot{
		bot:          bot,
		allowedUsers: conf.AllowedUsers,
		Provider:     provider,
	}, nil
}

func (s *SentinelBot) Run(ctx context.Context) {
	s.bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/sensors", tgbot.MatchTypeExact, s.sensorsHandler)

	s.bot.Start(ctx)
}

func callbackHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{ //nolint:errcheck,exhaustruct
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:errcheck,exhaustruct
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}

func (s *SentinelBot) sensorsHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	ids, err := s.Provider.SensorIDs(ctx)
	if err != nil {
		s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:errcheck,exhaustruct
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("error occurred %s, sentinel is offline", err.Error()),
		})

		log.Fatal(err)
	}

	kb := createSensorsKeyboard(ids)

	b.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:errcheck,exhaustruct
		ChatID:      update.Message.Chat.ID,
		Text:        "Available sensors",
		ReplyMarkup: kb,
	})
}

func createSensorsKeyboard(ids []string) *models.InlineKeyboardMarkup {
	buttons := make([][]models.InlineKeyboardButton, 0, len(ids))

	for _, id := range ids {
		button := []models.InlineKeyboardButton{
			{ //nolint:exhaustruct
				Text:         id,
				CallbackData: fmt.Sprintf("sensor_%s", id),
			},
		}

		buttons = append(buttons, button)
	}

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}
