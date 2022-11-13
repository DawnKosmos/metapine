package main

import . "github.com/DawnKosmos/metapine/backend/series/ta"

func Saphir(src1 Series, src2 Series, volume Series, len1, len2 int) (buy Condition, sell Condition) {
	outR1 := Sma(Rsi(src1, len1), 2)
	outR2 := Sma(Rsi(src2, len1), 2)
	outR := Div(AddF(outR2, outR1, 2), 3)
	outB1 := Sma(outR, len2)
	outB2 := Sma(outB1, len2)
	outB := SubF(outB1, outB2, 2)

	cc := Sub(outR, outB)
	c := Sma(cc, 2)

	c1, c2 := OffS(c, 1), OffS(c, 2)

	buy = And(Greater(c, c1), Smaller(c1, c2))
	sell = And(Smaller(c, c1), Greater(c1, c2))
	return
}

func dai(open, high, low, close, volume Series, l1, l2 int) Series {
	sh := Vwma(Sum(Add(high, close), l1), volume, l2)
	lo := Vwma(Sum(Sub(close, low), l1), volume, l2)
	return Div(Sub(lo, sh), sh)
}

func longCon(prozent float64, len int, close Series, low Series) Condition {
	lowest := Lowest(low, len)
	r1 := Sub(Div(close, lowest), 1)
	return Greater(r1, prozent/100)
}

func shortCon(prozent float64, len int, close, high Series) Condition {
	highest := Highest(high, len)
	r1 := Sub(Div(close, highest), 1)
	return Smaller(r1, -prozent/100)
}

func saphir(src1, src2, volume Series, l1, l2 int) (saphir Series) {
	r1 := Sma(Rsi(src1, l1), 2) //
	r2 := Sma(Rsi(src2, l2), 2)
	r := AddF(r1, r2, 2)
	b1 := Sma(r, l2)
	b2 := Sma(b1, l2)
	b := SubF(b1, b2, 2)
	cc := Sub(r, b)
	if volume == nil {
		saphir = Sma(cc, 2)
	} else {
		cc1, volume1 := OffS(cc, 1), OffS(volume, 1)
		vwma1 := Add((MultF(cc, volume, 2)), Mult(cc1, volume1))
		saphir = Div(vwma1, Add(Mult(volume, 2), volume1))
	}
	return
}

func MaCross(maFast func(s1 Series, len int) Series, maSlow func(s1 Series, len int) Series, src Series, fast int, slow int) (buy Condition, sell Condition) {
	fastMa := maFast(src, fast)
	slowMa := maSlow(src, slow)

	buy = Crossover(fastMa, slowMa)
	sell = Crossunder(fastMa, slowMa)
	return
}

//To iterate an Indicator you have to implement the Iterator interface
//Lets look at following indicator

func maCross(maFast func(s1 Series, len int) Series, maSlow func(s1 Series, len int) Series, src Series, fast int, slow int) (buy Condition, sell Condition) {
	fastMa := maFast(src, fast)
	slowMa := maSlow(src, slow)

	buy = Crossover(fastMa, slowMa)
	sell = Crossunder(fastMa, slowMa)
	return
}

func SolApeIter(oc2 Series, volume Series, ma func(s Series, l int) Series, len1, len2 int) (Condition, Condition) {
	outR := Sma(Roc(oc2, len1), 2)
	outB1 := ma(outR, len2)
	outB2 := ma(outB1, len2)
	outB := SubF(outB1, outB2, 2.0)
	cc := Sub(outR, outB)
	var c1 Series
	if volume == nil {
		c1 = Sma(cc, 2)
	} else {
		c1 = Vwma(cc, volume, 2)
	}
	c2, c3 := OffS(c1, 1), OffS(c1, 2)
	buy := And(Greater(c1, c2), Smaller(c2, c3))
	sell := And(Smaller(c1, c2), Greater(c2, c3))
	return buy, sell
}
