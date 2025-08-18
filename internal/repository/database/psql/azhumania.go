package psql

import (
	"azhumania/internal/repository/models"
	"context"
	"github.com/Masterminds/squirrel"
)

func (r *repository) GetAzhumania(ctx context.Context, userID int64) (azhumania models.Azhumania, err error) {
	query, args, err := r.builder.
		Select(
			"user_id",
			"date",
			"count",
		).
		From("azhumania").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo GetAzhumania.ToSql")
		return
	}

	if err = r.db.SelectContext(ctx, &azhumania, query, args...); err != nil {
		r.logger.Error().Err(err).Msg("error repo GetAzhumania.SelectContext")
		return
	}

	return
}

func (r *repository) AddAzhumania(ctx context.Context, azhumania models.Azhumania) error {
	query, args, err := r.builder.
		Insert("azhumania").
		Columns(
			"user_id",
			"date",
			"count",
		).
		Values(
			azhumania.UserID,
			azhumania.Date,
			azhumania.Count,
		).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("error repo AddAzhumania.ToSql")
		return 0, err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		r.logger.Error().Err(err).Msg("error repo AddAzhumania.ExecContext")
		return 0, err
	}

	return nil
}
