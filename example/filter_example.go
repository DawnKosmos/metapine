package main

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/iterator"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
	"sort"
	"time"
)
import . "github.com/DawnKosmos/metapine/backend/series/ta" //!! . importation means that no  nametag is needed for this package

func FilterExample(ch *OHCLV, p backtest.Parameter) {
	open := Open(ch)
	oc := OC2(ch)
	volume := Volume(ch)
	close := Close(ch)

	TE := tradeexecution.NewScaledLimit(0, 1, 5).Distribution(distribution.Normal).Size(1) // 1 means 100%

	strat := backtest.NewSimple(ch, TE, p)

	//Timestap to measure how long the permutation took
	tNow := time.Now()
	//Add Indicators
	strat.AddIndicator(close, Rsi(close, 14))

	rsiAbove50 := backtest.SmallerAs(1, 30)

	fn := func(src Series, maFunc MaFunc, fast, slow int) (buy, sell Condition) {
		return SolApeIter(src, volume, maFunc, fast, slow)
	}
	iter := IteratorExample{fn: fn}
	it := iterator.New(&iter)
	it.RegisterInt(0, 5, 12, 1)
	it.RegisterInt(1, 4, 14, 2)
	it.RegisterFunctions(0, Sma, WrappedRsi, Ema)
	it.RegisterSeries(0, oc, open, HL2(ch), OHCL4(ch))
	it.Run(strat)

	strat.Filter("Rsi<50", rsiAbove50)

	for _, v := range strat.Results {
		v.CalculatePNL()
	}
	//Sort the results regarding TotalPNL
	sort.Sort(backtest.BackTestStrategies(strat.Results))
	ts := time.Since(tNow) // Timepassed between tNow and after the Permutation
	// Print the results
	for _, v := range strat.Results[len(strat.Results)-25:] {
		fmt.Println(v.Name, fmt.Sprintf("Trades:%d Winrate:%.2f Pnl: %.2f", len(v.Trades()), v.Winrate, v.TotalPnl))
	}
	fmt.Println("Permutation took", ts)
}
