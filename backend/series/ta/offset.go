package ta

type offset[T any] struct {
	ERS[T]
}

func OffS(src Series, n int) Series {
	s := new(offset[float64])
	s.res = src.Resolution()
	s.st = s.res*int64(n) + src.StartTime()
	s.data = src.Data()[:len(src.Data())-n]
	return s
}

func OffC(src Condition, n int) Condition {
	s := new(offset[bool])
	s.res = src.Resolution()
	s.st = s.res*int64(n) + src.StartTime()
	s.data = src.Data()[:len(src.Data())-n]
	return s
}
