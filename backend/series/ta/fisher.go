package ta

import "math"

type fisher struct {
	ERS[float64]
}

/*
//Fisher Oscillator Calculations
//{

fish_h = ta.highest(i_addf_fish_src, i_addf_fish_l)
fish_l = ta.lowest(i_addf_fish_src, i_addf_fish_l)

fish_rounding(input) =>
    input > 0.99 ? 0.999 : input < -0.99 ? -0.999 : input

value = 0.0

var fisher_v = 0.0
fisher_v:=fish_rounding(0.66 * ((hl2 - fish_l) / math.max(fish_h - fish_l, 0.001) - 0.5) + 0.67 * nz(fisher_v[1]))

var fisher_vf = 0.0
fisher_vf:= 0.5 * math.log((1 + fisher_v) / math.max(1 - fisher_v, 0.001)) + 0.5 * nz(fisher_vf[1])
fisher_vfp = fisher_vf[1]
*/

func Fisher(src Series, l int) Series {
	s := new(fisher)
	s.res, s.st = src.Resolution(), src.StartTime()+src.Resolution()*int64(l)
	f := src.Data()

	d := make([]float64, 0, len(f)-l+1)
	low := Lowest(src, l).Data()
	high := Highest(src, l).Data()
	value1 := fishCalc1(f[l-1], low[0], high[0], 0)

	fish := fishCalc(fishCalc2(value1), 0)
	d = append(d, fish)

	s.name = "Fisher"
	for i, v := range f[l:] {
		value1 = fishCalc1(v, low[i], high[i], value1)
		fish = fishCalc(fishCalc2(value1), fish)
		d = append(d, fish)
	}

	s.data = d
	return s
}

func fishCalc1(hl2, lowest, highest, prev float64) float64 {
	return 0.66*((hl2-lowest)/(highest-lowest)-0.5) + 0.67*prev
}

func fishCalc2(value1 float64) (v float64) {
	if value1 > 0.999 {
		v = 0.9999
	} else {
		if value1 < -0.999 {
			v = -0.9999
		} else {
			v = value1
		}
	}
	return
}

func fishCalc(value2 float64, fish float64) float64 {
	return 0.5*math.Log((1+value2)/(1-value2)) + 0.5*fish
}
