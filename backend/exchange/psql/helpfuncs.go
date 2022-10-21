package psql

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"github.com/DawnKosmos/metapine/helper/formula"
	"strconv"
	"time"
)

var fnResolutionFunc = exchange.GenerateResolutionFunc(86400, 3600*3, 7200,
	3600, 900, 60, 15)

func checkResolution(res int64) int64 {
	return fnResolutionFunc(res)
}

func getDbName(exchange string, resolution int64) string {
	return exchange
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

func indexName(exchange string, ticker string) string {
	return fmt.Sprintf("%s:%s", exchange, ticker)
}

func minutesTable(id int32) string {
	return "minutes_" + strconv.Itoa(int(id))
}

type yearmonth struct {
	y int
	m int
}

func getMonthsAndYears(st, et time.Time) (o []yearmonth) {
	for st.Before(et) {
		o = append(o, yearmonth{st.Year(), int(st.Month())})
		st = st.AddDate(0, 1, 0)
	}
	return o
}

func getTimeStampsOfMonth(y int, m int) (time.Time, time.Time) {
	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC), time.Date(y, time.Month(m+1), 1, 0, -1, 0, 0, time.UTC)
}

func last[T any](a []T) T {
	return formula.Last(a)
}
