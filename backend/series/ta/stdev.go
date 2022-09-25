package ta

import "math"

type STDEV struct {
	ERS[float64]
}

func Stdev(src Series, l int) Series {
	var s STDEV
	s.res = src.Resolution()
	s.st = src.StartTime() + s.res*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)

	var m, oldM, S float64
	var m1, s1 float64

	for i := 0; i < l; i++ {
		x := f[i]
		oldM = m
		m = m + (x-m)/float64(i+1)
		S += (x - m) * (x - oldM)
	}

	lf := float64(l)
	d = append(d, math.Sqrt(S/lf))

	for i := l; i < len(f); i++ {
		x1 := f[i-l]
		m1 = (lf*m - x1) / (lf - 1)
		s1 = S - (x1-m1)*(x1-m)

		x := f[i]
		m = m1 + (x-m1)/lf
		S = s1 + (x-m)*(x-m1)

		//fmt.Println(S, s1)
		d = append(d, math.Sqrt(S/lf))
	}

	s.data = d
	return &s
}
