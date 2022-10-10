package psql

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"strings"
	"time"
)

type Exchange struct {
	ee       exchange.CandleProvider
	exchange string
}

func (e *Exchange) OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]exchange.Candle, error) {
	ticker = strings.ToLower(ticker)
	indexId, err := p.qq.GetIndexIdByName(ctx, indexName(e.Name(), ticker))
	if e.exchange == "index" {
		if err != nil {
			return nil, errors.New("Index does not exist:" + ticker)
		}

	}

	if err != nil || indexId == 0 {
		if "index" == e.exchange {
			return nil, errors.New("Index does not exist:" + ticker)
		}
		//initOhclv downloads the ticker(if exists) and saves it in the provided database
		return initOhclv(e.exchange, ticker, resolution, start, end) //Creates Ticker, Index, and Downloads OHCLV
	}
	rows, err := p.qq.ReturnIndex(ctx, indexId)
	for _, v := range rows {
		fmt.Println(v.Name)
	}
	return nil, nil
}

func (e *Exchange) Name() string {
	return e.exchange
}

// Creates a New CandleProvider. Data gets Saved in an SQL DB
func New(name string) (exchange.CandleProvider, error) {
	name = strings.ToLower(name)
	switch name {
	case "ftx", "deribit", "bybit", "index":
		return &Exchange{exchange: name}, nil
	case "bitmex", "coinbase", "phemex":
		return nil, errors.New("Exchange not yet implemented")
	default:
		return nil, errors.New("exchange does not exist")
	}
}
