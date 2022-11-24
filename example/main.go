package main

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange/deribit"
	"github.com/DawnKosmos/metapine/backend/exchange/psql"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/size"
	"os"
	"time"
)

/*
TODO
Seperate PSQL and DB logic
Add Enviroments to Iterators
Think of filtering

Visual representation
*/

var paras = backtest.Parameter{
	Modus:      mode.ALL,
	Pyramiding: 1,
	Fee: &backtest.Fee{
		Maker:    -0.00005,
		Taker:    0.0005,
		Slippage: 0,
	},
	Balance:  10000,
	SizeType: size.Account,
	PnlGraph: false,
}

func main() {
	psql.SetPSQL("localhost", "postgres", "metapine", "admin", 5432)
	//	exch, _ := psql.New("ftx")
	exch := deribit.New()

	ch := ta.NewOHCLV(exch, "BTC-PERPETUAL", time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 11, 23, 0, 0, 0, 0, time.UTC), 4*3600)
	if len(ch.Data()) == 0 {
		os.Exit(1)
	}

	FilterExample(ch, paras)
}

type IteratorRibbon struct {
	src1  ta.Series
	ma    func(s ta.Series, i int) ta.Series
	lenMa int
	lenHl int
	fn    func(src1 ta.Series, maFunc func(s ta.Series, i int) ta.Series, lenMa, lenHL int) (buy, sell ta.Condition)
}

func (i *IteratorRibbon) StructsAdresse() ([]*int, []*ta.Series, []*func(src ta.Series, l int) ta.Series) {
	return []*int{&i.lenMa, &i.lenHl}, []*ta.Series{&i.src1}, []*func(s ta.Series, i int) ta.Series{&i.ma}
}

func (i *IteratorRibbon) Calculation() (buy, sell ta.Condition) {
	return i.fn(i.src1, i.ma, i.lenMa, i.lenHl)
}

func (i *IteratorRibbon) Parameters() string {
	return fmt.Sprintf("%s \t%s\t%d \t%d", i.src1.Name(), MaFunc(i.ma).Name(), i.lenMa, i.lenHl)
}

func solape(oc2 ta.Series, volume ta.Series, len1, len2 int) (ta.Condition, ta.Condition) {
	outR := ta.Sma(ta.Roc(oc2, len1), 2)
	outB1 := ta.Sma(outR, len2)
	outB2 := ta.Sma(outB1, len2)
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

func WrappedRsi(src ta.Series, l int) ta.Series {
	return ta.Rsi(src, l)
}
