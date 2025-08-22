package psql

import (
	"azhumania/internal/repository/models"
	"context"
	"github.com/Masterminds/squirrel"
)

func (r *repository) GetUser(ctx context.Context, userID int64) (user models.User, err error) {
	query, args, err := r.builder.
		Select(
			"id",
			"phone",
			"nickname",
		).
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo GetUser.ToSql")
		return
	}

	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&user)
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo GetUser.QueryxContext")
		return
	}

	return
}

func (r *repository) AddUser(ctx context.Context, user models.User) (int64, error) {
	query, args, err := r.builder.
		Insert("users").
		Columns(
			"phone",
			"nickname",
		).
		Values(
			user.Phone,
			user.NickName,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo AddUser.ToSql")
		return 0, err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo AddUser.ExecContext")
		return 0, err
	}

	return user.ID, nil
}
