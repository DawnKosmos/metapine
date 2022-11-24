package main

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/iterator"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
	"math"
	"sort"
	"time"
)
import . "github.com/DawnKosmos/metapine/backend/series/ta" //!! . importation means that no  nametag is needed for this package

func FilterExample(ch *OHCLV, p backtest.Parameter) {

	open := Open(ch)
	oc := OC2(ch)
	volume := Volume(ch)
	close := Close(ch)

	TE := tradeexecution.NewScaledLimit(0, 5, 10).Distribution(distribution.Normal).Size(1) // 1 means 100%

	strat := backtest.NewSimple(ch, TE, p)

	//Timestap to measure how long the permutation took
	tNow := time.Now()
	//Add Indicators
	macd, signal, _ := MacdRelative(close, 12, 26, 9)
	strat.AddIndicator(close, Rsi(close, 14), Sma(close, 100), macd, signal)

	macdAboveBelow2 := func(sf []backtest.SafeFloat) bool {
		if sf[3].Safe {
			return math.Abs(sf[3].Value) < 2
		}
		return false
	}

	rsiAbove50 := backtest.GreaterAs(1, 50)

	fn := func(src Series, maFunc MaFunc, fast, slow int) (buy, sell Condition) {
		return SolApeIter(src, volume, maFunc, fast, slow)
	}
	strat.AddIndicator()
	iter := IteratorExample{fn: fn}
	it := iterator.New(&iter)
	it.RegisterInt(0, 4, 12, 1)
	it.RegisterInt(1, 3, 14, 2)
	it.RegisterFunctions(0, Sma, WrappedRsi, Ema)
	it.RegisterSeries(0, oc, open, HL2(ch), OHCL4(ch))
	it.Run(strat)

	strat.Split("Rsi>50", rsiAbove50)
	strat.Split("Macd><2", macdAboveBelow2)

	for _, v := range strat.Results {
		v.CalculatePNL()
		v.ChangeLessAlgo(backtest.LessAvgTrade)
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
