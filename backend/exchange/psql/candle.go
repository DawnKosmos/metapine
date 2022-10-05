package psql

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/jackc/pgx/v5"
)

type iterateForCandle struct {
	IndexId, Resolution  int64
	rows                 []exchange.Candle
	skippedFirstNextCall bool
}

func (r *iterateForCandle) Next() bool {
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

func (r *iterateForCandle) Values() ([]interface{}, error) {
	return []interface{}{
		r.IndexId,
		r.Resolution,
		r.rows[0].StartTime,
		r.rows[0].Open,
		r.rows[0].High,
		r.rows[0].Close,
		r.rows[0].Low,
		r.rows[0].Volume,
	}, nil
}

func (r *iterateForCandle) Err() error {
	return nil
}

func Init(tickerId int64, tableName string, ch []exchange.Candle) (int64, error) {
	res := ch[1].StartTime.Unix() - ch[0].StartTime.Unix()
	cc := iterateForCandle{
		IndexId:    tickerId,
		Resolution: res,
		rows:       ch,
	}
	return p.q.CopyFrom(ctx, []string{tableName}, []string{"index_id", "resolution", "starttime", "open", "high", "close", "low", "volume"}, &cc)
}

//With TX

const insertCandle = `-- name: InsertCandle :one
INSERT INTO %s
(index_id, resolution, starttime, open,high,close, low, volume)
VALUES ($1, $2,$3,$4,$5,$6,$7,$8)
RETURNING session_uuid
`

func InitTX(tickerId int64, tableName string, ch []exchange.Candle) (int64, error) {
	res := ch[1].StartTime.Unix() - ch[0].StartTime.Unix()

	tx, err := p.q.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	})
	if err != nil {
		return 0, err
	}

	ss := fmt.Sprintf(insertCandle, tableName)
	stmt, err := tx.Prepare(ctx, "insert"+tableName, ss)
	if err != nil {
		return 0, err
	}
	var count int64
	for _, v := range ch {
		_, err := tx.Exec(ctx, stmt.Name, tickerId, res, ch, v.StartTime, v.Open, v.High, v.Close, v.Low, v.Volume)
		if err != nil {
			tx.Rollback(ctx)
			return count, err
		}
		count++
	}

	return count, tx.Commit(ctx)
}
