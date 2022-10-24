package backtest

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
)

type TEInfo struct {
	Name             string
	Info             string
	CandlePnlSupport bool
}

type TradeExecution interface {
	CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat) (*Trade, error) //TradeExecution defines the strategy and gets as input an array from trade start to end
	GetInfo() TEInfo
}

func DefaultFeeInfo() *FeeInfo {
	return new(FeeInfo)
}

type ScaledLimit struct {
	Min          float64
	Max          float64
	OrderCount   int
	Size         *Size
	Distribution distribution.Func
	Parameter    *BacktestParameters
	stopLong     ta.Series // A stop is just a number.
	stopShort    ta.Series // A stop is just a number
	stopSize     float64
}

func (s *ScaledLimit) CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat) (*Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same candle")
	}

	var t *Trade
	var n, nMax int = 0, s.OrderCount
	if Side {
		dist := s.Distribution(ch[0].Open, s.Min, s.Max, s.OrderCount)
		for i := 0; i < exitCandle; {
			if nMax == n {
				t.PnlCalculation(ch[i])
				continue
			}
			if dist[n][1] > ch[i].Low {
				size := dist[n][0] * s.Size.Val
				var f Fill = Fill{
					Side:  Side,
					Type:  LIMIT,
					Price: dist[n][1],
					Size:  size,
					Time:  ch[i].StartTime,
					Fee:   dist[n][1] * s.Parameter.Fee.Maker,
				}
				if t == nil {
					t = NewTrade(f)
				} else {
					t.Add(f)
				}
				n++
			} else {
				t.PnlCalculation(ch[i])
				i++
			}
		}
	} else {
		dist := s.Distribution(ch[0].Open, -s.Min, -s.Max, s.OrderCount)
		for i := 0; i < exitCandle; {
			if nMax == n {
				t.PnlCalculation(ch[i])
				continue
			}
			if dist[n][1] > ch[i].Low {
				size := dist[n][0] * s.Size.Val
				var f Fill = Fill{
					Side:  Side,
					Type:  LIMIT,
					Price: dist[n][1],
					Size:  size,
					Time:  ch[i].StartTime,
					Fee:   dist[n][1] * s.Parameter.Fee.Maker,
				}
				if t == nil {
					t = NewTrade(f)
				} else {
					t.Add(f)
				}
				n++
			} else {
				t.PnlCalculation(ch[i])
				i++
			}
		}
	}

	if t == nil {
		return nil, errors.New("No trades got filled")
	}
	t.Close(ch[exitCandle].Open, s.Parameter.Fee.Slippage, ch[exitCandle].StartTime, MARKET, s.Parameter.Fee.Maker)
	return t, nil
}

/*
Trade Get created after the Signal Candle closed
func (s *ScaledLimit) CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat) (*Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same candle")
	}
	mp := ch[0].Open
	var n, nMax int = 0, s.OrderCount //
	var avgPrice, totalSize float64
	var fills []Fill
	var totalFee float64

	if Side {
		dist := s.Distribution(mp, s.Min, s.Max, s.OrderCount)
		pnl := make([]float64, 0, exitCandle) // fees just get calculated when trade closes
		pCandle := make([]CandlePNL, 0, exitCandle)

		for i := 0; i < exitCandle; {
			if nMax == n {
				pCandle = append(pCandle, PNLCalcCandle(1+totalSize, avgPrice, ch[i]))
				pnl = append(pnl, pCandle[i].Close)
				continue
			}
			if dist[n][1] > ch[i].Low {
				Size := dist[n][0] * s.Size.Val
				fee := s.Parameter.Fee.Maker * Size
				f := Fill{
					Side:  true,
					Type:  LIMIT,
					Price: dist[n][1],
					Size:  Size,
					Time:  ch[i].StartTime,
					Fee:   fee,
				}
				totalFee += fee
				fills = append(fills, f)
				avgPrice = (avgPrice*totalSize + f.Size*f.Price) / (totalSize + f.Size)
				totalSize += f.Size
				n++
			} else {
				pCandle = append(pCandle, PNLCalcCandle(1+totalSize, avgPrice, ch[i]))
				pnl = append(pnl, pCandle[i].Close)
				i++
			}
		}
		if totalSize == 0 {
			return nil, errors.New("trades not filled")
		}

		fee := totalSize * s.Parameter.Fee.Taker
		fills = append(fills, Fill{
			Side:  false,
			Type:  MARKET,
			Price: ch[exitCandle].Open,
			Size:  totalSize,
			Time:  ch[exitCandle].StartTime,
			Fee:   totalSize * s.Parameter.Fee.Taker,
		})
		totalFee += fee

		t := &Trade{
			Side:            Side,
			AvgBuy:          avgPrice,
			AvgSell:         ch[exitCandle].Open - s.Parameter.Fee.Slippage,
			EntrySignalTime: ch[0].StartTime,
			CloseSignalTime: ch[exitCandle].StartTime,
			Fills:           fills,
			Size:            totalSize,
			Pnl:             nil,
			Indicator:       nil,
			Fee:             totalFee,
		}
		return t, nil
	}

}
*/

func (s *ScaledLimit) GetInfo() TEInfo {
	return TEInfo{
		Name:             "Limit Orders",
		Info:             fmt.Sprintf("Min: %s, Max: %s, Orders: %s, StopLong	: %v, StopShort: %v", s.Min, s.Max, s.OrderCount, s.stopLong != nil, s.stopShort != nil),
		CandlePnlSupport: true,
	}
}

func NewScaledLimit(s *ScaledLimit) TradeExecution {
	if s.Size == nil {
		s.Size = DefaultSize()
	}
	if s.Distribution == nil {
		s.Distribution = distribution.Normal
	}
	if s.Min > s.Max {
		return nil
	}
	if s.OrderCount < 2 {
		s.OrderCount = 2
	}

	if s.Parameter == nil {
		DefaultParameters()
	}
	if s.stopSize == 0 {
		s.stopSize = 1
	}
	return s
}

/*

//TradeExecutionLimit executes Limit Orders
//A Distribution function returns at which prices the limit order get filled
type TradeExecutionLimit struct {
	o            *series.OHCLV
	Min          float64                                                        //minimum percent
	Max          float64                                                        //maximum percent
	division     int                                                            //divided
	Size         float64                                                        //total Size 1.0 equals to 100%
	Distribution func(marketPrice, Min, Max float64, division int) [][2]float64 // function that distributes the price and Size of a fill, [2]float{Size,price}
	fee          float64
}

func NewTradeExecutionLimit(minValue, maxValue float64, division int, Size float64, fee float64, distributionFn func(mp, Min, Max float64, division int) [][2]float64) TradeExecution {
	return TradeExecutionLimit{
		o:            nil,
		Min:          minValue,
		Max:          maxValue,
		division:     division,
		Size:         Size,
		Distribution: distributionFn,
		fee:          fee,
	}
}

func (p TradeExecutionLimit) CreateTrade(Side Side, ch []exchange.Candle) (Trade, error) {
	if len(ch) < 2 {
		return Trade{}, errors.New("Same Candle")
	}

	mp := ch[0].Open
	lc := ch[len(ch)-1]

	var n int // pointer to the number of fill in the array
	var nmax int = p.division
	var totalVolume float64
	var Fills []Fill
	var t Trade
	avgEntry := 1.0

	if Side { //long
		dist := p.Distribution(mp, p.Min, p.Max, p.division)
		pnl := make([]float64, 0, len(ch)-1)
		for i := 0; i < len(ch)-1; { //iterate the []Candle and see which limitorder [][2]Float64 gets filled
			if nmax == n { //all orders got filled, nothing to do anymore, besides calculating the PNL
				pnl = append(pnl, 1+totalVolume*PnlCalc(avgEntry, ch[i].Close))
				i++
				continue
			}
			if dist[n][1] > ch[i].Low {
				Size := dist[n][0] * p.Size
				f := Fill{
					Side:  true,
					Price: dist[n][1],
					Size:  Size,
					Time:  ch[i].StartTime,
				}
				Fills = append(Fills, f)
				avgEntry = (avgEntry*totalVolume + f.Size*f.Price) / (totalVolume + f.Size)
				totalVolume += f.Size
				n++
			} else {
				pnl = append(pnl, 1+totalVolume*PnlCalc(avgEntry, ch[i].Close))
				i++
			}
		}
		if totalVolume == 0 { // check if trades got filled
			return Trade{}, errors.New("trades not filled")
		}

		//Fill t
		t = Trade{
			StoppedOut:      false,
			Side:            true,
			AvgBuy:        avgFill(Fills),
			AvgSell:         lc.Open,
			EntrySignalTime: ch[0].StartTime,
			ExitSignalTime:  lc.StartTime,
			Fills: append(Fills, Fill{
				Side:  false,
				Price: lc.Open,
				Size:  totalVolume,
				Time:  lc.StartTime,
			}),
			Size:      totalVolume,
			ignore:    false,
			Indicator: []float64{},
			Pnl:       pnl,
		}
	} else { //Short
		dist := p.Distribution(mp, -p.Min, -p.Max, p.division)
		pnl := make([]float64, 0, len(ch)-1)
		for i := 0; i < len(ch)-1; { //iterate the ch and see which limit gets filled
			if nmax == n {
				pnl = append(pnl, 1-totalVolume*PnlCalc(avgEntry, ch[i].Close))
				i++
				break
			}
			if dist[n][1] < ch[i].High {
				Size := dist[n][0] * p.Size
				f := Fill{
					Side:  false,
					Price: dist[n][1],
					Size:  Size,
					Time:  ch[i].StartTime,
				}

				Fills = append(Fills, f)
				avgEntry = (avgEntry*totalVolume + f.Size*f.Price) / (totalVolume + f.Size)
				totalVolume += f.Size
				n++
			} else {
				pnl = append(pnl, 1-totalVolume*PnlCalc(avgEntry, ch[i].Close))
				i++
			}
		}
		if totalVolume == 0 { // check if trades got filled
			return Trade{}, errors.New("trades not filled")
		}

		//Fill t
		t = Trade{
			StoppedOut:      false,
			Side:            true,
			AvgBuy:        avgFill(Fills),
			AvgSell:         lc.Open,
			EntrySignalTime: ch[0].StartTime,
			ExitSignalTime:  lc.StartTime,
			Fills: append(Fills, Fill{
				Side:  true,
				Price: lc.Open,
				Size:  totalVolume,
				Time:  lc.StartTime,
			}),
			Size:      totalVolume,
			ignore:    false,
			Indicator: []float64{},
		}
	}

	return t, nil
}
*/
