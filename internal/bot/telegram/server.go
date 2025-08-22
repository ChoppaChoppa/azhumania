package telegram

import (
	"azhumania/internal/service"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	apiKey  string
	bot     *tgbotapi.BotAPI
	service service.IService
}

func New(key string, service service.IService) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramBot{
		apiKey:  key,
		bot:     bot,
		service: service,
	}
}
