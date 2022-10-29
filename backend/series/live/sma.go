package live

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/helper/data"
)

type sma struct {
	URS[float64]
	src Series

	l           int
	alpha       float64
	last, src14 float64
}

func Sma(src Series, l int) *sma {
	s := new(sma)
	s.src = src
	s.ug = src.GetUpdateGroup()
	s.ug.AddUpdater(s)

	sma := ta.Sma(src, l)
	src.SetLimit(l + 1)
	s.st, s.res = sma.StartTime(), sma.Resolution()
	s.l = l
	s.data = data.Array(sma.Data())
	s.alpha = 1 / float64(l)
	s.src14 = src.Val(s.l)
	return s
}

func (s *sma) OnTick(new bool) {
	src0 := s.src.Val(0)
	if new {
		s.src14 = s.src.Val(s.l)
		s.last = s.data.V(0)
		s.recent = src0 + 1
	}
	if s.recent != src0 {
		s.recent = src0

		avg := s.last + (src0-s.src14)*s.alpha
		if new {
			s.data.Append(avg)
		} else {
			s.data.SetValue(0, avg)
		}
	}
}
