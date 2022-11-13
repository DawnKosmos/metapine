package backtest

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"testing"
)

func TestExecution(t *testing.T) {
	sl := NewScaledLimit(0, 15, 10).Size(2).Distribution(distribution.Exponential)

	ff := ftx.New()

	ch := ta.NewOHCLV(ff, "FTT-perp", exchange.T2022, exchange.T2023, 3600*24)
	o, h, c, l, v := ta.ChartSources(ch)

	buy, sell := solape(o, ta.Roc, v, 4, 10)
	stoch := ta.Stoch(c, h, l, 14)
	sma := ta.Sma(o, 20)
	vwma := ta.Vwma(c, v, 15)
	strat := NewStrategy(ch).AddIndicator(stoch, sma, vwma)

	bt := strat.CreateStrategy("rsi cross over", buy, sell, sl, BTParameter{
		Modus:      ALL,
		Pyramiding: 1,
		Fee: &Fee{
			Maker:    0.00015,
			Taker:    0.0005,
			Slippage: 1,
		},
		Balance:  10000,
		SizeType: size.Account,
		PnlGraph: false,
	})

	var long, short float64
	var cc int
	for _, v := range bt.tr {
		if v.Side {
			long += v.RealisedPNL()
		} else {
			short += v.RealisedPNL()
		}
		cc++
	}

	fmt.Println(cc, long, short)
}

func solape(oc2 ta.Series, ma func(ta.Series, int) ta.Series, volume ta.Series, len1, len2 int) (ta.Condition, ta.Condition) {
	outR := ta.Sma(ma(oc2, len1), 2)
	outB1 := ta.Sma(outR, len2)
	outB2 := ta.Sma(outB1, len2)
	outB := ta.SubF(outB1, outB2, 2.0)
	cc := ta.Sub(outR, outB)
	var c1 ta.Series
	if volume == nil {
		c1 = ta.Sma(cc, 2)
	} else {
		c1 = ta.Vwma(cc, volume, 2)
	}
	c2, c3 := ta.OffS(c1, 1), ta.OffS(c1, 2)
	buy := ta.And(ta.Greater(c1, c2), ta.Smaller(c2, c3))
	sell := ta.And(ta.Smaller(c1, c2), ta.Greater(c2, c3))
	return buy, sell

}
