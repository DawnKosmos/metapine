-- name: CreateMinuteManager :exec
INSERT INTO
    minute_manager(index_id, tableName, dataArr)
VALUES ($1,$2,$3);

-- name: ReadMinuteManager :one
SELECT tableName, dataArr FROM minute_manager
WHERE index_id = $1;

-- name: UpdateMinuteManager :exec
UPDATE minute_manager
SET dataArr = $1 WHERE index_id = $2;