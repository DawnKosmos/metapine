package psql

import (
	"context"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"time"
)

const ohclvQueue = `-- name: GetOHCLV :many
SELECT starttime, open, high, close, low, volume
FROM %s
WHERE index_id = $1
  AND resolution = $2
  AND starttime > $3
  AND starttime < $4
`

type ohclvQueueParams struct {
	Exchange   string
	IndexId    int64
	Resolution int64
	StartTime  time.Time
	EndTime    time.Time
}

func (d *DB) ohclv(ctx context.Context, args ohclvQueueParams) ([]exchange.Candle, error) {
	qq := fmt.Sprintf(ohclvQueue, getDbName(args.Exchange, args.Resolution))
	rows, err := d.q.Query(ctx, qq, args.IndexId, args.Resolution, args.StartTime, args.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []exchange.Candle
	for rows.Next() {
		var i exchange.Candle
		if err := rows.Scan(
			&i.StartTime,
			&i.Open,
			&i.High,
			&i.Close,
			&i.Low,
			&i.Volume,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func getDbName(exchange string, resolution int64) string {
	resName := "high"
	if resolution < 3600 {
		resName = "low"
	}
	return exchange + resName
}

func (d *DB) WriteOHCLV(ctx context.Context, tableName string, indexId int64, res int64, args []exchange.Candle) (int64, error) {
	if len(args) > 10 {
		return d.copyFromOHCLV(ctx, tableName, indexId, res, args)
	}

}

// CopyFrom
func (d *DB) copyFromOHCLV(ctx context.Context, tableName string, indexId int64, res int64, args []exchange.Candle) (int64, error) {
	return d.q.CopyFrom(ctx, []string{tableName}, []string{"index_id", "resolution", "starttime", "open", "high", "close", "low", "volume"}, &iteratorForWriteOHCLV{
		indexId:              indexId,
		rows:                 args,
		res:                  res,
		skippedFirstNextCall: false,
	})
}

func checkResolution(res int64) int64 {
	fn := exchange.GenerateResolutionFunc(86400*7, 86400, 14400,
		3600, 900, 60, 15)
	return fn(res)
}

// CopyFrom
type iteratorForWriteOHCLV struct {
	indexId              int64
	rows                 []exchange.Candle
	res                  int64
	skippedFirstNextCall bool
}

func (r *iteratorForWriteOHCLV) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForWriteOHCLV) Values() ([]interface{}, error) {
	return []interface{}{
		r.indexId,
		r.res,
		r.rows[0].StartTime,
		r.rows[0].Open,
		r.rows[0].High,
		r.rows[0].Close,
		r.rows[0].Low,
		r.rows[0].Volume,
	}, nil
}

func (r iteratorForWriteOHCLV) Err() error {
	return nil
}

//WriteOHCLV

const writeOhclvQueue = `-- name: GetOHCLV :exec
INSERT INTO %s (index_id, resolution, starttime, open, high, close, low, volume)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
`

func (d *DB) writeOHCLV(ctx context.Context, tableName string, indexId int64, res int64, args []exchange.Candle) (int64, error) {
	qq := fmt.Sprintf(writeOhclvQueue, tableName)
	var c int64
	for _, v := range args {
		_, err := d.q.Exec(ctx, qq, indexId, res, v.StartTime, v.Open, v.High, v.Close, v.Low)
		if err != nil {
			fmt.Println(err)
		}
		c++
	}
	return c, nil
}
