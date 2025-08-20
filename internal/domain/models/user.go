package models

import (
	"azhumania/internal/domain/errors"
	"time"
)

// User представляет пользователя в домене
type User struct {
	ID         int64
	Phone      string
	NickName   string
	TelegramID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser создает нового пользователя
func NewUser(phone, nickname string, telegramID int64) *User {
	now := time.Now()
	return &User{
		Phone:      phone,
		NickName:   nickname,
		TelegramID: telegramID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsValid проверяет валидность пользователя
func (u *User) IsValid() error {
	if u.Phone == "" {
		return errors.ErrInvalidPhone
	}
	if u.NickName == "" {
		return errors.ErrInvalidNickname
	}
	if u.TelegramID <= 0 {
		return errors.ErrInvalidTelegramID
	}
	return nil
}

// UpdateNickname обновляет никнейм пользователя
func (u *User) UpdateNickname(nickname string) {
	u.NickName = nickname
	u.UpdatedAt = time.Now()
}
