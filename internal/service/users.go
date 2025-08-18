package service

import (
	"azhumania/internal/service/models"
	"context"
)

func (s *service) GetUser(ctx context.Context, userID int64) (user models.User, err error) {
	resp, err := s.cache.GetUser(ctx, userID)
	if err != nil || resp.ID == 0 {
		resp, err = s.db.GetUser(ctx, userID)
		if err != nil {
			s.logger.Error().Err(err).Msg("error svc GetUser.GetUser")
			return
		}
	}

	user.NewFromRepo(resp)

	return user, nil
}

func (s *service) CreateUser(ctx context.Context, user models.User) (int64, error) {
	id, err := s.db.AddUser(ctx, user.ToRepo())
	if err != nil {
		s.logger.Error().Err(err).Msg("error add user")
		return 0, err
	}

	go func() {
		ctx := context.Background()

		if err = s.cache.SetUser(ctx, user.ToRepo()); err != nil {
			s.logger.Error().Err(err).Msg("error set user")
		}
	}()

	return id, nil
}
