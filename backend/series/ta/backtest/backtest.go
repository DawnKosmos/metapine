package backtest

import (
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

type FeeInfo struct {
	Maker    float64
	Taker    float64
	Slippage float64
}

type BacktestParameters struct {
	Pyramiding int
	Fee        *FeeInfo
	Size       *Size
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
		Size: &Size{
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
