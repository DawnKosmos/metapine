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
	ch := e.Data()
	open := make([]float64, 0, len(ch))
	high := make([]float64, 0, len(ch))
	low := make([]float64, 0, len(ch))
	closes := make([]float64, 0, len(ch))
	volume := make([]float64, 0, len(ch))
	for _, c := range ch {
		open = append(open, c.Open)
		high = append(high, c.High)
		low = append(low, c.Low)
		closes = append(closes, c.Close)
		volume = append(volume, c.Volume)
	}
	o = empty(open, e.StartTime(), e.Resolution())
	h = empty(high, e.StartTime(), e.Resolution())
	c = empty(closes, e.StartTime(), e.Resolution())
	l = empty(low, e.StartTime(), e.Resolution())
	v = empty(volume, e.StartTime(), e.Resolution())
	return
}
