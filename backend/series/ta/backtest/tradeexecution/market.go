package tradeexecution

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"time"
)

type Market struct {
	size                float64
	stopLong, stopShort ta.Series
}

func NewMarketOrder(size float64) *Market {
	return &Market{size: size}
}

func (m *Market) Stop(long, short ta.Series) *Market {
	m.stopLong = long
	m.stopShort = short
	return m
}

func (m *Market) CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []backtest.SafeFloat, sizeInUsd float64, fee backtest.Fee, pnlgraph bool) (*backtest.Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same candle")
	}

	if Side {
		fillSize := m.size * sizeInUsd
		t := backtest.NewTrade(backtest.Fill{
			Side:  Side,
			Type:  backtest.MARKET,
			Price: ch[0].Open + fee.Slippage,
			Size:  fillSize / (ch[0].Open + fee.Slippage),
			Time:  time.Time{},
			Fee:   fillSize * fee.Taker,
		})
		t.EntrySignalTime = ch[0].StartTime
		t.Indicator = indicators
		t.Close(ch[exitCandle].Open, fee.Slippage, ch[exitCandle].StartTime, backtest.MARKET, fee.Maker)
		return t, nil
	} else {
		fillSize := m.size * sizeInUsd
		t := backtest.NewTrade(backtest.Fill{
			Side:  Side,
			Type:  backtest.MARKET,
			Price: ch[0].Open - fee.Slippage,
			Size:  fillSize / (ch[0].Open - fee.Slippage),
			Time:  time.Time{},
			Fee:   fillSize * fee.Taker,
		})
		t.EntrySignalTime = ch[0].StartTime
		t.Indicator = indicators
		t.Close(ch[exitCandle].Open, fee.Slippage, ch[exitCandle].StartTime, backtest.MARKET, fee.Taker)
		return t, nil
	}

}

func (m *Market) GetInfo() backtest.TEInfo {
	return backtest.TEInfo{
		Name:             "Market Orders",
		Info:             "",
		CandlePnlSupport: true,
	}
}
