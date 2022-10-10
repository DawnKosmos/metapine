package psql

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
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
