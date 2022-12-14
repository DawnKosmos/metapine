// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: 05_ticker_manager.sql

package gen

import (
	"context"
	"time"
)

const createTickerManager = `-- name: CreateTickerManager :exec
INSERT INTO ticker_manager(index_id, resolution, st, et)
VALUES ($1, $2, $3, $4)
`

type CreateTickerManagerParams struct {
	IndexID    int32
	Resolution int32
	St         time.Time
	Et         time.Time
}

func (q *Queries) CreateTickerManager(ctx context.Context, arg CreateTickerManagerParams) error {
	_, err := q.db.Exec(ctx, createTickerManager,
		arg.IndexID,
		arg.Resolution,
		arg.St,
		arg.Et,
	)
	return err
}

const readTickerManager = `-- name: ReadTickerManager :one
SELECT st, et
FROM ticker_manager
WHERE index_id = $1
  and resolution = $2
`

type ReadTickerManagerParams struct {
	IndexID    int32
	Resolution int32
}

type ReadTickerManagerRow struct {
	St time.Time
	Et time.Time
}

func (q *Queries) ReadTickerManager(ctx context.Context, arg ReadTickerManagerParams) (ReadTickerManagerRow, error) {
	row := q.db.QueryRow(ctx, readTickerManager, arg.IndexID, arg.Resolution)
	var i ReadTickerManagerRow
	err := row.Scan(&i.St, &i.Et)
	return i, err
}

const updateTickerManager = `-- name: UpdateTickerManager :exec
UPDATE ticker_manager
SET st = $1, et = $2
WHERE index_id = $3
  and resolution = $4
`

type UpdateTickerManagerParams struct {
	St         time.Time
	Et         time.Time
	IndexID    int32
	Resolution int32
}

func (q *Queries) UpdateTickerManager(ctx context.Context, arg UpdateTickerManagerParams) error {
	_, err := q.db.Exec(ctx, updateTickerManager,
		arg.St,
		arg.Et,
		arg.IndexID,
		arg.Resolution,
	)
	return err
}
