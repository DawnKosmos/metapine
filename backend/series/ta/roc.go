package ta

type ROC struct {
	ERS[float64]
}

// Roc the Rate of Change is the equivalent to pines roc(src, len)
func Roc(src Series, l int) Series {
	s := new(ROC)
	s.res = src.Resolution()
	s.st = src.StartTime() + s.res*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)

	for i := l; i < len(f); i++ {
		d = append(d, 100*(f[i]-f[i-l])/f[i-l])
	}
	s.data = d
	return s
}
