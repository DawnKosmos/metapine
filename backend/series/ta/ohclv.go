package ta

import (
	"time"

	"github.com/DawnKosmos/metapine/backend/exchange"
)

type Chart interface {
	Data() []exchange.Candle
	ResolutionStartTime
}

type OHCLV struct {
	ch  []exchange.Candle
	st  int64
	res int64
}

func NewOHCLV(e exchange.CandleProvider, ticker string, start time.Time, end time.Time, resolution int64) *OHCLV {
	o := new(OHCLV)
	o.res = resolution
	o.ch, _ = e.OHCLV(ticker, resolution, start, end)
	o.st = o.ch[0].StartTime.Unix()
	return o
}

func (o *OHCLV) Data() []exchange.Candle {
	return o.ch
}

func (o *OHCLV) StartTime() int64 {
	return o.st
}

func (o *OHCLV) Resolution() int64 {
	return o.res
}

func ChartSources(e Chart) (o, h, c, l, v Series) {
	var ff [5][]float64
	ch := e.Data()
	for i, c := range ch {
		ff[0][i], ff[1][i], ff[2][i], ff[3][i], ff[4][i] = c.Open, c.High, c.Close, c.Low, c.Volume
	}
	o = empty(ff[0], e.StartTime(), e.Resolution())
	h = empty(ff[1], e.StartTime(), e.Resolution())
	c = empty(ff[2], e.StartTime(), e.Resolution())
	l = empty(ff[3], e.StartTime(), e.Resolution())
	v = empty(ff[4], e.StartTime(), e.Resolution())
	return
}
