package psql

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"time"
)

type Ticker struct {
	Exchange      gen.Exchanges
	Ticker        string
	Weight        int32
	ExcludeVolume bool
}

func registerNewTicker(name gen.Exchanges, ticker string) (int32, error) {
	var ee exchange.CandleProvider
	switch name {
	case gen.ExchangesFtx:
		ee = ftx.New()
	case gen.ExchangesBinance:
	}
	tNow := time.Now()
	ch, _ := ee.OHCLV(ticker, 3600, tNow.Add(-7200*time.Second), tNow)
	if len(ch) == 0 {
		return 0, errors.New("ticker not found")
	}

	return p.qq.CreateTicker(ctx, gen.CreateTickerParams{Exchange: name, Ticker: ticker})
}

func (t Ticker) ohclv() {

}
