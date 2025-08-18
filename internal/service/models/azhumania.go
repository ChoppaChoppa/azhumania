package models

import (
	"azhumania/internal/repository/models"
	"time"
)

type Azhumania struct {
	UserID int64
	Date   time.Time
	Count  int
}

func (a Azhumania) NewFromRepo(r models.Azhumania) Azhumania {
	return Azhumania{
		UserID: r.UserID,
		Date:   r.Date,
		Count:  r.Count,
	}
}

func (a Azhumania) ToRepo() models.Azhumania {
	return models.Azhumania{
		UserID: a.UserID,
		Date:   a.Date,
		Count:  a.Count,
	}
}
