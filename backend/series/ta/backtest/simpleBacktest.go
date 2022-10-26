package backtest

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"github.com/DawnKosmos/metapine/helper/formula"
	"sort"
)

type BackTest struct {
	ch ta.Chart
	//PNL starting with first Candle
	Indicators [][]SafeFloat
}

type BackTestStrategy struct {
	//buy, sell  ta.Condition
	//TE         TradeExecution
	Name string
	Pnl  []float64

	Parameters BTParameter
	tr         []*Trade
}

type BTParameter struct {
	Modus      Mode
	Pyramiding int
	Fee        *FeeInfo
	Balance    float64
	SizeType   size.SizeBase
}

func NewStrategy(ch ta.Chart) *BackTest {
	return &BackTest{
		ch: ch,
	}
}

func (b *BackTest) AddIndicator(indicators ...ta.Series) *BackTest {
	if len(indicators) == 0 {
		return b
	}

	d := b.ch.Data()
	indi := make([][]SafeFloat, len(d), len(d))
	f := indicators[0].Data()
	l1 := len(indicators)

	for i := 0; i < len(d)-len(f); i++ {
		init := make([]SafeFloat, l1, l1)
		indi = append(indi, init)
	}

	var j int = len(d) - len(f)
	for _, v := range f {
		init := make([]SafeFloat, l1, l1)
		init[0] = SafeFloat{Safe: true, Value: v}
		indi[j] = init
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
	b.Indicators = indi
	return b
}

func (bt *BackTest) CreateStrategy(name string, buy, sell ta.Condition, TE TradeExecution, parameters BTParameter) *BackTestStrategy {
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
	}

	var indexLong, indexShort []int
	var tr []*Trade

	fmt.Println(len(ch))
	for j := 0; j < len(ch)-1; j++ {
		if l[j] {
			for i := 0; i < min(len(indexShort), p); i++ {
				index := indexShort[i]
				t, err := TE.CreateTrade(SHORT, ch[index+1:], j-index, indicators[index], parameters.Balance, *parameters.Fee)
				if err != nil {
					fmt.Println("Create Shorts at", i, err)
					continue
				}
				tr = append(tr, t)
			}
			indexShort = indexShort[:0]
			if parameters.Modus != OnlySHORT {
				indexLong = append(indexLong, j)
			}
		}
		if s[j] {
			for i := 0; i < min(len(indexLong), p); i++ {
				index := indexLong[i]
				t, err := TE.CreateTrade(LONG, ch[index+1:], j-index, indicators[index], parameters.Balance, *parameters.Fee)
				if err != nil {
					fmt.Println("Create Longs at", i, err)
					continue
				}
				tr = append(tr, t)
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
