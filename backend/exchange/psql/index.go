package psql

import (
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"strings"
)

type Index struct {
	name    string
	Tickers []Ticker
}

// Register a new Index
func NewIndex(name string, t ...Ticker) error {
	/*
		CreateIndex
		GetTicker
			Optional: RegisterNewTicker
		CreateTickerIndex
	*/
	id, err := p.qq.CreateIndex(ctx, "index_"+name)
	if err != nil {
		return err
	}
	for _, v := range t {
		tickerId, err := p.qq.GetTickerId(ctx, gen.GetTickerIdParams{
			Exchange: v.Exchange,
			Ticker:   strings.ToLower(v.Ticker),
		})
		if err != nil {
			p.qq.DeleteIndex(ctx, id)
			return err
		}
		if tickerId == 0 {
			tickerId, err = registerNewTicker(v.Exchange, v.Ticker)
			if err != nil {
				p.qq.DeleteIndex(ctx, id)
				return err
			}
		}
		if err = p.qq.CreateTickerIndex(ctx, gen.CreateTickerIndexParams{
			TickerID:      tickerId,
			IndexID:       id,
			Weight:        v.Weight,
			Excludevolume: v.ExcludeVolume,
		}); err != nil {
			p.qq.DeleteIndex(ctx, id)
			return err
		}
	}
	return nil
}
