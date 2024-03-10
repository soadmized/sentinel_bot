package api

import (
	"log"
	"sentinel_bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Sentinel struct {
	bot          *tgbotapi.BotAPI
	allowedUsers []int64
}

func New(conf config.Config) (*Sentinel, error) {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		return nil, errors.Wrap(err, "get bot api")
	}

	bot.Debug = conf.Debug
	bot.Self.SupportsInlineQueries = true

	users := []int64{conf.UserID} // temporary

	return &Sentinel{
		bot:          bot,
		allowedUsers: users,
	}, nil
}

func (s *Sentinel) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)
	s.handleUpdates(updates)
}

func (s *Sentinel) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		log.Print(update)
		received := update.Message
		msg := tgbotapi.NewMessage(received.Chat.ID, "")
		//query := update.CallbackQuery

		if !s.checkAuth(received) {
			msg.Text = "You are not authorized to use sentinel"
			s.sendMessage(msg)

			continue
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = s.inlineKeyboard()
		default:
			msg.Text = "I don't know that command"
		}

		s.sendMessage(msg)
	}

}

func (s *Sentinel) checkAuth(msg *tgbotapi.Message) bool {
	for _, id := range s.allowedUsers {
		if msg.From.ID == id {
			return true
		}
	}

	return false
}

func (s *Sentinel) checkIfMessage(msg *tgbotapi.Message) bool {
	if msg == nil || !msg.IsCommand() { // ignore any non-Message updates and non-command Messages
		return false
	}

	return true
}

func (s *Sentinel) inlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скажи привет", "hi"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скажи пока", "buy"),
		),
	)
}

func (s *Sentinel) sendMessage(msg tgbotapi.Chattable) {
	if _, err := s.bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func (s *Sentinel) handleCallbacks() {

}

func (s *Sentinel) handleCommands() {

}
