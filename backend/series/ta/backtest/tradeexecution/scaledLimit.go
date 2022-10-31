package tradeexecution

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	. "github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest/distribution"
)

type ScaledLimit struct {
	Min          float64
	Max          float64
	OrderCount   int
	size         float64
	distribution distribution.Func
	stopLong     ta.Series // A stop is just a number.
	stopShort    ta.Series // A stop is just a number
	stopSize     float64
}

func NewScaledLimit(min float64, max float64, OrderCount int) *ScaledLimit {
	return &ScaledLimit{
		Min:          min,
		Max:          max,
		OrderCount:   OrderCount,
		size:         1,
		distribution: distribution.Normal,
		stopLong:     nil,
		stopShort:    nil,
		stopSize:     1,
	}
}

func (s *ScaledLimit) CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []SafeFloat, sizeInUsd float64, fee FeeInfo, pnlgraph bool) (*Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same candle")
	}

	t := EmptyTrade(Side, ch[0].StartTime)
	t.Indicator = indicators

	var n, nMax int = 0, s.OrderCount
	if Side {
		dist := s.distribution(ch[0].Open, s.Min, s.Max, s.OrderCount)
		for i := 0; i < exitCandle; {
			if nMax == n {
				if pnlgraph {
					t.PnlCalculation(ch[i])
				}
				i++
				continue
			}
			if dist[n][1] > ch[i].Low {
				fillSize := dist[n][0] * s.size * sizeInUsd
				t.Add(Fill{
					Side:  Side,
					Type:  LIMIT,
					Price: dist[n][1],
					Size:  fillSize / dist[n][1],
					Time:  ch[i].StartTime,
					Fee:   fillSize * fee.Maker,
				})
				n++
			} else {
				if pnlgraph {
					t.PnlCalculation(ch[i])
				}
				i++
			}
		}
	} else {
		dist := s.distribution(ch[0].Open, -s.Min, -s.Max, s.OrderCount)
		for i := 0; i < exitCandle; {
			if nMax == n {
				if pnlgraph {
					t.PnlCalculation(ch[i])
				}
				i++
				continue
			}
			if dist[n][1] < ch[i].High {
				fillSize := dist[n][0] * s.size * sizeInUsd
				t.Add(Fill{
					Side:  Side,
					Type:  LIMIT,
					Price: dist[n][1],
					Size:  fillSize / dist[n][1],
					Time:  ch[i].StartTime,
					Fee:   fillSize * fee.Maker,
				})
				n++
			} else {
				if pnlgraph {
					t.PnlCalculation(ch[i])
				}
				i++
			}
		}
	}

	if len(t.Fills) == 0 {
		return nil, errors.New("No trades got filled")
	}
	t.Close(ch[exitCandle].Open, fee.Slippage, ch[exitCandle].StartTime, MARKET, fee.Taker)
	return t, nil
}

func (s *ScaledLimit) GetInfo() TEInfo {
	return TEInfo{
		Name:             "Limit Orders",
		Info:             fmt.Sprintf("Min: %f, Max: %f, Orders: %d, StopLong	: %v, StopShort: %v", s.Min, s.Max, s.OrderCount, s.stopLong != nil, s.stopShort != nil),
		CandlePnlSupport: true,
	}
}

// Setter Functions
func (s *ScaledLimit) Size(size float64) *ScaledLimit {
	s.size = size
	return s
}

func (s *ScaledLimit) Distribution(fn distribution.Func) *ScaledLimit {
	s.distribution = fn
	return s
}

func (s *ScaledLimit) Stop(long, short ta.Series, size float64) *ScaledLimit {
	s.stopLong, s.stopShort = long, short
	return s
}
