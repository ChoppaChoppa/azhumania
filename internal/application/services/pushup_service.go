package services

import (
	"azhumania/internal/domain/models"
	"azhumania/internal/domain/repositories"
	"context"
	"time"
)

// PushupService предоставляет бизнес-логику для работы с отжиманиями
type PushupService struct {
	pushupRepo repositories.PushupRepository
}

// NewPushupService создает новый экземпляр PushupService
func NewPushupService(pushupRepo repositories.PushupRepository) *PushupService {
	return &PushupService{
		pushupRepo: pushupRepo,
	}
}

// AddPushupApproach добавляет новый подход отжиманий
func (s *PushupService) AddPushupApproach(ctx context.Context, userID int64, count int) (*models.PushupSession, error) {
	// Получаем или создаем сессию за сегодня
	session, err := s.getOrCreateTodaySession(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Добавляем подход
	if err := session.AddApproach(count); err != nil {
		return nil, err
	}

	// Сохраняем сессию
	if err := s.pushupRepo.SaveSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetTodayStats получает статистику за сегодня
func (s *PushupService) GetTodayStats(ctx context.Context, userID int64) (*models.PushupSession, error) {
	session, err := s.pushupRepo.GetTodaySession(ctx, userID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		// Создаем пустую сессию для сегодня
		today := time.Now().Truncate(24 * time.Hour)
		session = models.NewPushupSession(userID, today)
	}

	return session, nil
}

// GetWeeklyStats получает статистику за неделю
func (s *PushupService) GetWeeklyStats(ctx context.Context, userID int64) (*models.WeeklyStats, error) {
	return s.pushupRepo.GetWeeklyStats(ctx, userID)
}

// GetMonthlyStats получает статистику за месяц
func (s *PushupService) GetMonthlyStats(ctx context.Context, userID int64) (*models.MonthlyStats, error) {
	return s.pushupRepo.GetMonthlyStats(ctx, userID)
}

// getOrCreateTodaySession получает или создает сессию за сегодня
func (s *PushupService) getOrCreateTodaySession(ctx context.Context, userID int64) (*models.PushupSession, error) {
	session, err := s.pushupRepo.GetTodaySession(ctx, userID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		// Создаем новую сессию за сегодня
		today := time.Now().Truncate(24 * time.Hour)
		session = models.NewPushupSession(userID, today)
	}

	return session, nil
}
