package live

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/helper/data"
)

type ema struct {
	URS[float64]
	src Series

	l           int
	alpha, last float64
}

func Ema(src Series, l int) *ema {
	s := new(ema)
	s.src = src
	s.ug = src.GetUpdateGroup()
	s.ug.AddUpdater(s)

	sma := ta.Ema(src, l)
	src.SetLimit(2)
	s.st, s.res = sma.StartTime(), sma.Resolution()
	s.data = data.Array(sma.Data())
	s.alpha = 2 / (float64(l) + 1)
	s.last = s.data.V(1)

	return s
}

func (s *ema) OnTick(new bool) {
	src0 := s.src.Val(0)
	if new {
		s.last = s.data.V(0)
		s.recent = src0 + 1
	}
	if s.recent != src0 {
		s.recent = src0
		avg := (src0-s.last)*s.alpha + s.last
		if new {
			s.data.Append(avg)
		} else {
			s.data.SetValue(0, avg)
		}
	}
}
