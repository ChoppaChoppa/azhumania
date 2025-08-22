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
		// Обрабатываем сообщения
		if update.Message != nil {
			answer, keyboard := t.service.Handle(update.Message)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

			// Добавляем клавиатуру, если она есть
			if keyboard != nil {
				switch k := keyboard.(type) {
				case tgbotapi.ReplyKeyboardMarkup:
					msg.ReplyMarkup = k
				case tgbotapi.ReplyKeyboardRemove:
					msg.ReplyMarkup = k
				}
			}

			if _, err := t.bot.Send(msg); err != nil {
				fmt.Println("error Send", err)
			}
		}

		// Обрабатываем callback-запросы от inline кнопок
		if update.CallbackQuery != nil {
			t.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

// handleCallbackQuery обрабатывает callback-запросы от inline кнопок
func (t *TelegramBot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	// Отвечаем на callback, чтобы убрать "часики" у кнопки
	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	if _, err := t.bot.Request(callbackResponse); err != nil {
		fmt.Println("error callback response", err)
	}

	// Создаем фейковое сообщение для обработки
	fakeMessage := &tgbotapi.Message{
		From: callback.From,
		Chat: callback.Message.Chat,
		Text: callback.Data,
	}

	answer, keyboard := t.service.Handle(fakeMessage)

	// Отправляем ответ в чат
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, answer)

	// Добавляем клавиатуру, если она есть
	if keyboard != nil {
		switch k := keyboard.(type) {
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = k
		case tgbotapi.ReplyKeyboardRemove:
			msg.ReplyMarkup = k
		}
	}

	if _, err := t.bot.Send(msg); err != nil {
		fmt.Println("error Send callback response", err)
	}
}
