package backtest

import (
	"github.com/DawnKosmos/metapine/helper/formula"
	"time"
)

type PNLGraph struct {
	Name string
	St   int64
	Res  int64
	Data []float64
}

func (t Trades) PnlGraph(name string, res int64, balance float64, end time.Time) *PNLGraph {
	st, et := t[0].EntrySignalTime.Unix(), end.Unix()
	arrLen := (et - st) / res
	pnl := make([]float64, arrLen, arrLen)

	mp := make(map[int]*Trade)
	for _, v := range t {
		mp[int((v.EntrySignalTime.Unix()-st)/res)] = v
	}

	pnl[0] = balance
	var index int
	for i := 0; i < int(arrLen); i++ {
		v, ok := mp[i]
		if ok {
			f := v.Pnl
			if index < i {
				for j := index; j <= i; j++ {
					pnl[j] = pnl[index]
					index = i
				}
			}
			for j, v := range f {
				pnl[j] += v
			}
			index = i + len(f)
		}

	}

	for i := index; i < int(arrLen); i++ {
		pnl[i] = pnl[index]
	}
	return &PNLGraph{
		Name: name,
		St:   st,
		Res:  res,
		Data: pnl,
	}
}

func (p *PNLGraph) ReturnPNL() []float64 {
	return p.Data
}

func Sync(pp ...*PNLGraph) {
	/*
		1) LowestStartTime
		2) FixLenght Beginning
		3) Filter KGV
		4) Fix ENDE
	*/

}

func pnlGraphMin(pp ...*PNLGraph) (int64, int) {
	val := pp[0].St
	var pos int
	for i, v := range pp[1:] {
		if v.St > val {
			pos = i
			val = v.St
		}
	}

	return val, pos
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func fixStartLen(pp *PNLGraph) {
	formula.Min()
}
