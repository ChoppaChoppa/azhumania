package handlers

import (
	"azhumania/internal/application/services"
	"azhumania/internal/domain/errors"
	"azhumania/internal/domain/models"
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// MessageHandler обрабатывает входящие сообщения
type MessageHandler struct {
	userService    *services.UserService
	pushupService  *services.PushupService
	commandHandler *CommandHandler
	logger         *zerolog.Logger
}

// NewMessageHandler создает новый обработчик сообщений
func NewMessageHandler(
	userService *services.UserService,
	pushupService *services.PushupService,
	commandHandler *CommandHandler,
	logger *zerolog.Logger,
) *MessageHandler {
	return &MessageHandler{
		userService:    userService,
		pushupService:  pushupService,
		commandHandler: commandHandler,
		logger:         logger,
	}
}

// Handle обрабатывает входящее сообщение
func (h *MessageHandler) Handle(ctx context.Context, msg *tgbotapi.Message) (string, interface{}) {
	if msg == nil {
		return "Ошибка: пустое сообщение", nil
	}

	// Получаем или создаем пользователя
	user, err := h.getOrCreateUser(ctx, msg)
	if err != nil {
		fmt.Println("@@@@@@@@@@@@@@@@@@error getOrCreateUser", err)
		h.logger.Error().Err(err).Int64("telegramID", msg.From.ID).Msg("failed to get or create user")
		return "Произошла ошибка. Попробуйте позже.", nil
	}

	// Проверяем, является ли сообщение командой или кнопкой
	if h.isCommand(msg.Text) || h.isButton(msg.Text) {
		return h.handleCommand(ctx, msg.Text, user)
	}

	// Обрабатываем количество отжиманий
	return h.handlePushupCount(ctx, msg.Text, user)
}

// getOrCreateUser получает или создает пользователя
func (h *MessageHandler) getOrCreateUser(ctx context.Context, msg *tgbotapi.Message) (*models.User, error) {
	phone := msg.From.UserName
	if phone == "" {
		phone = fmt.Sprintf("user_%d", msg.From.ID)
	}

	nickname := msg.From.FirstName
	if nickname == "" {
		nickname = "Пользователь"
	}

	return h.userService.GetOrCreateUser(ctx, msg.From.ID, phone, nickname)
}

// isCommand проверяет, является ли сообщение командой
func (h *MessageHandler) isCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}

// isButton проверяет, является ли сообщение нажатием кнопки
func (h *MessageHandler) isButton(text string) bool {
	buttons := []string{
		"📊 Статистика",
		"❓ Помощь",
		"🏠 Главное меню",
	}

	for _, button := range buttons {
		if text == button {
			return true
		}
	}
	return false
}

// handleCommand обрабатывает команды и кнопки
func (h *MessageHandler) handleCommand(ctx context.Context, command string, user *models.User) (string, interface{}) {
	switch command {
	case "/start", "🏠 Главное меню":
		return h.commandHandler.HandleStart(ctx, user)
	case "/help", "❓ Помощь":
		return h.commandHandler.HandleHelp(ctx, user)
	case "/stats", "📊 Статистика":
		return h.commandHandler.HandleStats(ctx, user)
	default:
		return h.commandHandler.HandleUnknownCommand(ctx, command)
	}
}

// handlePushupCount обрабатывает количество отжиманий
func (h *MessageHandler) handlePushupCount(ctx context.Context, text string, user *models.User) (string, interface{}) {
	// Парсим количество отжиманий
	count, err := strconv.Atoi(text)
	if err != nil {
		return "Пожалуйста, отправьте число отжиманий (например: 15) или используйте команду /help", nil
	}

	// Добавляем подход отжиманий
	session, err := h.pushupService.AddPushupApproach(ctx, user.ID, count)
	if err != nil {
		h.logger.Error().Err(err).Int64("userID", user.ID).Int("count", count).Msg("failed to add pushup approach")

		switch err {
		case errors.ErrInvalidPushupCount:
			return "Количество отжиманий должно быть больше 0", nil
		case errors.ErrPushupCountTooHigh:
			return "Количество отжиманий не может быть больше 1000 за раз", nil
		default:
			return "Произошла ошибка при сохранении данных. Попробуйте позже.", nil
		}
	}

	// Формируем ответ
	return h.formatPushupResponse(session, count), nil
}

// formatPushupResponse форматирует ответ с результатами отжиманий
func (h *MessageHandler) formatPushupResponse(session *models.PushupSession, lastCount int) string {
	response := fmt.Sprintf("✅ Сохранено %d отжиманий!\n\n", lastCount)
	response += "📊 Статистика за сегодня:\n"
	response += fmt.Sprintf("   • Всего отжиманий: %d\n", session.GetTotalCount())
	response += fmt.Sprintf("   • Подходов: %d\n", session.GetApproachCount())
	response += fmt.Sprintf("   • Среднее за подход: %.1f\n", session.GetAveragePerApproach())

	// Добавляем мотивацию
	response += h.getMotivationMessage(session.GetTotalCount())

	return response
}

// getMotivationMessage возвращает мотивационное сообщение
func (h *MessageHandler) getMotivationMessage(totalCount int) string {
	switch {
	case totalCount > 100:
		return "\n🔥 Невероятно! Вы настоящий чемпион!"
	case totalCount > 50:
		return "\n🔥 Отличная работа! Продолжайте в том же духе!"
	case totalCount > 20:
		return "\n💪 Хороший результат! Можете больше!"
	default:
		return "\n👍 Начинаем! Каждый подход важен!"
	}
}
