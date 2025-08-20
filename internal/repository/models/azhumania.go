package models

import (
	"fmt"
	"time"
)

type Azhumania struct {
	UserID int64     `json:"user_id" db:"user_id"`
	Date   time.Time `json:"date" db:"date"`
	Count  int       `json:"count" db:"count"`
}

func (a Azhumania) CacheKey() string {
	return fmt.Sprintf("azhumania:%d", a.UserID)
}
