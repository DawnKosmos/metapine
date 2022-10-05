// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: 02_ticker_index.sql

package gen

import (
	"context"
	"database/sql"
)

const createIndex = `-- name: CreateIndex :one
INSERT INTO index (name)
VALUES ($1)
RETURNING index_id
`

func (q *Queries) CreateIndex(ctx context.Context, name sql.NullString) (int64, error) {
	row := q.db.QueryRowContext(ctx, createIndex, name)
	var index_id int64
	err := row.Scan(&index_id)
	return index_id, err
}

const createTicker = `-- name: CreateTicker :one
INSERT INTO ticker (exchange, ticker)
VALUES ($1, $2)
RETURNING ticker_id
`

type CreateTickerParams struct {
	Exchange Exchanges
	Ticker   string
}

func (q *Queries) CreateTicker(ctx context.Context, arg CreateTickerParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createTicker, arg.Exchange, arg.Ticker)
	var ticker_id int32
	err := row.Scan(&ticker_id)
	return ticker_id, err
}

const createTickerIndex = `-- name: CreateTickerIndex :exec
INSERT INTO ticker_index (ticker_id, index_id, weight, excludevolume)
VALUES ($1, $2, $3, $4)
`

type CreateTickerIndexParams struct {
	TickerID      sql.NullInt32
	IndexID       sql.NullInt64
	Weight        int32
	Excludevolume bool
}

func (q *Queries) CreateTickerIndex(ctx context.Context, arg CreateTickerIndexParams) error {
	_, err := q.db.ExecContext(ctx, createTickerIndex,
		arg.TickerID,
		arg.IndexID,
		arg.Weight,
		arg.Excludevolume,
	)
	return err
}

const deleteTickerIndex = `-- name: DeleteTickerIndex :exec
DELETE
FROM index
WHERE index_id = $1
`

func (q *Queries) DeleteTickerIndex(ctx context.Context, indexID int64) error {
	_, err := q.db.ExecContext(ctx, deleteTickerIndex, indexID)
	return err
}

const returnIndex = `-- name: ReturnIndex :many
SELECT ticker.exchange, ticker.ticker, ticker_index.weight, ticker_index.excludevolume
FROM ticker_index
         JOIN ticker ON ticker.ticker_id = ticker_index.ticker_id
         JOIN index ON index.index_id = ticker_index.index_id
WHERE index.index_id = $1
`

type ReturnIndexRow struct {
	Exchange      Exchanges
	Ticker        string
	Weight        int32
	Excludevolume bool
}

func (q *Queries) ReturnIndex(ctx context.Context, indexID int64) ([]ReturnIndexRow, error) {
	rows, err := q.db.QueryContext(ctx, returnIndex, indexID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReturnIndexRow
	for rows.Next() {
		var i ReturnIndexRow
		if err := rows.Scan(
			&i.Exchange,
			&i.Ticker,
			&i.Weight,
			&i.Excludevolume,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}