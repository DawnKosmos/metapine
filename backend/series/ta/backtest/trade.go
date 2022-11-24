package backtest

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"time"
)

type Trade struct {
	Side                             bool    // True = Long, False = Short
	AvgBuy, AvgSell, avgPrice        float64 //the AvgBuy/Sell is needed to calculate realisedPNL, avgPrice is used for PNLCandle
	EntrySignalTime, CloseSignalTime time.Time

	Fills                      []Fill
	BuySize, SellSize, NetSize float64
	//The PNL starts with the EntrySignalTime. Every Tick represents 1 Candle
	//This information is needed to calculate the Overall PNL of the Indicator
	Pnl       []float64
	Indicator []SafeFloat
	UsdVolume float64
	Fee       float64
}

func EmptyTrade(Side bool, EntrySignalTime time.Time) *Trade {
	return &Trade{
		Side:            Side,
		EntrySignalTime: EntrySignalTime,
		Pnl:             []float64{},
	}
}

func NewTrade(f Fill) *Trade {
	t := &Trade{
		Side:            f.Side,
		EntrySignalTime: f.Time,
		Pnl:             []float64{},
	}
	t.Add(f)
	return t
}

// PnlCalculation TODO unit test
func (t *Trade) PnlCalculation(c exchange.Candle) {
	var realisedPNL float64
	if t.SellSize != 0 && t.BuySize != 0 {
		if t.Side {
			realisedPNL = (t.AvgSell - t.AvgBuy) * (t.BuySize - t.NetSize)
		} else {
			realisedPNL = -(t.AvgBuy - t.AvgSell) * (t.SellSize + t.NetSize)
		}
	}

	pnl := realisedPNL + (c.Close-t.avgPrice)*t.NetSize

	t.Pnl = append(t.Pnl, pnl)
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

	t.UsdVolume += f.Size * f.Price

	t.Fee += f.Fee
	t.Fills = append(t.Fills, f)
}

func (t *Trade) addTooShort(f Fill) {
	if !f.Side {
		t.avgPrice = (f.Size*f.Price - t.avgPrice*t.NetSize) / (f.Size - t.NetSize)
		t.AvgSell = (t.AvgSell*t.SellSize + f.Size*f.Price) / (t.SellSize + f.Size)
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
	t.UsdVolume += f.Size * f.Price

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
		fprice := price - slippage
		f = Fill{
			Side:  false,
			Type:  feeType,
			Price: fprice,
			Size:  t.NetSize,
			Time:  close,
			Fee:   fee * t.NetSize * fprice,
		}
	} else {
		fprice := price + slippage
		f = Fill{
			Side:  true,
			Type:  feeType,
			Price: fprice,
			Size:  -t.NetSize,
			Time:  close,
			Fee:   fprice * fee * -t.NetSize,
		}

	}

	t.CloseSignalTime = close

	t.Add(f)
}

func (t *Trade) Start() time.Time {
	return t.EntrySignalTime
}

func (t *Trade) RealisedPNL() float64 {
	var realisedPNL float64
	if t.Side {
		realisedPNL = (t.AvgSell - t.AvgBuy) * (t.BuySize - t.NetSize)
	} else {
		realisedPNL = -(t.AvgBuy - t.AvgSell) * (t.SellSize + t.NetSize)
	}

	return realisedPNL - t.Fee
}

// SimpleTrade Or SimpleTrade is used For FastBacktesting, this mode is used to iterate many parameters.
// Only the results are safed to Calculate
type SimpleTrade struct {
	Side                bool
	Entry, Exit         float64
	EntryTime, ExitTime int64
	indicators          []float64
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
	return 1 + x - (fee * x)
}

// TRADES
type Trades []*Trade

func (t Trades) Less(i, j int) bool {
	return t[i].EntrySignalTime.Unix() < t[j].EntrySignalTime.Unix()
}

func (t Trades) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Trades) Len() int {
	return len(t)
}
