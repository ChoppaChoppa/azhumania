package models

import "time"

type Azhumania struct {
	UserID int64     `json:"user_id" db:"user_id"`
	Date   time.Time `json:"date" db:"date"`
	Count  int       `json:"count" db:"count"`
}
