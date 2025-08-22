package handlers

import (
	"azhumania/internal/application/services"
	"azhumania/internal/domain/models"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// CommandHandler обрабатывает команды пользователей
type CommandHandler struct {
	userService   *services.UserService
	pushupService *services.PushupService
	logger        *zerolog.Logger
}

// NewCommandHandler создает новый обработчик команд
func NewCommandHandler(userService *services.UserService, pushupService *services.PushupService, logger *zerolog.Logger) *CommandHandler {
	return &CommandHandler{
		userService:   userService,
		pushupService: pushupService,
		logger:        logger,
	}
}

// HandleStart обрабатывает команду /start
func (h *CommandHandler) HandleStart(ctx context.Context, user *models.User) (string, interface{}) {
	message := fmt.Sprintf(`Привет, %s! 👋

Я помогу тебе отслеживать твои отжимания.

📝 Как использовать:
• Просто отправляй количество отжиманий в каждом подходе
• Например: "15", "20", "10"

📊 Команды:
/start - это сообщение
/stats - статистика за неделю
/help - помощь

💡 Советы:
• Делайте перерывы между подходами
• Постепенно увеличивайте нагрузку
• Регулярность важнее количества

Удачи в тренировках! 💪`, user.NickName)

	// Создаем клавиатуру с кнопками
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

	return message, keyboard
}

// HandleHelp обрабатывает команду /help
func (h *CommandHandler) HandleHelp(ctx context.Context, user *models.User) (string, interface{}) {
	message := `🤖 Помощь по использованию бота:

📝 Как использовать:
• Просто отправляй количество отжиманий в каждом подходе
• Например: "15", "20", "10"

📊 Команды:
/start - приветствие и инструкции
/stats - статистика за неделю
/help - это сообщение

💡 Советы:
• Делайте перерывы между подходами
• Постепенно увеличивайте нагрузку
• Регулярность важнее количества

Удачи в тренировках! 💪`

	return message, nil
}

// HandleStats обрабатывает команду /stats
func (h *CommandHandler) HandleStats(ctx context.Context, user *models.User) (string, interface{}) {
	weeklyStats, err := h.pushupService.GetWeeklyStats(ctx, user.ID)
	if err != nil {
		h.logger.Error().Err(err).Int64("userID", user.ID).Msg("failed to get weekly stats")
		return "Ошибка при получении статистики. Попробуйте позже.", nil
	}

	if weeklyStats.TotalCount == 0 {
		return "📈 У вас пока нет статистики за неделю.\n\nНачните тренировки, отправляя количество отжиманий!", nil
	}

	response := "📈 Статистика за неделю:\n\n"
	response += fmt.Sprintf("Всего отжиманий: %d\n", weeklyStats.TotalCount)
	response += fmt.Sprintf("Дней тренировок: %d\n", weeklyStats.TrainingDays)
	response += fmt.Sprintf("Среднее в день: %.1f\n", weeklyStats.AveragePerDay)
	response += fmt.Sprintf("Лучший день: %d отжиманий\n", weeklyStats.BestDay)

	// Добавляем мотивацию
	if weeklyStats.TotalCount > 200 {
		response += "\n🔥 Отличная неделя! Вы на правильном пути!"
	} else if weeklyStats.TotalCount > 100 {
		response += "\n💪 Хорошая работа! Можете больше!"
	} else {
		response += "\n👍 Начинаем! Каждый день важен!"
	}

	return response, nil
}

// HandleUnknownCommand обрабатывает неизвестные команды
func (h *CommandHandler) HandleUnknownCommand(ctx context.Context, command string) (string, interface{}) {
	message := fmt.Sprintf(`Неизвестная команда: %s

Отправь мне количество отжиманий (например: 15) или используй команду /help`, command)
	return message, nil
}
