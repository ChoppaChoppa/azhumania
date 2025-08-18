package service

import (
	"azhumania/internal/repository/cache/redis"
	"azhumania/internal/repository/database/psql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

var _ IService = &service{}

type IService interface {
	Handle(msg *tgbotapi.Message) string
}

type service struct {
	db    psql.IDatabase
	cache redis.ICache

	logger zerolog.Logger
}

func New(psql_dsn, redis_host, redis_username, redis_password string, redis_db int, logger zerolog.Logger) (IService, error) {
	db, err := psql.New(psql_dsn, logger)
	if err != nil {
		return nil, err
	}

	cache := redis.New(redis_host, redis_username, redis_password, redis_db, logger)
	if cache == nil {
		return nil, errors.New("redis connection error")
	}

	return &service{
		db:    db,
		cache: cache,

		logger: logger,
	}, nil
}
