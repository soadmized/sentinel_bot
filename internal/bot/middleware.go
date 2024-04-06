package bot

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Middleware struct {
	AllowedUsers []int64
}

func (m Middleware) checkAuth(next tgbot.HandlerFunc) tgbot.HandlerFunc {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		for _, user := range m.AllowedUsers {
			if user == update.Message.From.ID {
				next(ctx, b, update)
			}
		}

		log.Print(update.Message.From.ID)
	}
}
