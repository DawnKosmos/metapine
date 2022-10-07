package psql

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"strings"
	"time"
)

type Exchange struct {
	exchange string
}

func (e *Exchange) OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]exchange.Candle, error) {
	index, err := p.qq.GetIndexIdByName(ctx, fmt.Sprintf("%s:%s", e.exchange, strings.ToLower(ticker)))
	if err != nil || index == 0 {
		if "index" == e.exchange {
			return nil, errors.New("Index does not exist:" + ticker)
		}

		var ee exchange.CandleProvider
		switch e.exchange {
		case "ftx":
			ee = ftx.New()
		case "binance":
		case "bybit":
		default:
			return nil, errors.New("exchange not supported")
		}
		ee.OHCLV(ticker, resolution, start, end)

	}

	/*
		if err != nil || index == 0{
			if "index" == e.exchange{
				return err
			}
			get exchange ...
			//checke ob ticker existiert, exchange.OHCLV(ticker, resolution, st,end)
			CreateTicker, CreateIndex,
		}
	*/
}

func getOHCLV(name string, ticker string, resolution int64, start time.Time, end time.Time) ([]exchange.Candle, error) {
	var ee exchange.CandleProvider
	switch name {
	case "ftx":
		ee = ftx.New()
	case "binance":
	case "bybit":
	default:
		return nil, errors.New("exchange not supported")
	}
	return ee.OHCLV(ticker, resolution, start, end)
}

func (e *Exchange) Name() string {
	return e.exchange
}

type Index struct {
	Resolution int64
	Tickers    []Ticker
}

func New(name string) (*Exchange, error) {
	name = strings.ToLower(name)
	switch name {
	case "ftx", "deribit", "index":
		return &Exchange{exchange: name}, nil
	case "bybit", "bitmex", "coinbase", "phemex":
		return nil, errors.New("Exchange not implemented")
	default:
		return nil, errors.New("exchange does not exist")
	}
}
