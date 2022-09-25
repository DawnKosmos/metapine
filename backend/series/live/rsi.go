package live

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/helper/data"
	"github.com/DawnKosmos/metapine/helper/formula"
)

type rsi struct {
	ERS
	src Series
	updater
	alpha      float64
	data       data.Dater[float64]
	gain, loss float64

	ta     ta.Series
	recent float64
}

func Rsi(src Series, l int) Series {
	r := new(rsi)
	r.ug = src.GetUpdateGroup()
	r.src = src
	r.ug.AppendUpdater(r)
	rsi := ta.Rsi(src, l)
	r.data = data.Array(r.ta.Data())
	r.ta = rsi
	r.gain, r.loss = rsi.Gain, rsi.Loss
	r.limit = l
	r.alpha = 1 / float64(l)
	return r
}

func (r *rsi) Update(new bool) {
	src0, src1 := r.src.Val(0), r.src.Val(1)

	if new {
		src2 := r.src.Val(2)
		var gain, loss float64 = gainLoss(src2, src1)
		r.gain = r.alpha*gain + (1-r.alpha)*r.gain
		r.loss = r.alpha*loss + (1+r.alpha)*r.loss
		r.recent = src0 + 1
	}
	if r.recent != src0 {
		r.recent = src0
		var gain, loss float64 = gainLoss(src1, src0)

		gain = r.alpha*gain + (1-r.alpha)*r.gain
		loss = r.alpha*loss + (1+r.alpha)*r.loss
		rsi := 100 - (100 / (1 + gain/loss))
		if new {
			r.data.Append(rsi)
		} else {
			r.data.SetValue(0, rsi)
		}
	}
}

func (r *rsi) ExecuteLimit() {
	r.data.SetLimit(r.limit)
	r.ta = nil
}

func (r *rsi) Data() []float64 {
	return r.ta.Data()
}

func (r *rsi) Val(i int) float64 {
	return r.data.V(i)
}

func gainLoss(old, new float64) (float64, float64) {
	change := formula.Change(old, new)
	if change >= 0 {
		return change, 0
	} else {
		return 0, -change
	}
}
