package psql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

p *pgxpool.Pool

type DB struct {
	q *pgxpool.Pool
}

var ctx = context.Background()

func (d *DB) Ping() error {
	return d.q.Ping(ctx)
}

func SetPSQL()

//Being able to register indexes
