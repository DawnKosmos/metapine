package psql

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"time"
)

func checkResolution(res int64) int64 {
	fn := exchange.GenerateResolutionFunc(86400*7, 86400, 14400,
		3600, 900, 60, 15)
	return fn(res)
}

func getDbName(exchange string, resolution int64) string {
	resName := "high"
	if resolution < 3600 {
		resName = "low"
	}
	return exchange + resName
}

func stringToExchanges(s string) (e gen.Exchanges, err error) {
	switch s {
	case "ftx":
		e = gen.ExchangesFtx
	case "bybit":
		e = gen.ExchangesBybit
	case "deribit":
		e = gen.ExchangesDeribit
	case "binance":
		e = gen.ExchangesBinance
	case "bitmex":
		e = gen.ExchangesBitmex
	default:
		err = errors.New("invalid exchange")
	}
	return
}

func stringToCandleProvider(s string) (e exchange.CandleProvider) {
	switch s {
	case "ftx":
		e = ftx.New()
	case "binance":
	case "bybit":
	}
	return
}

type tupel struct {
	st time.Time
	et time.Time
}

func missingCandles(ch []exchange.Candle, res int64) []tupel {
	var t []tupel
	realLen := int64(len(ch))
	expectedLen := ch[realLen-1].StartTime.Unix() - ch[0].StartTime.Unix()/res
	if expectedLen > realLen {
		for i := 1; i < len(ch); i++ {
			if ch[i].StartTime.Unix()-ch[i-1].StartTime.Unix() != res {
				t = append(t, tupel{ch[i-1].StartTime, ch[i].StartTime})
			}
		}
	}
	return t
}

func fixMissingCandles(ch []exchange.Candle, res int64, ee exchange.CandleProvider)

1440*365
1500*300 = 480000*4 = 2 millionen