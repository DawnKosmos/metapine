-- name: CreateTickerManager :exec
INSERT INTO ticker_manager(index_id, resolution, st, et)
VALUES ($1, $2, $3, $4);

-- name: ReadTickerManager :one
SELECT st, et
FROM ticker_manager
WHERE index_id = $1
  and resolution = $2;

-- name: UpdateTickerManager :exec
UPDATE ticker_manager
SET st = $1, et = $2
WHERE index_id = $3
  and resolution = $4;