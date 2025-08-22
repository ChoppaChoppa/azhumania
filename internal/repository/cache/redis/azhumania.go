package redis

import (
	"azhumania/internal/repository/models"
	"context"
	"encoding/json"
)

func (r *repository) GetAzhumania(ctx context.Context, userID int64) ([]models.Azhumania, error) {
	azhumania := models.Azhumania{UserID: userID}

	data, err := r.cache.Get(ctx, azhumania.CacheKey()).Result()
	if err != nil {
		return nil, err
	}

	var azhumaniaList []models.Azhumania
	if err := json.Unmarshal([]byte(data), &azhumaniaList); err != nil {
		// Попробуем десериализовать как одну запись (для обратной совместимости)
		var singleAzhumania models.Azhumania
		if err := json.Unmarshal([]byte(data), &singleAzhumania); err != nil {
			return nil, err
		}
		return []models.Azhumania{singleAzhumania}, nil
	}

	return azhumaniaList, nil
}

func (r *repository) SetAzhumania(ctx context.Context, azhumania models.Azhumania) error {
	// Получаем существующие записи
	existingRecords, err := r.GetAzhumania(ctx, azhumania.UserID)
	if err != nil {
		// Если записи не найдены, создаем новый слайс
		existingRecords = []models.Azhumania{}
	}

	// Добавляем новую запись
	existingRecords = append(existingRecords, azhumania)

	// Сохраняем обновленный слайс
	data, err := json.Marshal(existingRecords)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, azhumania.CacheKey(), data, 0).Err()
}
