package services

import (
	"azhumania/internal/domain/errors"
	"azhumania/internal/domain/models"
	"azhumania/internal/domain/repositories"
	"context"
	"database/sql"
)

// UserService предоставляет бизнес-логику для работы с пользователями
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetOrCreateUser получает пользователя или создает нового
func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, phone, nickname string) (*models.User, error) {
	// Проверяем существование пользователя
	exists, err := s.userRepo.Exists(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if exists {
		// Получаем существующего пользователя
		user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// Создаем нового пользователя
	user := models.NewUser(phone, nickname, telegramID)
	if err := user.IsValid(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser получает пользователя по Telegram ID
func (s *UserService) GetUser(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserNickname обновляет никнейм пользователя
func (s *UserService) UpdateUserNickname(ctx context.Context, telegramID int64, nickname string) error {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	user.UpdateNickname(nickname)
	return s.userRepo.Update(ctx, user)
}
