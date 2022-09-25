package ta

type td struct {
	ERS[float64]
}

func TD(src Series, l int) Series {
	s := new(td)
	s.res = src.Resolution()
	s.st = src.StartTime() + s.res*int64(l)

	f := src.Data()

	out := make([]float64, 0, len(f))
	var val float64

	for i, v := range f[l:] {
		if f[i] < v {
			if val < 0 {
				val = 1.0
			} else {
				val += 1.0
			}
		} else {
			if val > 0 {
				val = -1.0
			} else {
				val -= 1.0
			}
		}
		out = append(out, val)
	}
	s.data = out
	return s
}
