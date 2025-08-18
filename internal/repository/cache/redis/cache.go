package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type repository struct {
	cache *redis.Client

	logger zerolog.Logger
}

func New(addr, username, password string, db int, logger zerolog.Logger) ICache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	return &repository{
		cache:  client,
		logger: logger,
	}
}
