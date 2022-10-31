package main

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange/psql"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
	"os"
	"time"
)

/*
TODO
Seperate PSQL and DB logic
CreateInterface for differentBacktestingOperations change parameters to string
Add Enviroments to Iterators
Think of filtering

Visual representation
*/

func main() {

	psql.SetPSQL("localhost", "postgres", "metapine", "admin", 5432)

	exch, _ := psql.New("ftx")

	ch := ta.NewOHCLV(exch, "BTC-PERP", time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 10, 29, 0, 0, 0, 0, time.UTC), 3600*12)
	if len(ch.Data()) == 0 {
		os.Exit(1)
	}
	open := ta.Open(ch)
	buy, sell := ta.Ribbon(open, open, ta.Sma, 36, 42)

	strat := backtest.NewStrategy(ch)

	for i := 2; i < 10; i++ {
		te := tradeexecution.NewScaledLimit(0, float64(i), 10).Distribution(distribution.Exponential).Size(1.5)
		bb := strat.CreateStrategy("Ribbon", buy, sell, te, backtest.BTParameter{
			Modus:      backtest.ALL,
			Pyramiding: 1,
			Fee: &backtest.FeeInfo{
				Maker:    -0.00005,
				Taker:    0.0005,
				Slippage: 1,
			},
			Balance:  10000,
			SizeType: size.Dollar,
			PnlGraph: false,
		})

		var lPnl, sPnl float64
		var totalVolume float64
		for _, v := range bb.Trades() {
			if v.Side {
				lPnl += v.RealisedPNL()
				//fmt.Println(v.Side, v.EntrySignalTime.Format("02/Jan/06"), v.CloseSignalTime.Format("02/Jan/06"), v.AvgBuy, v.AvgSell, v.UsdVolume/2, v.RealisedPNL())
			} else {
				//fmt.Println(v.Side, v.EntrySignalTime.Format("02/Jan/06"), v.CloseSignalTime.Format("02/Jan/06"), v.AvgSell, v.AvgBuy, v.UsdVolume/2, v.RealisedPNL())

				sPnl += v.RealisedPNL()
			}
			totalVolume += v.UsdVolume / 2
		}
		fmt.Println(fmt.Sprintf("%d, %d \t %.2f \t %2.f \t%2.f", i, len(bb.Trades()), lPnl, sPnl, totalVolume))
	}
}

func solape(oc2 ta.Series, volume ta.Series, len1, len2 int) (ta.Condition, ta.Condition) {
	outR := ta.Sma(ta.Roc(oc2, len1), 2)
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
