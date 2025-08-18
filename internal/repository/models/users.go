package models

type User struct {
	ID       int64  `json:"id" db:"id"`
	Phone    string `json:"phone" db:"phone"`
	NickName string `json:"nickname" db:"nickname"`
}
