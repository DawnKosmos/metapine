-- name: ReadOHCL :many
SELECT starttime, open, high, close, low, volume
FROM ohclv
WHERE index_id = $2
  AND resolution = $3
  AND starttime > $4
  AND starttime < $1;



-- name: WriteOHCLV :copyfrom
INSERT INTO ohclv (index_id, resolution, starttime,open,high,close,low,volume)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
