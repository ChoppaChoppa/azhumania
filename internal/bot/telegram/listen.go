package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TelegramBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		answer := t.service.Handle(update.Message)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

		if _, err := t.bot.Send(msg); err != nil {
			fmt.Println("error Send", err)
		}
	}
}
