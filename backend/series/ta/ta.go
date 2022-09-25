package ta

//Tsi is the true strenght indicator
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

//Macd is the equivalent of macd(source, fastLenght, slowLenght, signalLenght). Returns the macd, signal, histogram
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
	return mfi
}

//Atr gets an MA function, and a TR and len
func Atr(ma func(Series, int) Series, tr *TRange, l int) Series {
	return ma(tr, l)
}

//DoubleMA returns the double smoothed version of an MA
func DoubleMA(op func(Series, int) Series, src Series, l int) Series {
	e1 := op(src, l)
	return SubF(e1, op(e1, l), 2)
}
