package main

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/iterator"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
)
import . "github.com/DawnKosmos/metapine/backend/series/ta" //!! . importation means that no  nametag is needed for this package

/*
outR = sma(roc((open + close)/2, maR), 2)
outB1 = sma(outR, ma1)
outB2 = sma(outB1, ma1)
*/

func Example(ch *OHCLV, p backtest.Parameter) {
	oc := OC2(ch)
	c := Close(ch)
	l, h := Low(ch), High(ch)

	s14 := Sma(c, 14)
	e21 := Ema(c, 21)
	rsiOver := Rsi(c, 14)
	rsiEma := Ema(rsiOver, 9)
	vol := Volume(ch)
	b, s := iterator.SolApeIter(oc, vol, Sma, 4, 12)
	b = And(b, longCon(0.4, 6, c, l))
	s = And(s, shortCon(0.6, 6, c, h))

	//TE := tradeexecution.NewScaledLimit(0, 4, 5).Distribution(distribution.Exponential).Size(1.5) // 1 means 100%
	TE := tradeexecution.NewMarketOrder(1)
	bt := backtest.NewSimple(ch, TE, p)
	bt.AddIndicator(c, rsiOver, rsiEma, s14, e21)

	bt.AddStrategy(b, s, "Solape 4h")
	//bt.Split("Above SMA 14", backtest.Greater(1, 0))
	//	bt.Split("Rsi Above MA", backtest.Greater(1, 2))

	for _, v := range bt.Results {
		v.CalculatePNL()
	}

	for _, v := range bt.Results {
		fmt.Println(v.Name, fmt.Sprintf("Trades:%d Winrate:%.2f Pnl: %.2f", len(v.Trades()), v.Winrate, v.TotalPnl))
	}

}
