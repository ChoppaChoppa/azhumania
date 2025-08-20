package models

import "fmt"

type User struct {
	ID       int64  `json:"id" db:"id"`
	Phone    string `json:"phone" db:"phone"`
	NickName string `json:"nickname" db:"nickname"`
}

func (u User) CacheKey() string {
	return fmt.Sprintf("user:%d", u.ID)
}
