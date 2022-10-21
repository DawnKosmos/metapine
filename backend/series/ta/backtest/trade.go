package backtest

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
)

type Trade struct {
	Side            bool
	AvgEntry        float64
	AvgExit         float64
	EntrySignalTime int64
	CloseSignalTime int64

	Fills []Fill
	Size  float64
	//The PNL starts with the EntrySignalTime. Every Tick represents 1 Candle
	//This information is needed to calculate the Overall PNL of the Indicator
	Pnl       []float64
	PnlCandle []CandlePNL
	Indicator []SafeFloat
}

// STrade is used For FastBacktesting, this mode is used to iterate many parameters.
// Only the results are safed to Calculate
type STrade struct {
	Side                bool
	Entry, Exit         float64
	EntryTime, ExitTime int64
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
