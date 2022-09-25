package ta

import "github.com/DawnKosmos/metapine/helper/formula"

type comp struct {
	ERS[bool]
	constant bool
	c        float64
}

func Comp(op func(float64, float64) bool, src Series, v interface{}) Condition {
	var s comp
	s.res = src.Resolution()
	var d []bool

	switch v := v.(type) {
	case int:
		s.st = src.StartTime()
		d = make([]bool, 0, len(src.Data()))
		for _, vv := range src.Data() {
			d = append(d, op(vv, float64(v)))
		}
	case float64:
		s.st = src.StartTime()
		d = make([]bool, 0, len(src.Data()))
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
		panic("Comparison, not valid type")
	}
	s.data = d
	return &s
}

// Smaller (src,v) => src < v
func Smaller(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 < v2
	}
	return Comp(o, src, v)
}

// Greater (src,v) => src > v
func Greater(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 > v2
	}
	return Comp(o, src, v)
}

// Equal (src,v) => src == v
func Equal(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 == v2
	}
	return Comp(o, src, v)
}

// NotEqual (src,v) => src != v
func NotEqual(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 != v2
	}
	return Comp(o, src, v)
}

// SmallerEqual (src,v) => src <= v
func SmallerEqual(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 <= v2
	}
	return Comp(o, src, v)
}

// GreaterEqual (src,v) => src >= v
func GreaterEqual(src Series, v interface{}) Condition {
	o := func(v1 float64, v2 float64) bool {
		return v1 >= v2
	}
	return Comp(o, src, v)
}
