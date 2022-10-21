package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"time"
)

/*
Checkif table exists. Check empty spots
Download missing candles in months packets

*/

func ohclvMinute(indexId int32, ee exchange.CandleProvider, ticker string, st, et time.Time) ([]exchange.Candle, error) {
	mm, err := p.qq.ReadMinuteManager(ctx, indexId) //TODO fix that system in future
	if err != nil {
		if err == ErrNoRows {
			if err := createNewMinutesTable(indexId); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var dd DataArr
	if mm.Dataarr.Valid {
		json.Unmarshal([]byte(mm.Dataarr.String), &dd)
	}

	ym := getMonthsAndYears(st, et)
	for _, v := range ym[:1] {
		if !dd.Has(v.y, v.m) {
			tst, tet := getTimeStampsOfMonth(v.y, v.m)
			ch, err := ee.OHCLV(ticker, 60, tst, tet)
			if err != nil {
				fmt.Println(err)
			} else {
				if n, err := p.MinutesWriteOHCLV(ctx, minutesTable(indexId), ch); err != nil {
					fmt.Println(n)
					return nil, err
				}
				dd.AddMonth(v.y, v.m, true, true)
			}
		}
	}

	yy := ym[len(ym)-1]
	tNow := time.Now()
	if yy.y == tNow.Year() && yy.m == int(tNow.Month()) {

	} else {

	}

	p.qq.UpdateMinuteManager(ctx, gen.UpdateMinuteManagerParams{
		Dataarr: sql.NullString{},
		IndexID: indexId,
	})

	return nil, err
}
