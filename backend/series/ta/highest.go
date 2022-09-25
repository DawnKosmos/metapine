package ta

import (
	"log"
	"metapine/helper/formula"
)

type HIGHEST struct {
	ERS[float64]
}

func Highest(src Series, l int) Series {
	if l < 2 {
		log.Panicln("Invalid lenght lowest", l)
	}
	s := new(LOWEST)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)
	high, pos := formula.Highest(f[:l]...)
	d = append(d, high)
	for i := l; i < len(f); i++ {
		if pos < i-l {
			high, pos = formula.Highest(f[i-l : i]...)
			pos += i
		}
		if f[i] >= high {
			high = f[i]
			pos = i
		}
		d = append(d, high)
	}
	return s
}
