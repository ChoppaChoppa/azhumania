package redis

import (
	"azhumania/internal/repository/models"
	"context"
	"encoding/json"
)

func (r *repository) GetUser(ctx context.Context, userID int64) (models.User, error) {
	user := models.User{ID: userID}

	data, err := r.cache.Get(ctx, user.CacheKey()).Result()
	if err != nil {
		return models.User{}, err
	}

	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *repository) SetUser(ctx context.Context, user models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, user.CacheKey(), data, 0).Err()
}
