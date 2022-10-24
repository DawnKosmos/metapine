package backtest

import "github.com/DawnKosmos/metapine/backend/series/ta"

type BackTest struct {
	ch ta.Chart
	//PNL starting with first Candle
	PnlCandle  []CandlePNL
	Pnl        []float64
	Indicators [][]SafeFloat
}

type BackTestStrategy struct {
	//buy, sell  ta.Condition
	//TE         TradeExecution
	Name       string
	Parameters BacktestParameters
	tr         []*Trade
}

func NewStrategy(ch ta.Chart) *BackTest {
	return &BackTest{
		ch:        ch,
		PnlCandle: []CandlePNL{},
		Pnl:       []float64{},
	}
}

func (b *BackTest) AddIndicator(indicators ...ta.Series) *BackTest {
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
		}
		i++
	}
	b.Indicators = indi
	return b
}

func (bt *BackTest) CreateStrategy(name string, buy, sell ta.Condition, TE TradeExecution, parameters BacktestParameters) *BackTestStrategy {
	var b = new(BackTestStrategy)

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
		Parameters: BacktestParameters{},
		tr:         nil,
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
			Name: b.Name + info + "true",
			tr:   tt,
		}, &BackTestStrategy{

			Parameters: BacktestParameters{},
			Name:       b.Name + info + "false",
			tr:         tf,
		}
}
