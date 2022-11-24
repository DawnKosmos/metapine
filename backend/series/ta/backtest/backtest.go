package backtest

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
)

const LONG = true
const SHORT = false

type Backtester interface {
	AddStrategy(buy, sell ta.Condition, parameters string)
	Split(condition string, op Filter)
	Filter(condition string, op Filter)
}

// BacktestingParameters are needed to better simulate a real enviroment
type BacktestParameters struct {
	//Modus, means OnlyShort, onlyLongs or ALL
	Modus mode.Mode
	//On Default pyramiding is one
	Pyramiding int
	//Fees are described by Maker(Market Orders) Taker(limit orders) and Slippage(market orders)
	Fee *Fee
	//Size can be choosen in USD or Account size turns
	Size *size.Size
}

type Fee struct {
	Maker    float64
	Taker    float64
	Slippage float64
}

/*
type order struct {
	Side bool
	TradeExecution
}
*/

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

// Parameters are needed to better simulate a real enviroment
type Parameter struct {
	//Modus onlyLong, onlyShort, All
	Modus mode.Mode
	//The maximum number of entries allowed in the same direction
	Pyramiding int
	//The exchange fees. Remember that 0.05% are 0.0005
	Fee *Fee
	//The Account Balance
	Balance float64
	//Size is either Dollar or AccountSize
	SizeType size.SizeBase
	//right now PnlGraph isnt implemented. It will be useful in visual representation
	PnlGraph bool
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
