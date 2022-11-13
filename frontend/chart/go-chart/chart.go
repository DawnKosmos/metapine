package go_chart

import (
	"github.com/DawnKosmos/metapine/backend/series/ta/backtest"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func MultiLine(graph ...*backtest.PNLGraph) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "multi lines",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "shine",
			Width:  "1600px",
			Height: "900px",
		}),
	)
}
