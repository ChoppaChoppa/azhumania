package redis

import (
	"azhumania/internal/repository/models"
	"context"
	"encoding/json"
)

func (r *repository) GetAzhumania(ctx context.Context, userID int64) (models.Azhumania, error) {
	azhumania := models.Azhumania{UserID: userID}

	data, err := r.cache.Get(ctx, azhumania.CacheKey()).Result()
	if err != nil {
		return models.Azhumania{}, err
	}

	if err := json.Unmarshal([]byte(data), &azhumania); err != nil {
		return models.Azhumania{}, err
	}

	return azhumania, nil
}

func (r *repository) SetAzhumania(ctx context.Context, azhumania models.Azhumania) error {
	data, err := json.Marshal(azhumania)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, azhumania.CacheKey(), data, 0).Err()
}
