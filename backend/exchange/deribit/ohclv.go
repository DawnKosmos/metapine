package deribit

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"log"
	"strconv"
	"time"
)

func (d *Deribit) OHCLV(ticker string, resolution int64, start time.Time, endTime time.Time) ([]exchange.Candle, error) {
	var hp []exchange.Candle
	res := resolution / 60
	st, et := start.UnixMilli(), endTime.UnixMilli()
	var end int64 = et
	if time.Now().UnixMilli() < end {
		end = time.Now().UnixMilli()
	}
	newRes := checkResolution(res)
	for {
		c, err := d.getOHCLV(ticker, newRes, st, end)
		if err != nil {
			fmt.Println(st, end)
			log.Printf("Error OHCLV FTX %v", err)
			return hp, err
		}
		if len(c) < 2 {
			fmt.Println(c)
			break
		}

		hp = append(c, hp...)
		end = hp[0].StartTime.UnixMilli() - 1000
	}

	return exchange.ConvertChartResolution(newRes*60, resolution, hp)
}

func (d *Deribit) getOHCLV(ticker string, res int64, st int64, et int64) ([]exchange.Candle, error) {
	var hp []exchange.Candle
	resp, err := d.get("public/get_tradingview_chart_data?end_timestamp="+strconv.FormatInt(et, 10)+
		"&instrument_name="+ticker+
		"&resolution="+strconv.FormatInt(res, 10)+
		"&start_timestamp="+strconv.FormatInt(st, 10), []byte(""))
	if err != nil {
		log.Printf("Error OHCLV Deribit %v", err)
		return hp, err
	}

	var tv GetTradingViewChartResponse
	err = processResponse(resp, &tv)
	return deribitCandleToCandle(tv.Result), err
}

// checkResolution looking if the asked resolution is a valid one
func checkResolution(res int64) int64 {
	return fnRes(res)
}

var fnRes = exchange.GenerateResolutionFunc(1440, 720, 360, 180,
	120, 60, 30, 15, 10, 5, 3, 1)

type GetTradingViewChartResponse struct {
	UsOut  int64                       `json:"usOut,omitempty"`
	UsIn   int64                       `json:"usIn,omitempty"`
	UsDiff int64                       `json:"usDiff,omitempty"`
	Id     int64                       `json:"id,omitempty"`
	Result getTradingViewChartResponse `json:"result,omitempty"`
}

type getTradingViewChartResponse struct {
	Status string    `json:"status,omitempty"`
	Volume []float64 `json:"volume,omitempty"`
	Cost   []float64 `json:"cost,omitempty"`
	Ticks  []int64   `json:"ticks,omitempty"`
	Open   []float64 `json:"open,omitempty"`
	Close  []float64 `json:"close,omitempty"`
	High   []float64 `json:"high,omitempty"`
	Low    []float64 `json:"low,omitempty"`
}

func deribitCandleToCandle(c getTradingViewChartResponse) []exchange.Candle {
	newChart := make([]exchange.Candle, 0, len(c.Ticks))
	var ec exchange.Candle
	for i := 0; i < len(c.Ticks); i++ {
		ec = exchange.Candle{
			Close:     c.Close[i],
			Open:      c.Open[i],
			High:      c.High[i],
			Low:       c.Low[i],
			StartTime: time.Unix(c.Ticks[i]/1000, 0),
			Volume:    c.Volume[i],
		}
		newChart = append(newChart, ec)
	}
	return newChart
}
