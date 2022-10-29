package live

import "github.com/DawnKosmos/metapine/backend/series/ta"

type offsetS struct {
	URS[float64]
	src    Series
	offset int
}

func Offs(src Series, l int) Series {
	s := new(offsetS)

	off := ta.OffS(src, l)
	s.st, s.res = off.StartTime(), off.Resolution()
	s.offset = l
	s.src = src

	return s
}

func (s *offsetS) Val(i int) float64 {
	return s.src.Val(i + s.offset)
}

func (s *offsetS) SetLimit(limit int) {
	s.src.SetLimit(limit + s.offset)
}

func (s *offsetS) OnTick(new bool) {
	return
}

func (s *offsetS) ExecuteLimit() {
	return
}

func (s *offsetS) Data() []float64 {
	return s.src.Data()[:len(s.src.Data())-s.offset]
}
