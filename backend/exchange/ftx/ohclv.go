package ftx

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DawnKosmos/metapine/backend/exchange"
)

func (f *FTX) OHCLV(ticker string, resolution int64, startTime time.Time, endTime time.Time) ([]exchange.Candle, error) {
	var hp []exchange.Candle
	st, et := startTime.Unix(), endTime.Unix()
	var end int64 = et

	if time.Now().Unix() < end {
		end = time.Now().Unix()
	}
	newRes := checkResolution(resolution)
	for {
		c, err := f.getOHCLV(ticker, newRes, st, end)
		if err != nil {
			fmt.Println(st, end)
			log.Printf("Error OHCLV FTX %v", err)
			return hp, err
		}
		fmt.Println(len(c))
		if len(c) == 0 {
			break
		}

		hp = append(c, hp...)
		end = hp[0].StartTime.Unix() - 1
	}
	return exchange.ConvertChartResolution(newRes, resolution, hp)
}

type HistoricalPrices struct {
	Success bool              `json:"success"`
	Result  []exchange.Candle `json:"result"`
}

func (f *FTX) getOHCLV(ticker string, res int64, st int64, et int64) ([]exchange.Candle, error) {
	var h HistoricalPrices
	resp, err := f.get(
		"markets/"+ticker+
			"/candles?resolution="+strconv.FormatInt(res, 10)+
			"&start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		log.Printf("Error OHCLV FTX %v", err)
		return h.Result, err
	}
	err = processResponse(resp, &h)
	return h.Result, nil

}

// checkResolution looking if the asked resolution is a valid one
func checkResolution(res int64) int64 {
	return fnRes(res)
}

var fnRes = exchange.GenerateResolutionFunc(86400*7, 86400*2, 86400, 14400,
	3600, 900, 300, 60, 15)
