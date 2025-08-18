package main

import (
	"azhumania/internal/bot/telegram"
	"azhumania/internal/service"
	"github.com/rs/zerolog"
)

const (
	apiKey = "7887768155:AAHrGzYl8Qic0mygjmzFTfvsonY8NIuS9dg"

	psql_dsn       = "host=localhost port=5431 user=azhumania password=Asdflkjh12 dbname=azhumania"
	redis_host     = "localhost:55000"
	redis_username = "default"
	redis_password = "redispw"
	redis_db       = 0
)

func main() {

	svc, err := service.New(psql_dsn, redis_host, redis_username, redis_password, redis_db, zerolog.Logger{})
	if err != nil {
		panic(err)
	}

	tg_bot := telegram.New(apiKey, svc)

	tg_bot.Listen()
}
