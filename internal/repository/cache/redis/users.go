package redis

import (
	"azhumania/internal/repository/models"
	"context"
)

func (r *repository) GetUser(ctx context.Context, userID int64) (models.User, error) {

}

func (r *repository) SetUser(ctx context.Context, user models.User) error {

}
