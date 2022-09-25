package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/helper/formula"
)

const LONG = true
const SHORT = false

type BacktestParameters struct {
	FeeFlat    float64
	Pyramiding int
	Slippage   float64
	MakerFee   float64
	TakerFee   float64
	Size       Size
}

type TradeExecution interface {
	CreateTrade(Side bool, ch []exchange.Candle) (Trade, error) //TradeExecution defines the strategy and gets as input an array from trade start to end
	SetOHCLV(o ta.Chart)                                        //Just set the OHCLV, usually not used, but if needed to create a trade. it can be used
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
		FeeFlat:    0,
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

func NewStrategy(e ta.Chart, parameters BacktestParameters) *SimpleStrategy {

}

func NewFast(chart ta.Chart, buy, sell ta.Condition, te TradeExecution, Pyraminding int) SimpleStrategy {
	ch, l, s := chart.Data(), buy.Data(), sell.Data()
	te.SetOHCLV(chart)
	parameters := DefaultParameters()
	if Pyraminding > 1 {
		parameters.Pyramiding = Pyraminding
	}

	sl, _ := formula.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl+1:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]

	var trades []*Trade
	var tempOrderLong, tempOrderShort []exchange.Candle

	for i, c := range ch[:len(ch)-1] {
		if l[i] {
			for j := 0; j < min(len(tempOrderShort), parameters.Pyramiding); j++ {

			}
		}
	}

}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
