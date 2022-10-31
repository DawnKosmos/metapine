package ta

import (
	"log"

	"github.com/DawnKosmos/metapine/helper/formula"
)

type SMA struct {
	ERS[float64]
}

func Sma(src Series, l int) Series {
	if l < 2 {
		log.Panicln("sma invalid len", l)
	}

	s := new(SMA)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l-1)
	f := src.Data()
	l1 := len(f)
	d := make([]float64, 0, l1-l+1)
	s.name = "SMA"
	avg := formula.Average(f[:l]...)
	d = append(d, avg)
	alpha := 1 / float64(l)

	for i := l; i < len(f); i++ {
		avg = avg + (f[i]-f[i-l])*alpha
		d = append(d, avg)
	}
	s.data = d
	return s
}

type EMA struct {
	ERS[float64]
	Avg float64
}

func Ema(src Series, l int) Series {
	if l < 2 {
		log.Panicln("ema invalid len", l)
	}
	s := new(EMA)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)
	alpha := 2 / (float64(l) + 1)
	avg := formula.Average(f[:l]...)
	d = append(d, avg)
	for i := l; i < len(f); i++ {
		avg = (f[i]-avg)*alpha + avg
		d = append(d, avg)
	}
	s.name = "EMA"
	s.Avg = avg // For live trading
	s.data = d
	return s
}

type RMA struct {
	ERS[float64]
	Avg float64
}

func Rma(src Series, l int) *RMA {
	if l < 2 {
		log.Panicln("rma invalid len", l)
	}
	s := new(RMA)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)
	alpha := 1 / float64(l)
	avg := formula.Average(f[:l]...)
	d = append(d, avg)
	for i := l; i < len(f); i++ {
		avg = f[i]*alpha + (1-alpha)*avg
		d = append(d, avg)
	}
	s.name = "RMA"
	s.Avg = avg // For live
	s.data = d
	return s
}

type VWMA struct {
	ERS[float64]
}

func Vwma(src, volume Series, l int) Series {
	var s VWMA
	s.res = src.Resolution()
	s.st = src.StartTime() + s.res*int64(l)
	f := src.Data()
	d := make([]float64, 0, len(f))
	vol := volume.Data()
	vol = vol[len(vol)-len(f):]
	s.name = "VMWA"
	volSum := formula.Sum(vol[:l]...)
	volXsrcSum := formula.Sum(formula.ArrayOperation(formula.Mul[float64], vol[:l], f[:l])...)
	avg := volXsrcSum / volSum
	d = append(d, avg)
	for i := l; i < len(f); i++ {
		volSum = volSum + vol[i] - vol[i-l]
		volXsrcSum = volXsrcSum + vol[i]*f[i] - vol[i-l]*f[i-l]
		d = append(d, volXsrcSum/volSum)
	}
	s.data = d

	return &s
}

type WMA struct {
	ERS[float64]
}

func Wma(src Series, l int) Series {
	s := new(WMA)
	s.res = src.Resolution()
	s.st = src.StartTime() + s.res*int64(l)

	lf := float64(l)
	alpha := lf * (lf + 1) / 2.0
	s.name = "WMA"
	f := src.Data()
	d := make([]float64, 0, len(f)-l+1)
	for i := l; i < len(f); i++ {
		var sum float64
		for y := 1; y < l+1; y++ {
			sum += f[i-l+y] * float64(y)
		}
		d = append(d, sum/alpha)
	}
	s.data = d
	return s
}

func DoubleEma(src Series, l int) Series {
	return DoubleMA(Ema, src, l)
}

func DoubleSma(src Series, l int) Series {
	return DoubleMA(Sma, src, l)
}
func DoubleWma(src Series, l int) Series {
	return DoubleMA(Wma, src, l)
}
