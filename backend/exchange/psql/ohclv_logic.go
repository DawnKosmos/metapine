package psql

import (
	"errors"
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

func ohclvTicker(indexId int64, exname string, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	ch, err = p.ohclv(ctx, ohclvQueueParams{
		Exchange:   exname,
		IndexId:    indexId,
		Resolution: resolution,
		StartTime:  start,
		EndTime:    end,
	})
	if err != nil {
		return nil, err
	}

	ee := stringToCandleProvider(exname)

	if start.Unix()-ch[0].StartTime.Unix() > resolution {
		tch, err := ee.OHCLV(ticker, resolution, start, ch[0].StartTime.Add(-1*time.Second))
		if err != nil {
			return nil, err
		}
		if len(tch) > 1 {
			ch = append(tch, ch...)
			if _, err = p.WriteOHCLV(ctx, fmt.Sprintf("%s:%s", exname, ticker), indexId, resolution, tch); err != nil {
				return nil, err
			}
		}
	}
	if end.Unix()-ch[len(ch)-1].StartTime.Unix() > resolution {
		tch, err := ee.OHCLV(ticker, resolution, ch[len(ch)-1].StartTime.Add(1*time.Second), end)
		if err != nil {
			return nil, err
		}
		if len(tch) > 1 {
			ch = append(ch, tch...)
			if _, err := p.WriteOHCLV(ctx, fmt.Sprintf("%s:%s", exname, ticker), indexId, resolution, tch); err != nil {
				return nil, err
			}
		}
	}

	return nil, err
}

// Gets Called when the Ticker does not exist and needs to be initialized
func initOhclv(name string, ticker string, resolution int64, start time.Time, end time.Time) (ch []exchange.Candle, err error) {
	var ee exchange.CandleProvider
	newRes := checkResolution(resolution)

	ee = stringToCandleProvider(name)

	if newRes >= 3600 {
		ch, err = ee.OHCLV(ticker, newRes, time.Unix(1236085200, 0), time.Now())
	} else {
		ch, err = ee.OHCLV(ticker, newRes, start, end)
	}
	if err != nil {
		return nil, err
	}
	if len(ch) == 0 {
		return nil, errors.New(fmt.Sprintf("Ticker exists but no data, %v, %v", start, end))
	}

	//Database Init
	e, err := stringToExchanges(name)
	if err != nil {
		return nil, err
	}

	tickerId, err := p.qq.CreateTicker(ctx, gen.CreateTickerParams{
		Exchange: e,
		Ticker:   ticker,
	})
	if err != nil {
		return nil, err
	}

	indexId, err := p.qq.CreateIndex(ctx, indexName(name, ticker))
	if err != nil {
		return nil, err
	}

	err = p.qq.CreateTickerIndex(ctx, gen.CreateTickerIndexParams{
		TickerID:      tickerId,
		IndexID:       indexId,
		Weight:        1,
		Excludevolume: false,
	})
	if err != nil {
		return nil, err
	}

	n, err := p.WriteOHCLV(ctx, name, int64(indexId), resolution, ch)
	if err != nil {
		return nil, err
	}
	p.loggin.Println(n, "New Candles got written in the DB")

	return ch, nil
}
