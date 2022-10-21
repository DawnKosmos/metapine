package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
)

const LONG = true
const SHORT = false

type Mode int

const (
	OnlySHORT Mode = -1
	ALL            = 0
	OnlyLONG       = 1
)

type BacktestParameters struct {
	Pyramiding int
	Slippage   float64
	MakerFee   float64
	TakerFee   float64
	Size       Size
}

type order struct {
	Side bool
	TradeExecution
}

type SimpleStrategy struct {
	e          ta.Chart
	parameters *BacktestParameters
}

func DefaultParameters() BacktestParameters {
	return BacktestParameters{
		Pyramiding: 1,
		Slippage:   0,
		MakerFee:   0,
		TakerFee:   0,
		Size: Size{
			Type: AccountSize,
			Val:  100,
		},
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// A better representation of the PNL overtime
type CandlePNL struct {
	Open  float64
	High  float64
	Close float64
	Low   float64
}

func PNLCalcCandle(avgEntry float64, ch exchange.Candle) (c CandlePNL) {
	//TODO
	return
}
