package services

import (
	"azhumania/internal/domain/errors"
	"azhumania/internal/domain/models"
	"azhumania/internal/domain/repositories"
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

// UserService предоставляет бизнес-логику для работы с пользователями
type UserService struct {
	userRepo repositories.UserRepository
	logger   *zerolog.Logger
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo repositories.UserRepository, logger *zerolog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetOrCreateUser получает пользователя или создает нового
func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, phone, nickname string) (*models.User, error) {
	// Проверяем существование пользователя
	exists, err := s.userRepo.Exists(ctx, telegramID)
	if err != nil {
		s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to check user existence")
		return nil, err
	}

	if exists {
		// Получаем существующего пользователя
		user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
		if err != nil {
			s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to get user by telegramID")
			return nil, err
		}
		return user, nil
	}

	// Создаем нового пользователя
	user := models.NewUser(phone, nickname, telegramID)
	if err := user.IsValid(); err != nil {
		s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to validate user")
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to create user")
		return nil, err
	}

	return user, nil
}

// GetUser получает пользователя по Telegram ID
func (s *UserService) GetUser(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("user not found")
			return nil, errors.ErrUserNotFound
		}
		s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to get user by telegramID")
		return nil, err
	}
	return user, nil
}

// UpdateUserNickname обновляет никнейм пользователя
func (s *UserService) UpdateUserNickname(ctx context.Context, telegramID int64, nickname string) error {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to get user by telegramID")
		return err
	}

	user.UpdateNickname(nickname)
	return s.userRepo.Update(ctx, user)
}
