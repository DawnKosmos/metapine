package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
)

type TEInfo struct {
	Name             string
	Info             string
	CandlePnlSupport bool
}

/*
The TradeExecution Interface describes how a Trade is executed
The subfolder tradeexecution is showing some examples
From market orders to scaled limit orders it is implemented
*/
type TradeExecution interface {
	CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat, sizeInUsd float64, fee Fee, pnlgraph bool) (*Trade, error) //TradeExecution defines the strategy and gets as input an array from trade start to end
	GetInfo() TEInfo
}

func DefaultFeeInfo() *Fee {
	return new(Fee)
}
