package ta

type constant[T any] struct {
	ERS[T]
}

func constantS(src Condition, a float64) Series {
	s := new(constant[float64])
	s.st = src.StartTime()
	s.res = src.Resolution()
	s.data = make([]float64, len(src.Data()), len(src.Data()))
	for i := range s.data {
		s.data[i] = a
	}
	return s
}

func constantB(src Condition, a bool) Condition {
	s := new(constant[bool])
	s.st = src.StartTime()
	s.res = src.Resolution()
	s.data = make([]bool, len(src.Data()), len(src.Data()))
	if a {
		for i := range s.data {
			s.data[i] = a
		}
	}
	return s
}
