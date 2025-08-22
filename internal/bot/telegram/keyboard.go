package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetMainKeyboard возвращает основную клавиатуру с командами
func GetMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.KeyboardButton{Text: "📊 Статистика"},
				tgbotapi.KeyboardButton{Text: "❓ Помощь"},
			},
			{
				tgbotapi.KeyboardButton{Text: "🏠 Главное меню"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
		Selective:       false,
	}

	return keyboard
}

// GetInlineKeyboard возвращает inline клавиатуру для быстрых действий
func GetInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.InlineKeyboardButton{
					Text:         "📊 Статистика",
					CallbackData: &[]string{"stats"}[0],
				},
				tgbotapi.InlineKeyboardButton{
					Text:         "❓ Помощь",
					CallbackData: &[]string{"help"}[0],
				},
			},
			{
				tgbotapi.InlineKeyboardButton{
					Text:         "🏠 Главное меню",
					CallbackData: &[]string{"start"}[0],
				},
			},
		},
	}

	return keyboard
}

// RemoveKeyboard возвращает команду для скрытия клавиатуры
func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
}
