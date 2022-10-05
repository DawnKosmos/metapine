-- name: CreateIndex :one
INSERT INTO index (name)
VALUES ($1)
RETURNING index_id;

-- name: CreateTicker :one
INSERT INTO ticker (exchange, ticker)
VALUES ($1, $2)
RETURNING ticker_id;

-- name: CreateTickerIndex :exec
INSERT INTO ticker_index (ticker_id, index_id, weight, excludevolume)
VALUES ($1, $2, $3, $4);

-- name: ReturnIndex :many
SELECT ticker.exchange, ticker.ticker, ticker_index.weight, ticker_index.excludevolume
FROM ticker_index
         JOIN ticker ON ticker.ticker_id = ticker_index.ticker_id
         JOIN index ON index.index_id = ticker_index.index_id
WHERE index.index_id = $1;

-- name: DeleteTickerIndex :exec
DELETE
FROM index
WHERE index_id = $1;