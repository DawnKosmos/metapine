package chart

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/series/ta"
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
)

type Charter interface {
	WriteChart(e exchange.Candle) error
	AddPNL(graph backtest.PNLGraph) error
	AddIndicator(indicator ta.Series) error
}
