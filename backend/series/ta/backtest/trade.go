package backtest

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"time"
)

type Trade struct {
	Side                             bool
	AvgBuy, AvgSell, avgPrice        float64 //the AvgBuy/Sell is needed to calculate realisedPNL, avgPrice is used for PNLCandle
	EntrySignalTime, CloseSignalTime time.Time

	Fills                      []Fill
	BuySize, SellSize, NetSize float64
	//The PNL starts with the EntrySignalTime. Every Tick represents 1 Candle
	//This information is needed to calculate the Overall PNL of the Indicator
	Pnl       []float64
	PnlCandle []CandlePNL
	Indicator []SafeFloat
	Fee       float64
}

func NewTrade(f Fill) *Trade {
	t := &Trade{
		Side:  f.Side,
		Fills: []Fill{f},
		Fee:   f.Fee,
	}
	t.avgPrice = f.Price

	if f.Side {
		t.AvgBuy = f.Price
		t.BuySize = f.Size
		t.NetSize = f.Size

	} else {
		t.AvgSell = f.Price
		t.SellSize = f.Size
		t.NetSize -= f.Size
	}
	return t
}

// PnlCalculation TODO unit test
func (t *Trade) PnlCalculation(c exchange.Candle) {
	var realisedPNL float64
	if t.Side {
		realisedPNL = 1 + (t.BuySize-t.NetSize)*(PnlCalc(t.AvgBuy, t.AvgSell)) - t.Fee
	} else {
		realisedPNL = 1 - (t.SellSize+t.NetSize)*(PnlCalc(t.AvgSell, t.AvgBuy)) - t.Fee
	}

	pp := PNLCalcCandle(t.Side, realisedPNL, t.NetSize, t.avgPrice, c)
	t.PnlCandle = append(t.PnlCandle, pp)
	t.Pnl = append(t.Pnl, pp.Close)
}

func (t *Trade) Add(f Fill) {
	if t.Side {
		t.addTooLong(f)
	} else {
		t.addTooShort(f)
	}
}

func (t *Trade) addTooLong(f Fill) {
	if f.Side {
		t.avgPrice = (t.avgPrice*t.NetSize + f.Size*f.Price) / (t.NetSize + f.Size)
		t.AvgBuy = (t.AvgBuy*t.BuySize + f.Size*f.Price) / (t.BuySize + f.Size)
		t.BuySize += f.Size
		t.NetSize += f.Size
	} else {
		if f.Size > t.NetSize {
			f.Size = t.NetSize
			f.Fee = f.Fee * (t.NetSize / f.Size)
		}
		t.AvgSell = t.AvgSell*t.SellSize + f.Size*f.Price/(t.SellSize+f.Size)
		t.NetSize -= f.Size
		t.SellSize += f.Size
	}

	t.Fee += f.Fee
	t.Fills = append(t.Fills, f)
}

func (t *Trade) addTooShort(f Fill) {
	if !f.Side {
		t.avgPrice = (f.Size*f.Price - t.avgPrice*t.NetSize) / (f.Size - t.NetSize)
		t.AvgSell = t.AvgSell*t.SellSize + f.Size*f.Price/(t.SellSize+f.Size)
		t.SellSize += f.Size
		t.NetSize -= f.Size
	} else {
		if f.Size > -t.NetSize {
			f.Fee = f.Fee * (-t.NetSize / f.Size)
			f.Size = -t.NetSize
		}
		t.AvgBuy = (t.AvgBuy*t.BuySize + f.Size*f.Price) / (t.BuySize + f.Size)
		t.BuySize += f.Size
		t.NetSize += f.Size
	}
	t.Fills = append(t.Fills, f)
	t.Fee += f.Fee
}

/*
Check Netsize cant be lower Smaller 0 in a Long and cant be greater 0 in a short
*/
func (t *Trade) Close(price float64, slippage float64, close time.Time, feeType FillType, fee float64) {
	var f Fill
	//Trade Fertigstellen
	if t.Side {
		f = Fill{
			Side:  false,
			Type:  feeType,
			Price: price - slippage,
			Size:  t.NetSize,
			Time:  close,
			Fee:   fee * t.NetSize,
		}
	} else {
		f = Fill{
			Side:  true,
			Type:  feeType,
			Price: price + slippage,
			Size:  -t.NetSize,
			Time:  close,
			Fee:   fee * -t.NetSize,
		}
	}
	t.Add(f)
}

// SimpleTrade Or SimpleTrade is used For FastBacktesting, this mode is used to iterate many parameters.
// Only the results are safed to Calculate
type SimpleTrade struct {
	Side                bool
	Entry, Exit         float64
	EntryTime, ExitTime int64
}

func CreateSimpleTrade(side bool, entry, exit exchange.Candle) (SimpleTrade, error) {
	if entry.StartTime.Unix() >= exit.StartTime.Unix() {
		return SimpleTrade{}, errors.New("Error Candle")
	}
	return SimpleTrade{
		Side:  side,
		Entry: entry.Open,
		Exit:  exit.Open,
	}, nil
}

func (t *SimpleTrade) Pnl(fee float64) float64 {
	var x float64
	if t.Side {
		x = (t.Exit - t.Entry) / t.Entry
	} else {
		x = -1 * (t.Exit - t.Entry) / t.Entry
	}
	return x - (fee * 0.01)
}
