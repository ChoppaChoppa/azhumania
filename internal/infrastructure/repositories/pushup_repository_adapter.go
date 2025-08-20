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

// PushupRepositoryAdapter адаптирует существующие репозитории к доменному интерфейсу
type PushupRepositoryAdapter struct {
	db     psql.IDatabase
	cache  redis.ICache
	logger zerolog.Logger
}

// NewPushupRepositoryAdapter создает новый адаптер репозитория отжиманий
func NewPushupRepositoryAdapter(db psql.IDatabase, cache redis.ICache, logger zerolog.Logger) repositories.PushupRepository {
	return &PushupRepositoryAdapter{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// GetTodaySession получает сессию отжиманий за сегодня
func (r *PushupRepositoryAdapter) GetTodaySession(ctx context.Context, userID int64) (*domainModels.PushupSession, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// Сначала пробуем из кэша
	repoAzhumania, err := r.cache.GetAzhumania(ctx, userID)
	if err == nil && repoAzhumania.UserID != 0 && repoAzhumania.Date.Equal(today) {
		return r.convertToDomainSession(repoAzhumania), nil
	}

	// Если нет в кэше, получаем из БД
	repoAzhumania, err = r.db.GetAzhumania(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет данных за сегодня
		}
		r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to get azhumania from database")
		return nil, err
	}

	// Проверяем, что данные за сегодня
	if !repoAzhumania.Date.Equal(today) {
		return nil, nil // Данные не за сегодня
	}

	// Сохраняем в кэш асинхронно
	go func() {
		ctx := context.Background()
		if err := r.cache.SetAzhumania(ctx, repoAzhumania); err != nil {
			r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to cache azhumania")
		}
	}()

	return r.convertToDomainSession(repoAzhumania), nil
}

// SaveSession сохраняет сессию отжиманий
func (r *PushupRepositoryAdapter) SaveSession(ctx context.Context, session *domainModels.PushupSession) error {
	repoAzhumania := r.convertToRepoAzhumania(session)

	// Сохраняем в БД
	err := r.db.AddAzhumania(ctx, repoAzhumania)
	if err != nil {
		r.logger.Error().Err(err).Int64("userID", session.UserID).Msg("failed to save azhumania to database")
		return err
	}

	// Сохраняем в кэш асинхронно
	go func() {
		ctx := context.Background()
		if err := r.cache.SetAzhumania(ctx, repoAzhumania); err != nil {
			r.logger.Error().Err(err).Int64("userID", session.UserID).Msg("failed to cache azhumania")
		}
	}()

	return nil
}

// GetSessionsByDateRange получает сессии в диапазоне дат
func (r *PushupRepositoryAdapter) GetSessionsByDateRange(ctx context.Context, userID int64, from, to time.Time) ([]*domainModels.PushupSession, error) {
	// TODO: Реализовать получение сессий по диапазону дат
	// Пока возвращаем пустой список
	r.logger.Warn().Msg("GetSessionsByDateRange not implemented")
	return []*domainModels.PushupSession{}, nil
}

// GetWeeklyStats получает статистику за неделю
func (r *PushupRepositoryAdapter) GetWeeklyStats(ctx context.Context, userID int64) (*domainModels.WeeklyStats, error) {
	// TODO: Реализовать получение статистики за неделю
	// Пока возвращаем заглушку
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))

	stats := domainModels.NewWeeklyStats(userID, weekStart)
	stats.TotalCount = 0
	stats.TrainingDays = 0
	stats.AveragePerDay = 0
	stats.BestDay = 0

	return stats, nil
}

// GetMonthlyStats получает статистику за месяц
func (r *PushupRepositoryAdapter) GetMonthlyStats(ctx context.Context, userID int64) (*domainModels.MonthlyStats, error) {
	// TODO: Реализовать получение статистики за месяц
	// Пока возвращаем заглушку
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	stats := domainModels.NewMonthlyStats(userID, monthStart)
	stats.TotalCount = 0
	stats.TrainingDays = 0
	stats.AveragePerDay = 0
	stats.BestDay = 0
	stats.Streak = 0

	return stats, nil
}

// convertToDomainSession конвертирует репозиторную модель в доменную сессию
func (r *PushupRepositoryAdapter) convertToDomainSession(repoAzhumania repoModels.Azhumania) *domainModels.PushupSession {
	session := domainModels.NewPushupSession(repoAzhumania.UserID, repoAzhumania.Date)
	session.ID = repoAzhumania.UserID // Используем UserID как ID сессии

	// Создаем один подход с общим количеством отжиманий
	// TODO: Изменить структуру БД для хранения отдельных подходов
	approach := domainModels.PushupApproach{
		SessionID: session.ID,
		Count:     repoAzhumania.Count,
		CreatedAt: repoAzhumania.Date,
	}
	session.Approaches = append(session.Approaches, approach)

	return session
}

// convertToRepoAzhumania конвертирует доменную сессию в репозиторную модель
func (r *PushupRepositoryAdapter) convertToRepoAzhumania(session *domainModels.PushupSession) repoModels.Azhumania {
	return repoModels.Azhumania{
		UserID: session.UserID,
		Date:   session.Date,
		Count:  session.GetTotalCount(),
	}
}
