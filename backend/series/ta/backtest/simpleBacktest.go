package backtest

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"github.com/DawnKosmos/metapine/helper/formula"
	"sort"
)

type BackTest struct {
	ch ta.Chart
	//PNL starting with first Candle
	TE TradeExecution

	Parameters Parameter
	Indicators [][]SafeFloat
	Results    []*BackTestStrategy
}

type BackTestStrategy struct {
	Name string
	Pnl  []float64

	Parameters Parameter
	tr         []*Trade
	TotalPnl   float64
	Winrate    float64
	AvgTrade   float64
}

func (b *BackTestStrategy) Trades() []*Trade {
	return b.tr
}

func NewStrategy(ch ta.Chart, TE TradeExecution, parameters Parameter) *BackTest {
	return &BackTest{
		ch:         ch,
		TE:         TE,
		Parameters: parameters,
	}
}

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

	for _, vv := range indicators[1:] {

		var i int = 1
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

//Backtester Interface

func (bt *BackTest) AddStrategy(buy, sell ta.Condition, values string) {
	var b = new(BackTestStrategy)
	ch, l, s := bt.ch.Data(), buy.Data(), sell.Data()
	sl, _ := formula.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]

	if bt.TE.GetInfo().CandlePnlSupport && bt.Parameters.PnlGraph {
		b.Pnl = make([]float64, len(ch), len(ch))
		bt.Parameters.PnlGraph = false
	}
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
	b.Name = values
	bt.Results = append(bt.Results, b)
}

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
			bb = append(bb, &BackTestStrategy{
				Name:       vv.Name + condition + "true",
				Parameters: vv.Parameters,
				tr:         tt,
			}, &BackTestStrategy{

				Parameters: vv.Parameters,
				Name:       vv.Name + condition + "false",
				tr:         tf,
			})
		}
	}
	bt.Results = bb
}

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
		}, &BackTestStrategy{

			Parameters: vv.Parameters,
			Name:       vv.Name + "Shorts",
			tr:         tf,
		})
	}
	bt.Results = bb
}

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
			Name:       vv.Name + condition + " true",
			Parameters: vv.Parameters,
			tr:         tt,
		})
	}
	bt.Results = bb
}

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
}

/*
func (bt *BackTest) CreateStrategy(name string, buy, sell ta.Condition, TE TradeExecution, parameters Parameter) *BackTestStrategy {
	var b = new(BackTestStrategy)
	b.Name = name
	b.Parameters = parameters
	if b.Parameters.Pyramiding == 0 {
		b.Parameters.Pyramiding = 1
	}
	p := b.Parameters.Pyramiding

	ch, l, s := bt.ch.Data(), buy.Data(), sell.Data()
	sl, _ := formula.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]
	b.Pnl = make([]float64, len(ch), len(ch))
	var indicators [][]SafeFloat
	if bt.Indicators != nil {
		indicators = bt.Indicators[len(bt.Indicators)-sl:]
	} else {
		indicators = make([][]SafeFloat, sl, sl)
	}

	var indexLong, indexShort []int
	var tr []*Trade

	if !TE.GetInfo().CandlePnlSupport {
		b.Parameters.PnlGraph = false
	}
	balance := parameters.Balance
	var tempBalance float64

	for j := 0; j < len(ch)-1; j++ {
		if l[j] {
			for i := 0; i < min(len(indexShort), p); i++ {
				index := indexShort[i]
				t, err := TE.CreateTrade(SHORT, ch[index+1:], j-index, indicators[index], balance, *parameters.Fee, b.Parameters.PnlGraph)
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
			if parameters.Modus != OnlySHORT {
				indexLong = append(indexLong, j)
			}
		}
		if s[j] {
			for i := 0; i < min(len(indexLong), p); i++ {
				index := indexLong[i]
				t, err := TE.CreateTrade(LONG, ch[index+1:], j-index, indicators[index], balance, *parameters.Fee, b.Parameters.PnlGraph)
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
			if parameters.Modus != OnlyLONG {
				indexShort = append(indexShort, j)
			}
		}
	}
	sort.Sort(Trades(tr))
	b.tr = tr
	return b
}

func (b *BackTestStrategy) Filter(info string, op Filter) *BackTestStrategy {
	var tt []*Trade
	for _, v := range b.tr {
		if op(v.Indicator) {
			tt = append(tt, v)
		}
	}
	return &BackTestStrategy{
		Name:       b.Name + info,
		Parameters: b.Parameters,
		tr:         tt,
	}
}

func (b *BackTestStrategy) Split(info string, op Filter) (*BackTestStrategy, *BackTestStrategy) {
	var tt, tf []*Trade
	for _, v := range b.tr {
		if op(v.Indicator) {
			tt = append(tt, v)
		} else {
			tf = append(tf, v)
		}
	}

	return &BackTestStrategy{
			Name:       b.Name + info + "true",
			Parameters: b.Parameters,
			tr:         tt,
		}, &BackTestStrategy{

			Parameters: b.Parameters,
			Name:       b.Name + info + "false",
			tr:         tf,
		}
}

func (b *BackTestStrategy) LongsAndShorts() (*BackTestStrategy, *BackTestStrategy) {
	var tt, tf []*Trade
	for _, v := range b.tr {
		if v.Side {
			tt = append(tt, v)
		} else {
			tf = append(tf, v)
		}
	}
	return &BackTestStrategy{
			Name:       b.Name + "| longs",
			Parameters: b.Parameters,
			tr:         tt,
		}, &BackTestStrategy{

			Parameters: b.Parameters,
			Name:       b.Name + "| shorts",
			tr:         tf,
		}
}
*/

// TRADES

type BackTestStrategies []*BackTestStrategy

func (t BackTestStrategies) Less(i, j int) bool {
	return t[i].TotalPnl < t[j].TotalPnl
}

func (t BackTestStrategies) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t BackTestStrategies) Len() int {
	return len(t)
}
