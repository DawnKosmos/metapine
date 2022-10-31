package ta

type swma struct {
	ERS[float64]
}

func Swma(src Series) Series {
	s := new(swma)
	s.st = src.StartTime() + src.Resolution()*4
	s.res = src.Resolution()

	s.name = "SWMA"
	f := src.Data()
	d := make([]float64, 0, len(f))
	x3, x2, x1, x0 := 1/6*f[0], 2/6*f[1], 2/6*f[2], 1/6*f[3]
	for i := 3; i < len(f); i++ {
		x3 = x2 / 2
		x2 = x1
		x1 = x0 * 2
		x0 = f[i] * 1 / 6
		d = append(d, x0+x1+x2+x3)
	}
	s.data = d

	return s
}
