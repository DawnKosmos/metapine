package backtest

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/mode"
	"math"
	"time"

	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/helper/formula"
)

//fastBacktest should be used it you are permutation massive amounts of parameters

//Fast Backtest
type FastBacktest struct {
	ch         ta.Chart
	pyramiding int
	modus      mode.Mode
	st, et     int64
	cutEt      int
	res        int64
	fee        float64
	Indicators [][]SafeFloat
	P          []string
	results    []FastBacktestResult
	less       func(f *FastBacktestResult) float64
}

/*
NewFastBackTest should be used if you are iterating with a lot different parameters. It is not as detailed as an ordinary backtest but it does its job.
st, et -1 means trading whole Data
*/
func NewFastBackTest(ch ta.Chart, mode mode.Mode, pyramiding int, fee float64, st int64, et int64, parameter []string) *FastBacktest {
	var p int = 1
	if pyramiding > 1 {
		p = pyramiding
	}
	if st == -1 {
		st = ch.StartTime()
	}
	if et == -1 {
		et = time.Now().Unix()
	}
	var cutEt int
	var lastCandleTime int64 = ch.Data()[len(ch.Data())-1].StartTime.Unix()
	if et < lastCandleTime {
		cutEt = int((lastCandleTime - et) / ch.Resolution())
	}

	return &FastBacktest{
		ch:         ch,
		modus:      mode,
		pyramiding: p,
		st:         st,
		et:         et,
		fee:        fee,
		P:          parameter,
		res:        ch.Resolution(),
		cutEt:      cutEt,
		less:       SortPNL,
	}
}

func (bt *FastBacktest) AddIndicator(indicators ...ta.Series) *FastBacktest {
	if len(indicators) == 0 {
		return bt
	}

	d := bt.ch.Data()
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
	bt.Indicators = indi
	return bt
}

// Backtester Interface
func (f *FastBacktest) AddStrategy(buy ta.Condition, sell ta.Condition, paras string) {
	//Init
	ch, l, s := f.ch.Data(), buy.Data(), sell.Data()
	sl, pos := formula.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl+1:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]
	var st int64
	switch pos {
	case 0:
		st = f.ch.StartTime()
	case 1:
		st = buy.StartTime()
	case 2:
		st = sell.StartTime()
	}
	var iStart int
	if st < f.st {
		iStart = int((f.st - st) / f.res)
	}
	ch = ch[:len(ch)-f.cutEt]
	//

	var tempOrderLong, tempOrderShort []exchange.Candle //BackTest
	var trades []SimpleTrade

	for j := iStart; j < len(ch); j++ {
		if l[j] {
			for i := 0; i < min(len(tempOrderShort), f.pyramiding); i++ {
				t, err := CreateSimpleTrade(SHORT, tempOrderShort[i], ch[j])
				if err != nil {
					fmt.Println("Get Longs at", i, err)
					continue
				}
				trades = append(trades, t)
			}
			tempOrderShort = tempOrderShort[:0]
			if f.modus != mode.OnlySHORT {
				tempOrderLong = append(tempOrderLong, ch[j])
			}
		}
		if s[j] {
			for i := 0; i < min(len(tempOrderLong), f.pyramiding); i++ {
				t, err := CreateSimpleTrade(LONG, tempOrderLong[i], ch[j])
				if err != nil {
					fmt.Println("Get Longs at", i, err)
					continue
				}
				trades = append(trades, t)
			}
			tempOrderLong = tempOrderLong[:0]
			if f.modus != mode.OnlyLONG {
				tempOrderLong = append(tempOrderShort, ch[j])
			}
		}
	}

	f.results = append(f.results, newFastBacktestResult(trades, f.fee, f.less, paras))
}

type FastBacktestResult struct {
	name        string
	winrate     float64
	pnl         float64
	avgWin      float64
	TotalTrades int
	less        func(f *FastBacktestResult) float64
}

func newFastBacktestResult(tr []SimpleTrade, fee float64, less func(f *FastBacktestResult) float64, paras string) FastBacktestResult {
	gains := make([]float64, 0, len(tr))
	var wins int
	var pnl float64 = 1

	for i, v := range tr {
		gains = append(gains, v.Pnl(fee))
		if gains[i] > 1 {
			wins++
			pnl *= gains[i]
		} else {
			pnl *= gains[i]
		}
	}

	return FastBacktestResult{
		name:        paras,
		winrate:     float64(wins) / float64(len(tr)),
		pnl:         pnl,
		avgWin:      math.Pow(pnl, 1.0/float64(len(tr))),
		TotalTrades: len(tr),
		less:        less,
	}
}

func (p *FastBacktestResult) Print() {
	fmt.Println(fmt.Sprintf("%d %.2f, %.2f, %.2f", p.TotalTrades, p.winrate, p.pnl, p.avgWin))
}

// Sorting Algo
func (f *FastBacktest) SortingAlgo(less func(s *FastBacktestResult) float64) {
	f.less = less
}

func SortPNL(f *FastBacktestResult) float64 {
	return f.pnl
}

func SortWinrate(f *FastBacktestResult) float64 {
	return f.winrate
}

func SortAvgWin(f *FastBacktestResult) float64 {
	return f.avgWin
}

type FastBacktestResults []FastBacktestResult

func (a FastBacktestResults) Len() int           { return len(a) }
func (a FastBacktestResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FastBacktestResults) Less(i, j int) bool { return a[i].less(&a[i]) < a[j].less(&a[j]) }

func (f *FastBacktest) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}

func (f *FastBacktest) ReturnResults() []FastBacktestResult {
	return f.results
}
