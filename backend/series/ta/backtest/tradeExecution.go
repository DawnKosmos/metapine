package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
)

type TEInfo struct {
	Name             string
	Info             string
	CandlePnlSupport bool
}

type TradeExecution interface {
	CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat, sizeInUsd float64, fee FeeInfo, pnlgraph bool) (*Trade, error) //TradeExecution defines the strategy and gets as input an array from trade start to end
	GetInfo() TEInfo
}

func DefaultFeeInfo() *FeeInfo {
	return new(FeeInfo)
}
