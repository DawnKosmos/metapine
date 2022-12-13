package backtest

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"github.com/DawnKosmos/metapine/helper/formula"
	"sort"
)

// BackTest is what it says a struct to backtest
type BackTest struct {
	//OHCLV Data needed for calculation
	ch ta.Chart
	//TE describes how Trades get executed. You can choose simple market orders, or set multiple different limit orders or set stops
	TE TradeExecution
	//Lookup Parameters
	Parameters Parameter
	//Indicators are saved in a [][]SafeFloat(bool,float64) Array synched to the OHCLV data
	Indicators [][]SafeFloat
	//Lookup BackTestStrategy
	Results []*BackTestStrategy
}

// BacktestStrategy saves the Results
type BackTestStrategy struct {
	//Name has the strategy name in it, and saves different parameters
	Name string
	//PNL is the PNLchart, right now not implemented
	Pnl []float64

	//Lookup Parameters
	Parameters Parameter
	//All the trades that got executed
	tr []*Trade
	//Sum of Pnl of the Trades
	TotalPnl float64
	//Winrate of the Trades
	Winrate float64
	//AvgTrade Gain
	AvgTrade float64
	lessAlgo func(v *BackTestStrategy) float64
}

func (b *BackTestStrategy) Trades() []*Trade {
	return b.tr
}

func NewSimple(ch ta.Chart, TE TradeExecution, parameters Parameter) *BackTest {
	return &BackTest{
		ch:         ch,
		TE:         TE,
		Parameters: parameters,
	}
}

// AddIndicators, fills the [][]SafeFloat with Series, OHCLV is also a Series
func (bt *BackTest) AddIndicator(indicators ...ta.Series) *BackTest {
	if len(indicators) == 0 {
		return bt
	}

	d := bt.ch.Data()
	indi := make([][]SafeFloat, 0, len(d)) //init array
	f := indicators[0].Data()              // Data of first array
	l1 := len(indicators)

	for i := 0; i < len(d)-len(f); i++ {
		init := make([]SafeFloat, l1, l1) //init t
		indi = append(indi, init)
	}

	var j int = len(d) - len(f)
	for _, v := range f {
		init := make([]SafeFloat, l1, l1)
		init[0] = SafeFloat{Safe: true, Value: v}
		indi = append(indi, init)
		j++
	}

	var i int = 1
	for _, vv := range indicators[1:] {
		f = vv.Data()
		j = len(d) - len(f)
		for _, v := range f {
			indi[j][i] = SafeFloat{Safe: true, Value: v}
			j++
		}
		i++
	}
	bt.Indicators = indi
	return bt
}

// AddStrategy , Strategies are created with buy and sell signals(just bool) They open a Position and close the  contrarian one.
func (bt *BackTest) AddStrategy(buy, sell ta.Condition, Name string) {
	var b = new(BackTestStrategy)
	//OHCLV, buy and sell have to match up the same size
	ch, l, s := bt.ch.Data(), buy.Data(), sell.Data()
	sl, _ := formula.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]
	b.lessAlgo = LessPnl

	//Check if PnlGraph is supported
	if bt.TE.GetInfo().CandlePnlSupport && bt.Parameters.PnlGraph {
		b.Pnl = make([]float64, len(ch), len(ch))
		bt.Parameters.PnlGraph = false
	}

	//Check if Indicators were added
	var indicators [][]SafeFloat
	if bt.Indicators != nil {
		indicators = bt.Indicators[len(bt.Indicators)-sl:]
	} else {
		indicators = make([][]SafeFloat, sl, sl)
	}
	var indexLong, indexShort []int
	var tr []*Trade

	balance := bt.Parameters.Balance
	parameters := bt.Parameters
	var tempBalance float64

	//Trades get Created here, this is a Simple Backtest. It does not support, having buy and sell strategies running next to each other.
	for j := 0; j < len(ch)-1; j++ {
		if l[j] {
			for i := 0; i < min(len(indexShort), parameters.Pyramiding); i++ {
				index := indexShort[i]
				t, err := bt.TE.CreateTrade(SHORT, ch[index+1:], j-index, indicators[index], balance, *bt.Parameters.Fee, b.Parameters.PnlGraph)
				if err != nil {
					//fmt.Println("Create Shorts at", j, err)
					continue
				}
				tr = append(tr, t)
				if parameters.SizeType == size.Account {
					tempBalance += t.RealisedPNL()
				}
			}
			if parameters.SizeType == size.Account {
				balance += tempBalance
				tempBalance = 0
			}

			indexShort = indexShort[:0]
			if parameters.Modus != mode.OnlySHORT {
				indexLong = append(indexLong, j)
			}
		}
		if s[j] {
			for i := 0; i < min(len(indexLong), parameters.Pyramiding); i++ {
				index := indexLong[i]
				t, err := bt.TE.CreateTrade(LONG, ch[index+1:], j-index, indicators[index], balance, *parameters.Fee, b.Parameters.PnlGraph)
				if err != nil {
					//fmt.Println("Create Longs at", j, err)
					continue
				}
				tr = append(tr, t)
				if parameters.SizeType == size.Account {
					tempBalance += t.RealisedPNL()
				}
			}
			if parameters.SizeType == size.Account {
				balance += tempBalance
				tempBalance = 0
			}

			indexLong = indexLong[:0]
			if parameters.Modus != mode.OnlyLONG {
				indexShort = append(indexShort, j)
			}
		}
	}

	sort.Sort(Trades(tr))
	b.tr = tr
	b.Name = Name
	bt.Results = append(bt.Results, b)
}

// Split turns 2 Results Sets with Conditions
func (bt *BackTest) Split(condition string, op Filter) {
	var bb []*BackTestStrategy
	for _, vv := range bt.Results {
		var tt, tf []*Trade
		for _, v := range vv.tr {
			if op(v.Indicator) {
				tt = append(tt, v)
			} else {
				tf = append(tf, v)
			}
		}

		if len(tt) > 0 {
			bb = append(bb, &BackTestStrategy{
				Name:       vv.Name + "\t" + condition + " true",
				Parameters: vv.Parameters,
				tr:         tt,
				lessAlgo:   vv.lessAlgo,
			})
		}
		if len(tf) > 0 {
			bb = append(bb, &BackTestStrategy{
				Parameters: vv.Parameters,
				Name:       vv.Name + "\t" + condition + " false",
				tr:         tf,
				lessAlgo:   vv.lessAlgo,
			})
		}
	}
	bt.Results = bb
}

// Splits LongsFromShorts
func (bt *BackTest) LongShort() {
	var bb []*BackTestStrategy
	for _, vv := range bt.Results {
		var tt, tf []*Trade
		for _, v := range vv.tr {
			if v.Side {
				tt = append(tt, v)
			} else {
				tf = append(tf, v)
			}
		}
		bb = append(bb, &BackTestStrategy{
			Name:       vv.Name + "Longs",
			Parameters: vv.Parameters,
			tr:         tt,
			lessAlgo:   vv.lessAlgo,
		}, &BackTestStrategy{

			Parameters: vv.Parameters,
			Name:       vv.Name + "Shorts",
			tr:         tf,
			lessAlgo:   vv.lessAlgo,
		})
	}
	bt.Results = bb
}

// Deletes The Trades if given conditions is false
func (bt *BackTest) Filter(condition string, op Filter) {
	var bb []*BackTestStrategy
	for _, vv := range bt.Results {
		var tt []*Trade
		for _, v := range vv.tr {
			if op(v.Indicator) {
				tt = append(tt, v)
			}
		}
		bb = append(bb, &BackTestStrategy{
			Name:       vv.Name + "\t" + condition + " true",
			Parameters: vv.Parameters,
			tr:         tt,
			lessAlgo:   vv.lessAlgo,
		})
	}
	bt.Results = bb
}

// CalculatePNL Total PNL
func (bt *BackTestStrategy) CalculatePNL() {
	var tpnl float64
	var win int

	for _, v := range bt.tr {
		pnl := v.RealisedPNL()
		if pnl > 0 {
			win++
		}
		tpnl += pnl
	}

	bt.Winrate = float64(win) / float64(len(bt.tr)) * 100
	bt.TotalPnl = tpnl
	bt.AvgTrade = bt.TotalPnl / float64(len(bt.tr))
}

func (bt *BackTestStrategy) ChangeLessAlgo(fn func(b *BackTestStrategy) float64) {
	bt.lessAlgo = fn
}

func LessPnl(b *BackTestStrategy) float64 {
	return b.TotalPnl
}

func LessWinrate(b *BackTestStrategy) float64 {
	return b.Winrate
}

func LessAvgTrade(b *BackTestStrategy) float64 {
	return b.AvgTrade
}

type BackTestStrategies []*BackTestStrategy

func (t BackTestStrategies) Less(i, j int) bool {
	return t[i].lessAlgo(t[i]) < t[j].lessAlgo(t[j])
}

func (t BackTestStrategies) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t BackTestStrategies) Len() int {
	return len(t)
}
