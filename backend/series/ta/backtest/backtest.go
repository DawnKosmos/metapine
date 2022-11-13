package backtest

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
)

const LONG = true
const SHORT = false

type Fee struct {
	Maker    float64
	Taker    float64
	Slippage float64
}

type BacktestParameters struct {
	Modus      mode.Mode
	Pyramiding int
	Fee        *Fee
	Size       *size.Size
}

type order struct {
	Side bool
	TradeExecution
}

type SimpleStrategy struct {
	e          ta.Chart
	parameters *BacktestParameters
}

func DefaultParameters() *BacktestParameters {
	return &BacktestParameters{
		Pyramiding: 1,
		Fee:        DefaultFeeInfo(),
		Size:       size.New(size.Account, 1),
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type Parameter struct {
	Modus      mode.Mode
	Pyramiding int
	Fee        *Fee
	Balance    float64
	SizeType   size.SizeBase
	PnlGraph   bool
}

type Backtester interface {
	AddStrategy(buy, sell ta.Condition, parameters string)
	Split(condition string, op Filter)
	Filter(condition string, op Filter)
}

/*
BTParameter, TradeExecution

NewFastBackTest -> FastBacktest
	bt.AddStrategy(buy,sell)
NewStrategy -> Backtest
	bt.CreateStrategy("name", buy, sell, TE, parameters BTParameter
NewMultiTicker -> MultiTicker
	bt.CreateResult(tickers []string, ee exchange.CandleProvider, st, et,res)

*/
