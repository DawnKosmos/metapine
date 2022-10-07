-- name: GetTickerId :one
SELECT ticker_id FROM ticker
WHERE exchange = $1 and ticker = $2;


