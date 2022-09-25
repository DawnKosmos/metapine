package ta

import "github.com/DawnKosmos/metapine/helper/formula"

type AVG struct {
	ERS[float64]
}

func Avg(src ...Series) Series {
	s := new(AVG)

	var f [][]float64
	for _, v := range src {
		f = append(f, v.Data())
	}
	minLen, position := formula.ShortestLenOfArray(f)
	for i, v := range f {
		f[i] = v[len(v)-minLen:]
	}
	s.st, s.res = src[position].StartTime(), src[position].Resolution()
	d := make([]float64, 0, minLen)
	var sum float64
	for i := 0; i < minLen; i++ {
		sum = 0
		for j := 0; j < len(f); i++ {
			sum = sum + f[i][j]
		}
		d = append(d, sum/float64(len(f)))
	}
	s.data = d
	return s
}
