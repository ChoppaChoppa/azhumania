package repositories

import (
	domainModels "azhumania/internal/domain/models"
	"azhumania/internal/domain/repositories"
	"azhumania/internal/repository/cache/redis"
	"azhumania/internal/repository/database/psql"
	repoModels "azhumania/internal/repository/models"
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog"
)

// UserRepositoryAdapter адаптирует существующие репозитории к доменному интерфейсу
type UserRepositoryAdapter struct {
	db     psql.IDatabase
	cache  redis.ICache
	logger *zerolog.Logger
}

// NewUserRepositoryAdapter создает новый адаптер репозитория пользователей
func NewUserRepositoryAdapter(db psql.IDatabase, cache redis.ICache, logger *zerolog.Logger) repositories.UserRepository {
	return &UserRepositoryAdapter{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// GetByTelegramID получает пользователя по Telegram ID
func (r *UserRepositoryAdapter) GetByTelegramID(ctx context.Context, telegramID int64) (*domainModels.User, error) {
	// Сначала пробуем из кэша
	repoUser, err := r.cache.GetUser(ctx, telegramID)
	if err == nil && repoUser.ID != 0 {
		return r.convertToDomainUser(repoUser), nil
	}

	// Если нет в кэше, получаем из БД
	repoUser, err = r.db.GetUser(ctx, telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		r.logger.Error().Err(err).Int64("telegramID", telegramID).Msg("failed to get user from database")
		return nil, err
	}

	// Сохраняем в кэш асинхронно
	go func() {
		ctx := context.Background()
		if err := r.cache.SetUser(ctx, repoUser); err != nil {
			r.logger.Error().Err(err).Int64("userID", repoUser.ID).Msg("failed to cache user")
		}
	}()

	return r.convertToDomainUser(repoUser), nil
}

// GetByID получает пользователя по ID
func (r *UserRepositoryAdapter) GetByID(ctx context.Context, id int64) (*domainModels.User, error) {
	// Сначала пробуем из кэша
	repoUser, err := r.cache.GetUser(ctx, id)
	if err == nil && repoUser.ID != 0 {
		return r.convertToDomainUser(repoUser), nil
	}

	// Если нет в кэше, получаем из БД
	repoUser, err = r.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		r.logger.Error().Err(err).Int64("userID", id).Msg("failed to get user from database")
		return nil, err
	}

	// Сохраняем в кэш асинхронно
	go func() {
		ctx := context.Background()
		if err := r.cache.SetUser(ctx, repoUser); err != nil {
			r.logger.Error().Err(err).Int64("userID", repoUser.ID).Msg("failed to cache user")
		}
	}()

	return r.convertToDomainUser(repoUser), nil
}

// Create создает нового пользователя
func (r *UserRepositoryAdapter) Create(ctx context.Context, user *domainModels.User) error {
	repoUser := r.convertToRepoUser(user)

	id, err := r.db.AddUser(ctx, repoUser)
	if err != nil {
		r.logger.Error().Err(err).Interface("user", user).Msg("failed to create user in database")
		return err
	}

	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Сохраняем в кэш асинхронно
	go func() {
		ctx := context.Background()
		repoUser.ID = id
		if err := r.cache.SetUser(ctx, repoUser); err != nil {
			r.logger.Error().Err(err).Int64("userID", id).Msg("failed to cache new user")
		}
	}()

	return nil
}

// Update обновляет пользователя
func (r *UserRepositoryAdapter) Update(ctx context.Context, user *domainModels.User) error {
	user.UpdatedAt = time.Now()
	repoUser := r.convertToRepoUser(user)

	// Обновляем в БД (здесь нужно будет добавить метод Update в интерфейс БД)
	// TODO: Добавить метод Update в psql.IDatabase
	r.logger.Warn().Msg("Update method not implemented in database interface")

	// Обновляем в кэше
	go func() {
		ctx := context.Background()
		if err := r.cache.SetUser(ctx, repoUser); err != nil {
			r.logger.Error().Err(err).Int64("userID", user.ID).Msg("failed to update user in cache")
		}
	}()

	return nil
}

// Exists проверяет существование пользователя
func (r *UserRepositoryAdapter) Exists(ctx context.Context, telegramID int64) (bool, error) {
	// Проверяем в кэше
	repoUser, err := r.cache.GetUser(ctx, telegramID)
	if err == nil && repoUser.ID != 0 {
		return true, nil
	}

	// Проверяем в БД
	repoUser, err = r.db.GetUser(ctx, telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return repoUser.ID != 0, nil
}

// convertToDomainUser конвертирует репозиторную модель в доменную
func (r *UserRepositoryAdapter) convertToDomainUser(repoUser repoModels.User) *domainModels.User {
	return &domainModels.User{
		ID:         repoUser.ID,
		Phone:      repoUser.Phone,
		NickName:   repoUser.NickName,
		TelegramID: repoUser.ID, // Предполагаем, что ID в БД это TelegramID
		CreatedAt:  time.Now(),  // TODO: Добавить поля CreatedAt/UpdatedAt в репозиторную модель
		UpdatedAt:  time.Now(),
	}
}

// convertToRepoUser конвертирует доменную модель в репозиторную
func (r *UserRepositoryAdapter) convertToRepoUser(domainUser *domainModels.User) repoModels.User {
	return repoModels.User{
		ID:       domainUser.TelegramID, // Используем TelegramID как ID в БД
		Phone:    domainUser.Phone,
		NickName: domainUser.NickName,
	}
}
