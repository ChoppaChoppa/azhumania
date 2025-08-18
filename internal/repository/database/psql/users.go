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

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo GetUser.QueryxContext")
		return
	}

	defer rows.Close()

	if err = rows.StructScan(&user); err != nil {
		r.logger.Error().Err(err).Msg("error repo GetUser.StructScan")
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

	if err = r.db.SelectContext(ctx, user.ID, query, args...); err != nil {
		r.logger.Error().Err(err).Msg("error repo AddUser.ExecContext")
		return 0, err
	}

	return user.ID, nil
}
