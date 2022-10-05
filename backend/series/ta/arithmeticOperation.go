package ta

import (
	"log"
	"math"
	"github.com/DawnKosmos/metapine/helper/formula"
)

type arit struct {
	ERS[float64]
}

func Arit(op func(float64, float64) float64, src Series, v interface{}) Series {
	s := new(arit)
	s.res = src.Resolution()
	var d []float64
	switch v := v.(type) {
	case float64:
		s.st = src.StartTime()
		d = make([]float64, 0, len(src.Data()))
		for _, vv := range src.Data() {
			d = append(d, op(vv, v))
		}
	case int:
		s.st = src.StartTime()
		d = make([]float64, 0, len(src.Data()))
		for _, vv := range src.Data() {
			d = append(d, op(vv, float64(v)))
		}
	case Series:
		var f, f1 []float64 = src.Data(), v.Data()
		l, pos := formula.ShortestLenOfArray(f, f1)
		if pos == 0 {
			s.st = src.StartTime()
			f1 = f1[len(f1)-l:]
		} else {
			s.st = v.StartTime()
			f = f[len(f)-l:]
		}
		for i := 0; i < len(f); i++ {
			d = append(d, op(f[i], f1[i]))
		}
	default:
		log.Panicln("wrong type", v)
	}
	s.data = d
	return s
}

//Abs return the Absolute of its interface{
func Abs(src Series) Series {
	o := func(v1 float64, _ float64) float64 {
		return math.Abs(v1)
	}
	return Arit(o, src, 0)
}

//Add (a,b) => a + b
func Add(src Series, v interface{}) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 + v2
	}
	return Arit(o, src, v)
}

//Mult (a,b) => a*b
func Mult(src Series, v interface{}) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 * v2
	}
	return Arit(o, src, v)
}

//Sub (src,v) => src - v
func Sub(src interface{}, v interface{}) Series {
	s, ok := src.(Series)

	if ok {
		o := func(v1 float64, v2 float64) float64 {
			return v1 - v2
		}
		return Arit(o, s, v)
	}

	f, ok := src.(float64)

	if ok {
		v, ok := src.(Series)
		if ok {
			o := func(v1 float64, v2 float64) float64 {
				return v2 - v1
			}
			return Arit(o, v, f)
		}
	}

	kek, ok := src.(int)
	if ok {
		v, ok := src.(Series)
		if ok {
			o := func(v1 float64, v2 float64) float64 {
				return v2 - v1
			}
			return Arit(o, v, kek)
		}
	}

	return nil

}

//Div (src,v) => src/v
func Div(src Series, v interface{}) Series {
	o := func(v1 float64, v2 float64) float64 {
		if v2 == 0 {
			v1 = 0
			v2 = 1
		}
		return v1 / v2
	}
	return Arit(o, src, v)
}

//Mod (src,v) => src%v
func Mod(src Series, v interface{}) Series {
	o := math.Mod
	return Arit(o, src, v)
}

//Pow (src,v) => src^v
func Pow(src Series, v interface{}) Series {
	o := math.Pow
	return Arit(o, src, v)
}

//Round (src) => round(src)
func Round(src Series) Series {
	o := func(v1 float64, _ float64) float64 {
		return math.Round(v1)
	}
	return Arit(o, src, 0.0)
}

//Min (src,v) => min(src,v)
func Min(src Series, v interface{}) Series {
	o := math.Min /*func(v1 float64, v2 float64) float64 {
		return math.Min(v1, v2)
	}*/
	return Arit(o, src, v)
}

//Max (src,v) => max(src,v)
func Max(src Series, v interface{}) Series {
	o := math.Max
	return Arit(o, src, v)
}

//Remainder (src,v) => return remainder of src/v
func Remainder(src Series, v interface{}) Series {
	o := math.Remainder
	return Arit(o, src, v)
}

//Hypot (src, v) => sqrt(src²+v²)
func Hypot(src Series, v interface{}) Series {
	o := math.Hypot
	return Arit(o, src, v)
}

//AddF (src, v, factor) => src*factor + v
func AddF(src Series, v interface{}, factor float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1*factor + v2
	}
	return Arit(o, src, v)
}

//SubF (src,v,factor) => src*factor - v
func SubF(src Series, v interface{}, factor float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1*factor - v2
	}
	return Arit(o, src, v)
}

//MultF (src,v,factor) => src*factor * v
func MultF(src Series, v interface{}, factor float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 * factor * v2
	}
	return Arit(o, src, v)
}

//DivF (src,v,factor) => src*factor/v
func DivF(src Series, v interface{}, factor float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 * factor / v2
	}
	return Arit(o, src, v)
}

//AddC (src,v,constant) => src+v+constant
func AddC(src Series, v interface{}, constant float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 + v2 + constant
	}
	return Arit(o, src, v)
}

//SubC (src,v,constant) => src-v+constant
func SubC(src Series, v interface{}, constant float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1 - v2 + constant
	}
	return Arit(o, src, v)
}

//DivC (src,v,constant) => src/v+constant
func DivC(src Series, v interface{}, constant float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1/v2 + constant
	}
	return Arit(o, src, v)
}

//MultC (src,v,constant) => src*v+constant
func MultC(src Series, v interface{}, constant float64) Series {
	o := func(v1 float64, v2 float64) float64 {
		return v1*v2 + constant
	}
	return Arit(o, src, v)
}

//Sqrt (src) => sqrt(src)
func Sqrt(src Series) Series {
	o := func(v1, v2 float64) float64 {
		return math.Sqrt(v1)
	}
	return Arit(o, src, 0)
}

//Neg(src) = -src
func Neg(src Series) Series {
	o := func(v1, _ float64) float64 {
		return -v1
	}
	return Arit(o, src, 0)
}

//Log returns the natural logarithm
func Log(src Series) Series {
	o := func(v1, _ float64) float64 {
		return math.Log(v1)
	}
	return Arit(o, src, 0)
}

//Sin (src) => sin(src)
func Sin(src Series) Series {
	o := func(v1, v2 float64) float64 {
		return math.Sin(v1)
	}
	return Arit(o, src, 0)
}

//Cos returns the Cos
func Cos(src Series) Series {
	o := func(v1, v2 float64) float64 {
		return math.Cos(v1)
	}
	return Arit(o, src, 0)
}

//Asin return the asin
func Asin(src Series) Series {
	o := func(v1, v2 float64) float64 {
		return math.Asin(v1)
	}
	return Arit(o, src, 0)
}

//Acos return the acos
func Acos(src Series) Series {
	o := func(v1, _ float64) float64 {
		return math.Acos(v1)
	}
	return Arit(o, src, 0)
}

//Floor rounds a float number to the lower or equal integer
func Floor(src Series) Series {
	o := func(v1, _ float64) float64 {
		return math.Floor(v1)
	}
	return Arit(o, src, 0)
}

//The Change between the older and newer value
func Change(src Series, len int) Series {
	srcOffset := OffS(src, len)
	return Arit(formula.Change[float64], srcOffset, src)
}

//Invert(src) => 1/src
func Invert(src Series) Series {
	o := func(v1, _ float64) float64 {
		return 1 / v1
	}
	return Arit(o, src, 0)
}
