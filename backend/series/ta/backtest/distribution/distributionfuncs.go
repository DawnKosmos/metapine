package distribution

type Func func(price, min, max float64, orderCount int) [][2]float64

func Normal(mp, min, max float64, div int) [][2]float64 {
	out := make([][2]float64, 0, div)
	size := 1.0 / float64(div)
	minp := mp - mp*(min/100)
	maxp := mp - mp*(max/100)
	b := (minp - maxp) / float64(div-1)
	out = append(out, [2]float64{size, minp})

	for i := 0; i < div-1; i++ {
		minp = minp - b
		out = append(out, [2]float64{size, minp})
	}
	return out
}

func Exponential(mp, min, max float64, division int) [][2]float64 {
	div := float64(division)
	out := make([][2]float64, 0, division)
	minp := mp - mp*(min/100)
	maxp := mp - mp*(max/100)
	b := (minp - maxp) / (div - 1)

	sum := (div + 1) / 2
	fn := func(iterate int) float64 {
		return (float64(iterate+1) / div) / sum
	}
	for i := 0; i < division; i++ {
		out = append(out, [2]float64{fn(i), minp})
		minp = minp - b
	}

	return out
}

func OpenAndDipsNormal(openSize float64) func(mp, min, max float64, div int) [][2]float64 {
	b := func(mp, min, max float64, division int) [][2]float64 {
		var os float64 = openSize
		//div := float64(division)
		out := make([][2]float64, 0, division+1)

		out = append(out, [2]float64{0.5, mp})
		kek := Normal(mp, min, max, division)

		rz := 1 - os // rest size
		for _, v := range kek {
			v[0] = v[0] * rz
			out = append(out, v)
		}
		return out
	}
	return b
}

func OpenAndDipsExponential(openSize float64) func(mp, min, max float64, div int) [][2]float64 {
	b := func(mp, min, max float64, division int) [][2]float64 {
		var os float64 = openSize
		//div := float64(division)
		out := make([][2]float64, 0, division+1)
		out = append(out, [2]float64{0.5, mp})
		kek := Exponential(mp, min, max, division)

		rz := 1 - os // rest size
		for _, v := range kek {
			v[0] = v[0] * rz
			out = append(out, v)
		}
		return out
	}
	return b
}
