package psql

import (
	"azhumania/internal/repository/models"
	"context"
)

type IDatabase interface {
	IUsersDatabase
	IAzhumaniaDatabase
}

type IUsersDatabase interface {
	GetUser(context.Context, int64) (models.User, error)
	AddUser(context.Context, models.User) (int64, error)
}

type IAzhumaniaDatabase interface {
	GetAzhumania(context.Context, int64) (models.Azhumania, error)
	AddAzhumania(context.Context, models.Azhumania) error
}
