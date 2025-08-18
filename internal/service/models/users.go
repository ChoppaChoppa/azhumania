package models

import "azhumania/internal/repository/models"

type User struct {
	ID         int64
	Phone      string
	NickName   string
	TelegramID int64
}

func (u User) NewFromRepo(r models.User) User {
	return User{
		ID:       r.ID,
		Phone:    r.Phone,
		NickName: r.NickName,
	}
}

func (u User) ToRepo() models.User {
	return models.User{
		ID:       u.ID,
		Phone:    u.Phone,
		NickName: u.NickName,
	}
}
