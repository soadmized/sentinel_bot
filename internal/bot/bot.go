package bot

import (
	"context"
	"fmt"
	"log"
	"time"

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
	s.bot.RegisterHandler(tgbot.HandlerTypeMessageText, "/sensors", tgbot.MatchTypeExact, s.sensorsHandler)
	s.bot.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "sensor", tgbot.MatchTypePrefix, s.sensorCallbackHandler)

	s.bot.Start(ctx)
}

// sensorCallbackHandler is callback handler for exact sensor query from inline keyboard.
func (s *SentinelBot) sensorCallbackHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	data, err := s.Provider.LastValues(ctx, "first") // TODO sensorID
	if err != nil {
		s.errorHandler(ctx, update.Message.Chat.ID, err)
	}

	payload := fmt.Sprintf("Last values (at %s) from %s sensor: temp = %.2f, light = %d, motion = %t",
		data.UpdatedAt.Local().Format(time.RFC822), //nolint:gosmopolitan
		data.ID,
		data.Temp,
		data.Light,
		data.Motion,
	)

	_, err = b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{ //nolint:exhaustruct
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		s.errorHandler(ctx, update.Message.Chat.ID, err)
	}

	_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:exhaustruct
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   payload,
	})
	if err != nil {
		s.errorHandler(ctx, update.Message.Chat.ID, err)
	}
}

// sensorsHandler is handler for /sensors command.
func (s *SentinelBot) sensorsHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	ids, err := s.Provider.SensorIDs(ctx)
	if err != nil {
		s.errorHandler(ctx, update.Message.Chat.ID, err)
	}

	_, err = s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:exhaustruct
		ChatID:      update.Message.Chat.ID,
		Text:        "Available sensors",
		ReplyMarkup: createSensorsKeyboard(ids),
	})
	if err != nil {
		s.errorHandler(ctx, update.Message.Chat.ID, err)
	}
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

func (s *SentinelBot) errorHandler(ctx context.Context, chatID int64, err error) {
	s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:errcheck,exhaustruct
		ChatID: chatID,
		Text:   fmt.Sprintf("error occurred %s, sentinel is offline", err.Error()),
	})

	log.Fatal(err)
}
