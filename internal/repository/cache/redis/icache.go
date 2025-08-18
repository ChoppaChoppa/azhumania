package redis

import (
	"azhumania/internal/repository/models"
	"context"
)

type ICache interface {
	IUsersCache
	IAzhumaniaCache
}

type IUsersCache interface {
	GetUser(context.Context, int64) (models.User, error)
	SetUser(context.Context, models.User) error
}

type IAzhumaniaCache interface {
	GetAzhumania(context.Context, int64) (models.Azhumania, error)
	SetAzhumania(context.Context, models.Azhumania) error
}
