package ta

type offsetS struct {
	ERS[float64]
	src    Series
	offset int
}

func OffS(src Series, n int) Series {
	s := new(offsetS)
	s.res = src.Resolution()
	s.st = s.res*int64(n) + src.StartTime()
	s.src = src
	return s
}

func (s *offsetS) Data() []float64 {
	f := s.src.Data()
	return f[:len(f)-s.offset-1]
}

type offsetC struct {
	ERS[bool]
	src    Condition
	offset int
}

func OffC(src Condition, n int) Condition {
	s := new(offsetC)
	s.res = src.Resolution()
	s.st = s.res*int64(n) + src.StartTime()
	s.src = src
	return s
}

func (s *offsetC) Data() []bool {
	f := s.src.Data()
	return f[:len(f)-s.offset-1]
}
