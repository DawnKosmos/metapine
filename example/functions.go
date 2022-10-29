package main

import . "github.com/DawnKosmos/metapine/backend/series/ta"

func Saphir(src1 Series, src2 Series, volume Series, len1, len2, len3 int) (buy Condition, sell Condition) {
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

/*
maR = input(title="SMI 1", defval=12)
ma1 = input(title="SMI  2", defval=4)
outR1 = sma(rsi(haOpen + haClose, maR), 2)
outR2 = sma(rsi(hahaHigh+hahaLow,maR),2)
outR = (2*outR2+outR1)/3
outB1 = sma(outR, ma1)
outB2 = sma(outB1, ma1)
outB = 2 * outB1 - 1 * outB2
loa = input(2)
cc = outR - outB
vwma_1 = (cc*haVolume*loa+cc[1]*haVolume[1])/(haVolume*loa+haVolume[1])
sma_1 = sma(cc, 2)
ccc = useVol ? vwma_1 : sma_1
c1 = ccc > ccc[1] and ccc[1] < ccc[2]
c2 = ccc < ccc[1] and ccc[1] > ccc[2]
*/
