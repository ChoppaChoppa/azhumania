package errors

import "errors"

// Доменные ошибки
var (
	ErrInvalidPhone       = errors.New("invalid phone number")
	ErrInvalidNickname    = errors.New("invalid nickname")
	ErrInvalidTelegramID  = errors.New("invalid telegram ID")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPushupCount = errors.New("invalid pushup count")
	ErrPushupCountTooHigh = errors.New("pushup count too high")
)
