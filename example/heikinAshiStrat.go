package main

import (
	"fmt"
	. "github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/tradeexecution"
)

func HeikinAshiStrategy(ch *OHCLV) {
	ha := HeikinAshi(ch)
	o, c := Open(ha), Close(ha)

	TE := tradeexecution.NewScaledLimit(0, 2, 10).Distribution(distribution.Exponential).Size(0.25)
	bt := backtest.NewSimple(ch, TE, paras)

	o1, c1 := OffS(o, 1), OffS(c, 1)

	buy := And(Greater(c, o), Smaller(c1, o1))
	sell := And(Greater(c1, o1), Smaller(c, o))

	_, _, macd := MacdRelative(Close(ch), 12, 26, 9)

	bt.AddIndicator(macd)

	bt.AddStrategy(buy, sell, "strat")

	bt.Filter("macd", func(sf []backtest.SafeFloat) bool {
		return true // math.Abs(sf[0].Value) > 0.05
	})

	for _, v := range bt.Results {
		v.CalculatePNL()
	}

	for _, v := range bt.Results {
		fmt.Println(v.Name, fmt.Sprintf("Trades:%d Winrate:%.2f Pnl: %.2f", len(v.Trades()), v.Winrate, v.TotalPnl))
	}

}
