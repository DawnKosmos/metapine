package ta

import (
	"log"

	"github.com/DawnKosmos/metapine/helper/formula"
)

type LOWEST struct {
	ERS[float64]
}

func Lowest(src Series, l int) Series {
	if l < 2 {
		log.Panicln("Invalid lenght lowest", l)
	}
	s := new(LOWEST)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)
	lo, pos := formula.Lowest(f[:l]...)
	d = append(d, lo)
	for i := l; i < len(f); i++ {
		if pos < i-l {
			lo, pos = formula.Lowest(f[i-l : i]...)
			pos += i
		}
		if f[i] <= lo {
			lo = f[i]
			pos = i
		}
		d = append(d, lo)
	}

	s.data = d
	return s
}
