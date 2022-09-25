package exchange

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type DB struct {
	q *pgx.Conn
}

var ctx = context.Background()

func (d *DB) Ping() error {
	return d.q.Ping(ctx)
}

/*
	BackTestModus:
	OHCLV(ticker string => "FTX,Btc-perp"



//Being able to register indexes
*/

func (d *DB) OHCLV(ticker string, resolution int64, start time.Time, end time.Time) ([]Candle, error) {
	//Separate exchange,ticker => theoretical also having index,ticker
	//CheckIf table exists

}
