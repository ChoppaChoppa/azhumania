package psql

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	db *sqlx.DB

	builder squirrel.StatementBuilderType
	logger  *zerolog.Logger
}

func New(dsn string, logger *zerolog.Logger) (IDatabase, error) {
	db, err := sqlx.Connect("pgx", dsn)
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
