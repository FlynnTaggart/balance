package databases

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxDB struct {
	*pgxpool.Pool
	Logger pgx.Logger
}

func NewPgxDB(pool *pgxpool.Pool, logger pgx.Logger) *PgxDB {
	return &PgxDB{
		Pool:   pool,
		Logger: logger,
	}
}
