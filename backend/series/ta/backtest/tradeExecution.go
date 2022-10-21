package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
)

type TradeExecution interface {
	CreateTrade(Side bool, ch []exchange.Candle) (Trade, error) //TradeExecution defines the strategy and gets as input an array from trade start to end
	CandlePnlSupport() bool
}
