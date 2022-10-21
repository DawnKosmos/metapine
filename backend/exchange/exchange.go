package exchange

import (
	"time"
)

const HOUR int = 3600

var T2020 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
var T2021 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var T2022 = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
var T2023 = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

type LiveTa interface {
	Val(index int) Candle
	CandleProvider
}

type LivaTrading interface {
	LiveTa
	PrivateExchange
}

type PrivateExchange interface {
	SetOrder(input SetOrderInput) (Order, error)
	OpenOrder(ticker ...string) ([]Order, error)
	OpenTriggerORder(ticker ...string) ([]TriggerOrder, error)
	Cancel(Side bool, Ticker string) (int, error)
	CancelOrderById(id int64) error
	CancelTrigger(Side bool, ticker Ticker) (int, error)
	CancelTriggerOrderById(id int64) error
	FundingPayments(ticker string, starttime int64, endtime int64) ([]FundingPayments, error)
	Fills(starttime int64, endtime int64) ([]Fill, error)
	OpenPosition() ([]Position, error)
}

type Index struct {
	Tickers []Ticker
}

type Ticker struct {
	Exchange      string
	Ticker        string
	Weight        uint
	ExcludeVolume bool
}

type CandleProvider interface {
	OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]Candle, error)
	Name() string
}

/*
	NewIndex(in Index) error
	Init(exchange string, ticker string, resolution int64, ch []Candle) error
	AddCandle(exchange string, ticker string, resolution int64, ch []Candle) error
	UpdateCandle(exchange string, ticker string, resolution int64, ch []Candle) error
*/

/*
Ticker sollen abspeichbar sein. In den Einstellungen kann man wählen wie Abgespeichert werden soll. eine Json Datei soll standard sein.
PSQL Datenbank ist die bessere wahl, aber erfordert eine Initialisierung
Es sollte Möglichkeiten geben Indizes zu definieren

*/
