package tradeexecution

import (
	"errors"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
)

type MarketMaker struct {
	Spread  float64
	Size    float64
	MaxSize float64
}

func (s *MarketMaker) CreateTrade(Side bool, ch []exchange.Candle, exitCandle int, indicators []backtest.SafeFloat, sizeInUsd float64, fee backtest.Fee) (*backtest.Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same Candle")
	}

	t := backtest.EmptyTrade(Side, ch[0].StartTime)

	//mp := ch[0].Open
	if Side {
		//	targetBuy := mp - mp*(s.Spread/100)
		for i := 0; i < exitCandle; {

		}
	}
	return t, nil

}

func (s *MarketMaker) GetInfo() backtest.TEInfo {
	//TODO implement me
	panic("implement me")
}

func mmdistrubution(mp, size, maxsize float64) [][2]float64 {
	var as float64
	for as > maxsize {

	}
	return nil
}
