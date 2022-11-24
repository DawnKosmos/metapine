package psql

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/deribit"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"strings"
	"time"
)

type Exchange struct {
	ee       exchange.CandleProvider
	exchange string
}

// Creates a New CandleProvider. Data gets Saved in an PSQL DB
func New(name string) (*Exchange, error) {
	name = strings.ToLower(name)
	switch name {
	case "ftx":
		return &Exchange{exchange: name, ee: ftx.New()}, nil
	case "index":
	case "deribit":
		return &Exchange{exchange: name, ee: deribit.New()}, nil
	case "bitmex", "coinbase", "phemex", "bybit":
		return nil, errors.New("Exchange not yet implemented")
	default:
		return nil, errors.New("exchange does not exist")
	}
	return nil, nil
}

func (e *Exchange) OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]exchange.Candle, error) {
	ticker = strings.ToLower(ticker)
	if end.After(time.Now()) {
		end = time.Now()
	}
	if start.After(end) {
		return nil, errors.New("StartTime has to be Lower than EndTime or cant be in future")
	}

	indexId, err := p.qq.GetIndexIdByName(ctx, indexName(e.Name(), ticker))
	if e.exchange == "index" {
		if err != nil {
			return nil, errors.New("Index does not exist:" + ticker)
		}
	}

	if err != nil || indexId == 0 {
		n, _ := stringToExchanges(e.Name())
		tickerId, err := registerNewTicker(n, ticker)
		if err != nil {
			return nil, err
		}
		indexId, err := p.qq.CreateIndex(ctx, indexName(e.Name(), ticker))
		if err != nil {
			return nil, err
		}
		err = p.qq.CreateTickerIndex(ctx, gen.CreateTickerIndexParams{
			TickerID:      tickerId,
			IndexID:       indexId,
			Weight:        1,
			Excludevolume: false,
		})
		return initOhclv(indexId, e.ee, ticker, resolution, start, end)
	}

	return ohclvTicker(indexId, e.ee, ticker, resolution, start, end)
}

func (e *Exchange) Name() string {
	return e.exchange
}
