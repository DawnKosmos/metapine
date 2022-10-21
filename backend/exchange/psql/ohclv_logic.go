package psql

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"time"
)

/*
Funkionweise:
	1. Index existiert nicht:
		downloade alle möglichen daten runter OHCLV(0,now)
	2. Letzten 50 bars fehlern
		GetOHCLV -> Download rest
	3. Ticker hat lücken
		Auf Index zugreifen

	Problem:
		aktuellste Kerze muss doublechecked werden:
			lsg: jedes mal die top kerzen checken
			lsg2: höchste kerze nicht eintragen.

1 min probleme:
	lsg. Monatsweise Datenrunterladen, abspeichern in eigener Tabelle 4jahre 2 millionen einträge sind


*/

func ohclvTicker(indexId int32, ee exchange.CandleProvider, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	if resolution == 60 {
		return ohclvMinute(indexId, ee, ticker, start, end)
	}
	return lowOHCLV(indexId, ee, ticker, resolution, start, end)
}

func lowOHCLV(index int32, ee exchange.CandleProvider, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	newRes := checkResolution(resolution)

	args, err := p.qq.ReadTickerManager(ctx, gen.ReadTickerManagerParams{
		IndexID:    index,
		Resolution: int32(newRes),
	})

	if err != nil {
		return initOhclv(index, ee, ticker, resolution, start, end)
	}
	//The most actual candle isn't closed, therefore we don't write it into the database
	var lc exchange.Candle //LastCandle
	//Add candles to the database if missing

	if args.St.Unix() > start.Unix() {
		tch, err := ee.OHCLV(ticker, newRes, start, args.St.Add(-1*time.Second))
		if err != nil {
			fmt.Println(err)
		} else {
			if len(tch) == 0 {
				args.St = time.Unix(0, 0).UTC()
			} else {
				n, err := p.WriteOHCLV(ctx, index, ee.Name(), newRes, tch)
				if err != nil {
					return nil, err
				}
				args.St = tch[0].StartTime
				p.loggin.Println(ticker, newRes, "|", n, "| lines got added")
			}
		}
	}
	var n int64
	//Add candles that are missing
	if args.Et.Unix() < end.Unix() {
		tch, err := ee.OHCLV(ticker, newRes, args.Et.Add(1*time.Second), end)
		if err != nil || len(tch) == 0 {
			fmt.Println(err)
		} else {
			fmt.Println("first", len(tch))
			lc = tch[len(tch)-1]
			if int64(time.Now().Sub(lc.StartTime)/time.Second) < resolution {
				n, err = p.WriteOHCLV(ctx, index, ee.Name(), newRes, tch[:len(tch)-1])
				if len(tch) > 1 {
					args.Et = tch[len(tch)-2].StartTime
				}
			} else {
				fmt.Println("second", len(tch))

				n, err = p.WriteOHCLV(ctx, index, ee.Name(), newRes, tch)
				args.Et = lc.StartTime
			}
			p.loggin.Println(ticker, newRes, "|", n, "| lines got added")
		}
	}

	err = p.qq.UpdateTickerManager(ctx, gen.UpdateTickerManagerParams{
		St:         args.St,
		Et:         args.Et,
		IndexID:    index,
		Resolution: int32(newRes),
	})
	if err != nil {
		fmt.Println("updatey not workugn", err)
	}

	out, err := p.ohclv(ctx, ohclvQueueParams{
		Exchange:   ee.Name(),
		IndexId:    index,
		Resolution: newRes,
		StartTime:  start,
		EndTime:    end,
	})
	if err != nil {
		return nil, err
	}

	if time.Now().Unix()-last(out).StartTime.Unix() < newRes {
		out = append(out, lc)
	}

	return exchange.ConvertChartResolution(newRes, resolution, out)
}

// INIT functions
// Gets Called when the Ticker+Resolution does not exist and needs to be initialized
func initOhclv(indexId int32, ee exchange.CandleProvider, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	newRes := checkResolution(resolution)

	switch {
	case newRes >= 3600*3:
		ch, err = highInitOhclv(indexId, ee, ticker, newRes, start, end)
		if err != nil {
			return nil, err
		}
	case 60 == newRes:
	default:
		ch, err = lowInitOhclv(indexId, ee, ticker, newRes, start, end)
		if err != nil {
			return nil, err
		}
	}

	return exchange.ConvertChartResolution(newRes, resolution, ch)
}

func lowInitOhclv(indexId int32, ee exchange.CandleProvider, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	ch, err = ee.OHCLV(ticker, resolution, start, end)
	if err != nil {
		return nil, err
	}
	lc := ch[len(ch)-1]

	var n int64
	if int64(time.Now().Sub(lc.StartTime)/time.Second) < resolution {
		n, err = p.WriteOHCLV(ctx, indexId, ee.Name(), resolution, ch[:len(ch)-1])
		lc = ch[len(ch)-2]
	} else {
		n, err = p.WriteOHCLV(ctx, indexId, ee.Name(), resolution, ch)
	}
	if err != nil {
		return nil, err
	}
	p.loggin.Println(n, "candles got inserted")

	p.qq.CreateTickerManager(ctx, gen.CreateTickerManagerParams{
		IndexID:    indexId,
		Resolution: int32(resolution),
		St:         ch[0].StartTime,
		Et:         lc.StartTime,
	})
	return ch, err
}

func highInitOhclv(indexId int32, ee exchange.CandleProvider, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	ch, err = ee.OHCLV(ticker, resolution, time.Unix(0, 0), time.Now())
	if err != nil {
		return nil, err
	}
	lc := ch[len(ch)-1]
	n, err := p.WriteOHCLV(ctx, indexId, ee.Name(), resolution, ch[:len(ch)-1])
	if err != nil {
		return nil, err
	}
	p.loggin.Println(n, "candles got inserted")

	p.qq.CreateTickerManager(ctx, gen.CreateTickerManagerParams{
		IndexID:    indexId,
		Resolution: int32(resolution),
		St:         time.Unix(0, 0), // all candles
		Et:         ch[len(ch)-2].StartTime,
	})

	tch, err := p.ohclv(ctx, ohclvQueueParams{
		Exchange:   ee.Name(),
		IndexId:    indexId,
		Resolution: resolution,
		StartTime:  start,
		EndTime:    end,
	})
	if err != nil {
		return nil, err
	}

	if lc.StartTime.Unix()-last(tch).StartTime.Unix() < resolution {
		tch = append(tch, lc)
	}
	return tch, err
}

func FindMissing(st, et, start, end time.Time) {

}
