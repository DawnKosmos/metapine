package ta

// Tsi is the true strenght indicator
func Tsi(src Series, r, s int) Series {
	src1 := OffS(src, 1)
	m := Sub(src, src1)

	t1 := Ema(Ema(m, r), s)
	t2 := Ema(Ema(Abs(m), r), s)
	return DivF(t1, t2, 100)
}

func BollingerBands(src Series, len int, mul float64) (LowerBand, Basis, UpperBand Series) {
	Basis = Sma(src, len)
	std := Mult(Stdev(src, len), mul)
	LowerBand = Sub(Basis, std)
	UpperBand = Add(Basis, std)
	return
}

func BollingerBandsWidth(src Series, len int, mul float64) Series {
	l, b, u := BollingerBands(src, len, mul)
	return Div(Sub(u, l), b)
}

// Macd is the equivalent of macd(source, fastLenght, slowLenght, signalLenght). Returns the macd, signal, histogram
func Macd(src Series, fastLen, slowLen, SignalLen int) (macd Series, signal Series, histogram Series) {
	f := Ema(src, fastLen)
	s := Ema(src, slowLen)
	// macd = f - s
	macd = Sub(f, s)
	signal = Ema(macd, SignalLen)
	//histogram macd - signal
	histogram = Sub(macd, signal)
	return
}

func MacdRelative(src Series, fastLen, slowLen, SignalLen int) (macd Series, signal Series, histogram Series) {
	f := Ema(src, fastLen)
	s := Ema(src, slowLen)
	// macd = f - s
	macd = DivF(Sub(f, s), s, 100)
	signal = Ema(macd, SignalLen)
	//histogram macd - signal
	histogram = Sub(macd, signal)
	return
}

/*
TrendRibonNoro, The link for this script is posted in the src Code
Source Code:
https://www.tradingview.com/script/ZsKsLiUU-noro-s-trend-ribbon-strategy/
*/
func TrendRibonNoro(MAFunction func(src Series, len int) Series, src Series, len int) (lowerLine, upperLine Series) {
	ma := MAFunction(src, len)
	upperLine = Highest(ma, len)
	lowerLine = Lowest(ma, len)
	return
}

func Momentum(src Series, l int) Series {
	return Sub(src, OffS(src, l))
}

func Range(src Series, l int) Series {
	return Sub(Highest(src, l), Lowest(src, l))
}

func WilliamsR(close, high, low Series, l int) Series {
	a := Sub(Highest(high, l), close)
	b := Sub(Highest(high, l), Lowest(low, l))
	return Div(a, b)
}

func MFI(src Series, volume Series, len int) Series {
	ch := Change(src, 1)
	con := SmallerEqual(ch, 0.0)
	upper := Sum(Mult(volume, IfS(con, 0, src)), len)
	lower := Sum(Mult(volume, IfS(Not(con), 0, src)), len)
	mfr := Div(upper, lower)
	mfi := DivF(mfr, Add(mfr, 1), 100)
	mfi.SetName("MFI")
	return mfi
}

// Atr gets an MA function, and a TR and len
func Atr(ma func(Series, int) Series, tr *TRange, l int) Series {
	return ma(tr, l)
}

// DoubleMA returns the double smoothed version of an MA
func DoubleMA(op func(Series, int) Series, src Series, l int) Series {
	e1 := op(src, l)
	return SubF(e1, op(e1, l), 2)
}

func SmoothedStoch(high, low, close Series, signalFunc func(s Series, l int) Series, l1, l2, l3, lenSignal int) (d, slowD, signal Series) {
	highS1 := Ema(high, 3)
	lowS1 := Ema(low, 3)

	highS := SubF(highS1, Ema(highS1, 3), 2)
	lowS := SubF(lowS1, Ema(lowS1, 3), 2)

	k := Stoch(close, highS, lowS, l1)
	d = Sma(k, l2)
	slowD = Sma(d, l3)

	signal = signalFunc(slowD, lenSignal)
	return
}

func Ribbon(src1, src2 Series, maFunc func(s Series, i int) Series, lenMa int, lenHL int) (Condition, Condition) {
	ma := maFunc(src1, lenMa)
	h := Highest(ma, lenHL)
	low := Lowest(ma, lenHL)

	h1 := OffS(h, 1)
	l1 := OffS(low, 1)

	trendT := IfS(Greater(src2, h1), 1, -1)

	trend := IfS(Equal(trendT, -1), IfS(Smaller(src2, l1), -1, 1), 1)

	trend1 := OffS(trend, 1)

	buy := And(Equal(trend, 1), Equal(trend1, -1))
	sell := And(Equal(trend, -1), Equal(trend1, 1))
	return buy, sell
}
