package service

import (
	"azhumania/internal/service/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

var azhumaniaByDay map[time.Time]int = map[time.Time]int{}

func (s *service) Handle(msg *tgbotapi.Message) string {
	if msg == nil {
		return ""
	}

	userID := msg.From.ID

	ctx := context.Background()

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err := s.CreateUser(ctx, models.User{
				Phone:      msg.From.UserName,
				TelegramID: msg.From.ID,
			})
			if err != nil {
				s.logger.Error().Err(err).Msg("error svc Handle.CreateUser")
				return err.Error()
			}
		}
	}

	count, err := strconv.Atoi(msg.Text)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error converting message to int: %s", err.Error()))
		return "не число"
	}

	return ""
}
