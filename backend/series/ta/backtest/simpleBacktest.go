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

/*



 */

func (b *BackTest) AddIndicator(series ...ta.Series) *BackTest {
	return b
}

func (b *BackTest) SetStrategy(buy, sell ta.Condition) *BackTestStrategy {

	return nil
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
