-- name: ReadOHCL :many
SELECT starttime, open, high, close, low, volume
FROM ohclv
WHERE index_id = $2
  AND resolution = $3
  AND starttime > $4
  AND starttime < $1;