package psql

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"time"
)

/*
Checkif table exists. Check empty spots
Download missing candles in months packets

*/

func ohclvMinute(indexId int32, ee exchange.CandleProvider, ticker string, st, et time.Time) ([]exchange.Candle, error) {
	mm, err := p.qq.ReadMinuteManager(ctx, indexId)
	if err != nil {
		return nil, err
	}
	fmt.Println(mm.Tablename)
	return nil, err
}
