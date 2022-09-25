package ta

import "github.com/DawnKosmos/metapine/helper/formula"

type iff[T any] struct {
	ERS[T]
}

// IfS is user when the output is a Series | float64 value
func IfS(con Condition, Then interface{}, Else interface{}) Series {
	s := new(iff[float64])
	//Change Types
	var src1, src2 Series
	switch v := Then.(type) {
	case Series:
		src1 = v
	case float64:
		src1 = constantS(con, v)
	case int:
		src1 = constantS(con, float64(v))
	}
	switch v := Else.(type) {
	case Series:
		src2 = v
	case float64:
		src2 = constantS(con, v)
	case int:
		src2 = constantS(con, float64(v))
	}
	arr := []ResolutionStartTime{con, src1, src2}
	c1, f1, f2 := con.Data(), src1.Data(), src2.Data()
	l1, pos := formula.MinInt(len(c1), len(f1), len(f2))
	s.st = arr[pos].StartTime()
	s.res = con.Resolution()
	c1 = c1[len(c1)-l1:]
	f1 = f1[len(f1)-l1:]
	f2 = f2[len(f2)-l1:]
	d := make([]float64, 0, l1)
	for i, v := range c1 {
		if v {
			d = append(d, f1[i])
		} else {
			d = append(d, f2[i])
		}
	}
	s.data = d
	return s
}

// IfC is used when the output has to be a Condition|boolean value
func IfC(con Condition, Then interface{}, Else interface{}) Condition {
	s := new(iff[bool])
	//Change Types
	var src1, src2 Condition
	switch v := Then.(type) {
	case Condition:
		src1 = v
	case bool:
		src1 = constantB(con, v)
	}
	switch v := Else.(type) {
	case Condition:
		src2 = v
	case bool:
		src2 = constantB(con, v)
	}
	arr := []ResolutionStartTime{con, src1, src2}
	c1, f1, f2 := con.Data(), src1.Data(), src2.Data()
	l1, pos := formula.MinInt(len(c1), len(f1), len(f2))
	s.st = arr[pos].StartTime()
	s.res = con.Resolution()
	c1 = c1[len(c1)-l1:]
	f1 = f1[len(f1)-l1:]
	f2 = f2[len(f2)-l1:]
	d := make([]bool, 0, l1)
	for i, v := range c1 {
		if v {
			d = append(d, f1[i])
		} else {
			d = append(d, f2[i])
		}
	}
	s.data = d
	return s
}
