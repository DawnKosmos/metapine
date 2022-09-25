package ta

import "github.com/DawnKosmos/metapine/helper/formula"

type LOGICAL struct {
	ERS[bool]
}

func logical(op func(bool, bool) bool, con1, con2 Condition) Condition {
	s := new(LOGICAL)
	f, f1 := con1.Data(), con2.Data()
	s.res = con1.Resolution()
	minLen, position := formula.ShortestLenOfArray(f, f1)
	if position == 0 {
		s.st = con1.StartTime()
		f1 = f1[len(f1)-minLen:]
	} else {
		s.st = con2.StartTime()
		f = f[len(f)-minLen:]
	}
	d := make([]bool, 0, minLen)
	for i := 0; i < minLen; i++ {
		d = append(d, op(f[i], f1[i]))
	}
	s.data = d
	return s
}

func And(c1, c2 Condition) Condition {
	op := func(a, b bool) bool {
		return a && b
	}
	return logical(op, c1, c2)
}

func Or(c1, c2 Condition) Condition {
	op := func(a, b bool) bool {
		return a || b
	}
	return logical(op, c1, c2)
}

func Xor(c1, c2 Condition) Condition {
	op := func(x, y bool) bool {
		return (x || y) && !(x && y)
	}
	return logical(op, c1, c2)
}

type not struct {
	ERS[bool]
}

func Not(con Condition) Condition {
	s := new(not)
	s.st = con.StartTime()
	s.res = con.Resolution()
	f := con.Data()
	d := make([]bool, 0, len(f))
	for _, v := range f {
		d = append(d, !v)
	}
	s.data = d
	return s
}
