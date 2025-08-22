package main

import (
	"azhumania/internal/bot/telegram"
	"azhumania/internal/service"
	"os"

	"github.com/rs/zerolog"
)

const (
	apiKey = "7887768155:AAHrGzYl8Qic0mygjmzFTfvsonY8NIuS9dg"

	psql_dsn       = "host=192.168.0.14 port=5431 user=azhumania password=Asdflkjh12 dbname=azhumania"
	redis_host     = "192.168.0.14:55000"
	redis_username = "default"
	redis_password = "redispw"
	redis_db       = 0
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	svc, err := service.New(psql_dsn, redis_host, redis_username, redis_password, redis_db, &logger)
	if err != nil {
		panic(err)
	}

	tg_bot := telegram.New(apiKey, svc)

	tg_bot.Listen()
}
