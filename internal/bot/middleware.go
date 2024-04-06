package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Middleware struct {
	AllowedUsers []int64
}

func (m Middleware) checkAuth(next tgbot.HandlerFunc) tgbot.HandlerFunc {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		var id int64

		if update.Message != nil {
			id = update.Message.From.ID
		} else {
			id = update.CallbackQuery.From.ID
		}

		for _, user := range m.AllowedUsers {
			if user == id {
				next(ctx, b, update)
			}
		}
	}
}
