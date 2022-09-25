package ta

import (
	"math"
)

type TRange struct {
	ERS[float64]
}

func TrueRange(close, high, low Series) *TRange {
	s := new(TRange)
	c, h, l := close.Data(), high.Data(), low.Data()
	s.st = close.StartTime()
	s.res = close.Resolution()

	d := make([]float64, 0, len(c))
	var hlasb, abslc float64

	for i := 1; i < len(c); i++ {
		hlasb = math.Max(h[i]-l[i], math.Abs(h[i]-c[i-1]))
		abslc = math.Abs(l[i] - c[i-1])
		d = append(d, math.Max(hlasb, abslc))
	}
	s.data = d
	return s
}
