package models

import "time"

// WeeklyStats представляет статистику за неделю
type WeeklyStats struct {
	UserID        int64
	WeekStart     time.Time
	WeekEnd       time.Time
	TotalCount    int
	TrainingDays  int
	AveragePerDay float64
	BestDay       int
	BestDayDate   time.Time
}

// MonthlyStats представляет статистику за месяц
type MonthlyStats struct {
	UserID        int64
	Month         time.Time
	TotalCount    int
	TrainingDays  int
	AveragePerDay float64
	BestDay       int
	BestDayDate   time.Time
	Streak        int // Текущая серия дней
}

// NewWeeklyStats создает новую статистику за неделю
func NewWeeklyStats(userID int64, weekStart time.Time) *WeeklyStats {
	weekEnd := weekStart.AddDate(0, 0, 6)
	return &WeeklyStats{
		UserID:    userID,
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
	}
}

// NewMonthlyStats создает новую статистику за месяц
func NewMonthlyStats(userID int64, month time.Time) *MonthlyStats {
	return &MonthlyStats{
		UserID: userID,
		Month:  month,
	}
}
