package backtest

import (
	"errors"
	"metapine/backend/exchange"
	"time"
)

type SafeFloat struct {
	Safe  bool
	Value float64
}

type Fill struct {
	Side  bool
	Price float64
	Size  float64
	Time  time.Time
}

type Trade struct {
	Side      bool
	AvgEntry  float64
	AvgExit   float64
	EntryTime time.Time
	CloseTime time.Time

	Fills     []Fill
	Size      float64
	Pnl       []float64
	Indicator []SafeFloat
}

type STrade struct {
	Side  bool
	Entry float64
	Exit  float64
}

func CreateSimpleTrade(side bool, entry, exit exchange.Candle) (STrade, error) {
	if entry.StartTime.Unix() >= exit.StartTime.Unix() {
		return STrade{}, errors.New("Error Candle")
	}

	return STrade{
		Side:  side,
		Entry: entry.Open,
		Exit:  exit.Open,
	}, nil
}

func (t *STrade) Pnl(fee float64) float64 {
	var x float64
	if t.Side {
		x = (t.Exit - t.Entry) / t.Entry
	} else {
		x = -1 * (t.Exit - t.Entry) / t.Entry
	}
	return x - (fee * 0.01)
}
