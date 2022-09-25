package exchange

import (
	"time"
)

const HOUR int = 3600

var t2020 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
var t2021 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var t2022 = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var t2023 = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

type LiveTrading interface {
	Val(index int) Candle
	CandleProvider
}

type CandleProvider interface {
	OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]Candle, error)
	Name() string
}
