package repositories

import (
	"azhumania/internal/domain/models"
	"context"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// GetByTelegramID получает пользователя по Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)

	// GetByID получает пользователя по ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) error

	// Update обновляет пользователя
	Update(ctx context.Context, user *models.User) error

	// Exists проверяет существование пользователя
	Exists(ctx context.Context, telegramID int64) (bool, error)
}
