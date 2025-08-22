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
	logger *zerolog.Logger
}

// NewPushupRepositoryAdapter создает новый адаптер репозитория отжиманий
func NewPushupRepositoryAdapter(db psql.IDatabase, cache redis.ICache, logger *zerolog.Logger) repositories.PushupRepository {
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
	repoAzhumaniaList, err := r.cache.GetAzhumania(ctx, userID)
	if err == nil && len(repoAzhumaniaList) > 0 {
		// Ищем запись за сегодня
		for _, azhumania := range repoAzhumaniaList {
			if azhumania.Date.Equal(today) {
				return r.convertToDomainSession(azhumania), nil
			}
		}
	}

	// Если нет в кэше, получаем из БД
	repoAzhumaniaList, err = r.db.GetAzhumania(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет данных за сегодня
		}
		r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to get azhumania from database")
		return nil, err
	}

	// Ищем запись за сегодня
	for _, azhumania := range repoAzhumaniaList {
		if azhumania.Date.Equal(today) {
			// Сохраняем в кэш асинхронно
			go func() {
				ctx := context.Background()
				if err := r.cache.SetAzhumania(ctx, azhumania); err != nil {
					r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to cache azhumania")
				}
			}()

			return r.convertToDomainSession(azhumania), nil
		}
	}

	return nil, nil // Нет данных за сегодня
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
	// Получаем все записи пользователя
	repoAzhumaniaList, err := r.db.GetAzhumania(ctx, userID)
	if err != nil {
		r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to get azhumania for date range")
		return nil, err
	}

	var sessions []*domainModels.PushupSession

	// Создаем карту для группировки записей по дням
	dailySessions := make(map[time.Time]*domainModels.PushupSession)

	// Обрабатываем записи в диапазоне дат
	for _, azhumania := range repoAzhumaniaList {
		// Проверяем, что запись входит в диапазон
		if azhumania.Date.After(from) && azhumania.Date.Before(to) {
			dayKey := azhumania.Date.Truncate(24 * time.Hour)

			// Создаем или получаем сессию для этого дня
			session, exists := dailySessions[dayKey]
			if !exists {
				session = domainModels.NewPushupSession(userID, dayKey)
				session.ID = userID
				dailySessions[dayKey] = session
			}

			// Добавляем подход
			approach := domainModels.PushupApproach{
				SessionID: session.ID,
				Count:     azhumania.Count,
				CreatedAt: azhumania.Date,
			}
			session.Approaches = append(session.Approaches, approach)
		}
	}

	// Преобразуем карту в слайс
	for _, session := range dailySessions {
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetWeeklyStats получает статистику за неделю
func (r *PushupRepositoryAdapter) GetWeeklyStats(ctx context.Context, userID int64) (*domainModels.WeeklyStats, error) {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)

	// Получаем все записи пользователя
	repoAzhumaniaList, err := r.db.GetAzhumania(ctx, userID)
	if err != nil {
		r.logger.Error().Err(err).Int64("userID", userID).Msg("failed to get azhumania for weekly stats")
		return nil, err
	}

	stats := domainModels.NewWeeklyStats(userID, weekStart)

	// Создаем карту для подсчета отжиманий по дням
	dailyCounts := make(map[time.Time]int)

	// Обрабатываем записи за неделю
	for _, azhumania := range repoAzhumaniaList {
		// Проверяем, что запись входит в неделю
		if azhumania.Date.After(weekStart) && azhumania.Date.Before(weekEnd) {
			dayKey := azhumania.Date.Truncate(24 * time.Hour)
			dailyCounts[dayKey] += azhumania.Count
		}
	}

	// Вычисляем статистику
	stats.TrainingDays = len(dailyCounts)
	stats.TotalCount = 0
	stats.BestDay = 0

	for _, count := range dailyCounts {
		stats.TotalCount += count
		if count > stats.BestDay {
			stats.BestDay = count
		}
	}

	if stats.TrainingDays > 0 {
		stats.AveragePerDay = float64(stats.TotalCount) / float64(stats.TrainingDays)
	}

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
