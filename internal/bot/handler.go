package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
)

// lastValuesCallback is callback handler for exact sensor last values query from inline keyboard.
func (s *SentinelBot) lastValuesCallback(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	id, err := sensorIDFromCallback(update.CallbackQuery.Data)
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}

	data, err := s.Provider.LastValues(ctx, id)
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}

	payload := fmt.Sprintf("Last values (at %s) from %s sensor: temp = %.2f, light = %d, motion = %t",
		data.UpdatedAt.Local().Format(time.RFC822), //nolint:gosmopolitan
		data.ID,
		data.Temp,
		data.Light,
		data.Motion,
	)

	_, err = s.bot.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{ //nolint:exhaustruct
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}

	_, err = s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:exhaustruct
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   payload,
	})
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}
}

// lastValues is handler for /last_values command.
func (s *SentinelBot) lastValues(ctx context.Context, b *tgbot.Bot, update *models.Update) {
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

// TODO
// triggerSentinelCallback is callback handler for exact sensor triggering query from inline keyboard.
func (s *SentinelBot) triggerSentinelCallback(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	_, err := s.bot.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{ //nolint:exhaustruct
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}

	_, err = s.bot.SendMessage(ctx, &tgbot.SendMessageParams{ //nolint:exhaustruct
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "",
	})
	if err != nil {
		s.errorHandler(ctx, update.CallbackQuery.Message.Message.Chat.ID, err)
	}
}

// TODO
// triggerSentinel is handler for /last_values command.
func (s *SentinelBot) triggerSentinel(ctx context.Context, b *tgbot.Bot, update *models.Update) {
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

func sensorIDFromCallback(data string) (string, error) {
	_, id, found := strings.Cut(data, queryPrefix)
	if !found {
		return "", errors.New("cannot separate id")
	}

	return id, nil
}
