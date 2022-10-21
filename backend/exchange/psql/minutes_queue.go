package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"time"
)

/*
To manage 1 min data efficient. if we import them always a whole month gets downloaded
A own table manages the name and the months used by this database
a month starts with time.Data(YEAR, time.Month, 0,0,0,0,0,time.UTC)
and ends with time.Data(YEAR, time.Month+1, 0,0,0,-1,0,time.UTC)
the current month get treaten differently

*/

const createNewTableQueue = `
CREATE TABLE IF NOT EXISTS %s
(
    starttime timestamp unique not null,
    open      float4           not null,
    high      float4           not null,
    close     float4           not null,
    low       float4           not null,
    volume    float4           not null
);
`

func createNewMinutesTable(indexId int32) error {
	s := minutesTable(indexId)
	fmt.Println(s)
	qq := fmt.Sprintf(createNewTableQueue, s)
	_, err := p.q.Exec(ctx, qq)
	if err != nil {
		return err
	}
	err = p.qq.CreateMinuteManager(ctx, gen.CreateMinuteManagerParams{
		IndexID:   indexId,
		Tablename: minutesTable(indexId),
		Dataarr: sql.NullString{
			Valid: false,
		},
	})
	return err
}

const minutesOhclvqueue = `-- name: GetOHCLV :many
SELECT starttime, open, high, close, low, volume
FROM %s
WHERE index_id = $1
  AND starttime > $3
  AND starttime < $4
`

func (d *DB) minutesOhclv(ctx context.Context, tableName string, IndexId int32, st time.Time, et time.Time) ([]exchange.Candle, error) {
	qq := fmt.Sprintf(minutesOhclvqueue, tableName)
	rows, err := d.q.Query(ctx, qq, IndexId, st, et)
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

func (d *DB) MinutesWriteOHCLV(ctx context.Context, tableName string, args []exchange.Candle) (int64, error) {
	return d.q.CopyFrom(ctx, []string{tableName}, []string{"starttime", "open", "high", "close", "low", "volume"}, &iteratorMinutesOHCLV{
		rows:                 args,
		skippedFirstNextCall: false,
	})
}

// CopyFrom
type iteratorMinutesOHCLV struct {
	rows                 []exchange.Candle
	skippedFirstNextCall bool
}

func (r *iteratorMinutesOHCLV) Next() bool {
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

func (r iteratorMinutesOHCLV) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].StartTime,
		r.rows[0].Open,
		r.rows[0].High,
		r.rows[0].Close,
		r.rows[0].Low,
		r.rows[0].Volume,
	}, nil
}

func (r iteratorMinutesOHCLV) Err() error {
	return nil
}

// ==== Last Candle
const minuteLastCandle = `-- name: GetOHCLV :one
SELECT starttime
FROM %s
ORDER BY starttime DESC LIMIT 1
`

func (d *DB) lastCandle(tb string) (st int64, err error) {
	qq := fmt.Sprintf(minuteLastCandle, tb)

	res := d.q.QueryRow(ctx, qq)
	err = res.Scan(&st)
	return st, err
}
