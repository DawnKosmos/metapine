package iterator

import "github.com/DawnKosmos/metapine/backend/series/ta"

/*
4h 0,5 normal distribution
SMA Open 10 5 Trades:2204 Winrate:65.25 Pnl: 50452.32
SMA HL2 10 7 Trades:1883 Winrate:63.73 Pnl: 51972.12
SMA HL2 11 13 Trades:1459 Winrate:62.65 Pnl: 55295.05
EMA OHCL4 4 5 Trades:1925 Winrate:63.79 Pnl: 55838.25
SMA OHCL4 10 5 Trades:2084 Winrate:62.81 Pnl: 80527.78
SMA HL2 10 5 Trades:2110 Winrate:62.75 Pnl: 95532.03


*/

func SolApeIter(oc2 ta.Series, volume ta.Series, ma func(s ta.Series, l int) ta.Series, len1, len2 int) (ta.Condition, ta.Condition) {
	outR := ta.Sma(ta.Roc(oc2, len1), 2)
	outB1 := ma(outR, len2)
	outB2 := ma(outB1, len2)
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
