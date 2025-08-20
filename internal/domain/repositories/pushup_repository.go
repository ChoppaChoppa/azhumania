package repositories

import (
	"azhumania/internal/domain/models"
	"context"
	"time"
)

// PushupRepository определяет интерфейс для работы с отжиманиями
type PushupRepository interface {
	// GetTodaySession получает сессию отжиманий за сегодня
	GetTodaySession(ctx context.Context, userID int64) (*models.PushupSession, error)

	// SaveSession сохраняет сессию отжиманий
	SaveSession(ctx context.Context, session *models.PushupSession) error

	// GetSessionsByDateRange получает сессии в диапазоне дат
	GetSessionsByDateRange(ctx context.Context, userID int64, from, to time.Time) ([]*models.PushupSession, error)

	// GetWeeklyStats получает статистику за неделю
	GetWeeklyStats(ctx context.Context, userID int64) (*models.WeeklyStats, error)

	// GetMonthlyStats получает статистику за месяц
	GetMonthlyStats(ctx context.Context, userID int64) (*models.MonthlyStats, error)
}
