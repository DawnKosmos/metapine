package iterator

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"testing"
)

type Example struct {
	src1 ta.Series
	l1   int
	l2   int
}

func (e *Example) StructsAdresse() ([]*int, []*ta.Series, []*func(series ta.Series, l int) ta.Series) {
	return []*int{&e.l1, &e.l2}, []*ta.Series{&e.src1}, nil
}

func (e *Example) Calculation() (buy, sell ta.Condition) {
	fmt.Println(e.src1.Name(), e.l1, e.l2)
	return solape(e.src1, e.l1, e.l2)
}

func (e *Example) Parameters() string {
	return fmt.Sprintf("%s \t %d \t%d", e.src1.Name(), e.l1, e.l2)
}

func solape(oc2 ta.Series, len1, len2 int) (ta.Condition, ta.Condition) {
	outR := ta.Sma(ta.Roc(oc2, len1), 2)
	outB1 := ta.Sma(outR, len2)
	outB2 := ta.Sma(outB1, len2)
	outB := ta.SubF(outB1, outB2, 2.0)
	cc := ta.Sub(outR, outB)
	var c1 ta.Series
	c1 = ta.Sma(cc, 2)
	c2, c3 := ta.OffS(c1, 1), ta.OffS(c1, 2)
	buy := ta.And(ta.Greater(c1, c2), ta.Smaller(c2, c3))
	sell := ta.And(ta.Smaller(c1, c2), ta.Greater(c2, c3))

	return buy, sell
}

func TestIterator(t *testing.T) {
	ee := ftx.New()
	ch := ta.NewOHCLV(ee, "BTC-PERP", exchange.T2022, exchange.T2023, 3600*24)
	close, oc2, hl2, open := ta.Close(ch), ta.OC2(ch), ta.HL2(ch), ta.Open(ch)

	bt := backtest.InitFastBackTest(ch, backtest.ALL, 4, 0, -1, -1, []string{})

	it := New(&Example{})
	it.RegisterInt(0, 5, 10, 1)
	it.RegisterInt(1, 4, 8, 1)
	it.RegisterSeries(0, close, oc2, hl2, open)

	it.Run(bt)

	fmt.Println(len(bt.ReturnResults()))
}
