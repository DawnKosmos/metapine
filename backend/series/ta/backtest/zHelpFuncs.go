package backtest

import "github.com/DawnKosmos/metapine/backend/exchange"

/*
SafeFloat
Pretty often indicator have no Value while others have it. To still store everything in a File we have this struct.
SafeFloat is used when we use filters for our Strategies
*/
type SafeFloat struct {
	Safe  bool
	Value float64
}

// A better representation of the PNL overtime
type CandlePNL struct {
	Open  float64
	High  float64
	Close float64
	Low   float64
}

func PNLCalcCandle(Side bool, realisedPnl, netSize, avgEntry float64, c exchange.Candle) CandlePNL {

	if Side {
		return CandlePNL{
			Open:  realisedPnl + netSize*PnlCalc(avgEntry, c.Open),
			High:  realisedPnl + netSize*PnlCalc(avgEntry, c.High),
			Close: realisedPnl + netSize*PnlCalc(avgEntry, c.Close),
			Low:   realisedPnl + netSize*PnlCalc(avgEntry, c.Low),
		}
	} else {
		return CandlePNL{
			Open:  realisedPnl - netSize*PnlCalc(avgEntry, c.Open),
			High:  realisedPnl - netSize*PnlCalc(avgEntry, c.High),
			Close: realisedPnl - netSize*PnlCalc(avgEntry, c.Close),
			Low:   realisedPnl - netSize*PnlCalc(avgEntry, c.Low),
		}
	}
}

/*

		10   -8    +10
		100 120 130 140 150
realised	160 160 160 160
unrealized	40	60  100 200

entry	100			120
calc				240
*/

func PnlCalc(entry float64, exit float64) float64 {
	return (exit - entry) / entry
}
