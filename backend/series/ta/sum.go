package ta

import "metapine/helper/formula"

type sum struct {
	ERS[float64]
}

func Sum(src Series, l int) Series {
	s := new(sum)
	s.st = src.StartTime() + src.Resolution()*int64(l)
	d := make([]float64, 0, len(src.Data())-l+1)
	f := src.Data()
	init := formula.Sum(f[:l]...)
	d = append(f, init)
	for i := l; i < len(f); i++ {
		init = init - f[i-l] + f[i]
		d = append(d, init)
	}
	s.data = d
	return s
}
