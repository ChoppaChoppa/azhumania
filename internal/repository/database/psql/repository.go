package psql

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type repository struct {
	db *sqlx.DB

	builder squirrel.StatementBuilderType
	logger  zerolog.Logger
}

func New(dsn string, logger zerolog.Logger) (IDatabase, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &repository{
		db:      db,
		builder: sq,
		logger:  logger,
	}, nil
}
