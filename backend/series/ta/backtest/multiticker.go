package backtest

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"time"
)

type MultiTicker struct {
	Name       string
	Algo       func(ch ta.Chart) (buy ta.Condition, sell ta.Condition)
	TE         TradeExecution
	Indicator  []func(ch ta.Chart) ta.Series
	Parameters BTParameter
}

func NewMultiTicker(Name string, Algo func(ch ta.Chart) (ta.Condition, ta.Condition), te TradeExecution, parameters BTParameter) *MultiTicker {
	return &MultiTicker{
		Name:       Name,
		Algo:       Algo,
		TE:         te,
		Parameters: parameters,
	}
}

func (s *MultiTicker) AddIndicator(indicator ...func(ch ta.Chart) ta.Series) {
	s.Indicator = indicator
}

func (s *MultiTicker) CreateResult(tickers []string, ee exchange.CandleProvider, st, et time.Time, res int64) []*BackTestStrategy {
	var bb []*BackTestStrategy
	for _, v := range tickers {
		ch := ta.NewOHCLV(ee, v, st, et, res)
		var indis []ta.Series
		for _, vv := range s.Indicator {
			indis = append(indis, vv(ch))
		}

		ss := NewStrategy(ch).AddIndicator(indis...)
		buy, sell := s.Algo(ch)
		bb = append(bb, ss.CreateStrategy(v, buy, sell, s.TE, s.Parameters))
	}
	return bb
}
