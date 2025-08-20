package models

import (
	"azhumania/internal/domain/errors"
	"time"
)

// PushupSession представляет сессию отжиманий за день
type PushupSession struct {
	ID         int64
	UserID     int64
	Date       time.Time
	Approaches []PushupApproach
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// PushupApproach представляет один подход отжиманий
type PushupApproach struct {
	ID        int64
	SessionID int64
	Count     int
	CreatedAt time.Time
}

// NewPushupSession создает новую сессию отжиманий
func NewPushupSession(userID int64, date time.Time) *PushupSession {
	now := time.Now()
	return &PushupSession{
		UserID:     userID,
		Date:       date,
		Approaches: make([]PushupApproach, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// AddApproach добавляет новый подход к сессии
func (ps *PushupSession) AddApproach(count int) error {
	if err := validatePushupCount(count); err != nil {
		return err
	}

	approach := PushupApproach{
		SessionID: ps.ID,
		Count:     count,
		CreatedAt: time.Now(),
	}

	ps.Approaches = append(ps.Approaches, approach)
	ps.UpdatedAt = time.Now()

	return nil
}

// GetTotalCount возвращает общее количество отжиманий в сессии
func (ps *PushupSession) GetTotalCount() int {
	total := 0
	for _, approach := range ps.Approaches {
		total += approach.Count
	}
	return total
}

// GetApproachCount возвращает количество подходов
func (ps *PushupSession) GetApproachCount() int {
	return len(ps.Approaches)
}

// GetAveragePerApproach возвращает среднее количество отжиманий за подход
func (ps *PushupSession) GetAveragePerApproach() float64 {
	if len(ps.Approaches) == 0 {
		return 0
	}
	return float64(ps.GetTotalCount()) / float64(len(ps.Approaches))
}

// IsToday проверяет, является ли сессия за сегодня
func (ps *PushupSession) IsToday() bool {
	today := time.Now().Truncate(24 * time.Hour)
	return ps.Date.Equal(today)
}

// validatePushupCount валидирует количество отжиманий
func validatePushupCount(count int) error {
	if count <= 0 {
		return errors.ErrInvalidPushupCount
	}
	if count > 1000 {
		return errors.ErrPushupCountTooHigh
	}
	return nil
}
