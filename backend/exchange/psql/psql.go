package psql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	q *pgxpool.Pool
}

var ctx = context.Background()

func (d *DB) Ping() error {
	return d.q.Ping(ctx)
}

//Being able to register indexes
