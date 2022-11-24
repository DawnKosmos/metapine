package main

import (
	"fmt"
	. "github.com/DawnKosmos/metapine/backend/series/ta" //!! . importation means that no  nametag is needed for this package
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/iterator"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
	"sort"
	"time"
)

func IterExample(ch *OHCLV, p backtest.Parameter) {
	// Get Different sources
	open := Open(ch)
	oc := OC2(ch)
	high := High(ch)
	volume := Volume(ch)

	//Wrap the MFI indicator to use it for iteration
	WrappedMFI := func(src Series, l int) Series {
		return MFI(src, volume, l)
	}

	//TE executes the trades by placing 10 orders from 0,10 with equal size, which sum is equal to the total trade size
	TE := tradeexecution.NewScaledLimit(0, 5, 10).Distribution(distribution.Normal).Size(1) // 1 means 100%
	//A New Startegy
	strat := backtest.NewStrategy(ch, TE, p)

	//Timestap to measure how long the permutation took
	tNow := time.Now()
	//
	fn := func(src Series, maFunc MaFunc, fast, slow int) (buy, sell Condition) {
		return SolApeIter(src, volume, maFunc, fast, slow)
	}
	iter := IteratorExample{fn: fn}
	it := iterator.New(&iter)
	// Iterating the first int which parameter name is "fast" from 4 to 12 in +1 steps
	it.RegisterInt(0, 4, 12, 1)
	// Iterating the second int which parameter name is "slow" from 3 to 14 in +2 steps
	it.RegisterInt(1, 3, 14, 2)
	// Iterating the fastMa
	it.RegisterFunctions(0, Sma, WrappedRsi, Ema, WrappedMFI)
	// Iteration The Sources
	it.RegisterSeries(0, oc, open, HL2(ch), OHCL4(ch), high)

	it.Run(strat)

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

/*
This is an Example implementation of the Iterator interface
*/
type IteratorExample struct {
	src    Series                                                            //The Source for the calculation of the indicator
	maFast func(s1 Series, l int) Series                                     //A Moving Average Function
	fast   int                                                               // length of the indictor
	slow   int                                                               // second length of the indicator
	fn     func(src Series, maFunc MaFunc, l1, l2 int) (buy, sell Condition) //The Function for calculating the Startegies buy & sell signals
}

// StructsAdresse return every Pointer to Parameter(Struct Field), you want to change, in arrays
func (it *IteratorExample) StructsAdresse() ([]*int, []*Series, []*func(src Series, l int) Series) {
	return []*int{&it.fast, &it.slow}, []*Series{&it.src}, []*func(src Series, l int) Series{&it.maFast}
}

func (it *IteratorExample) Calculation() (buy, sell Condition) {
	return it.fn(it.src, it.maFast, it.fast, it.slow)
}

// Parameters returns a String of the Parameters Value, which is needed to separate the different Iterations
func (it *IteratorExample) Parameters() string {
	//Every Indicator has its own Name(), but maFast(func(s1 Series, l int) Series) does not have it
	//Therefore we have to Wrap this function to have  a way to differentiate. See below
	return fmt.Sprintf("%s %s %d %d", MaFunc(it.maFast).Name(), it.src.Name(), it.fast, it.slow)
}

type MaFunc func(s Series, i int) Series

func (m MaFunc) Name() string {
	return m(Constant([]float64{1, 2, 3, 4}, 0, 3600, ""), 2).Name()
}
