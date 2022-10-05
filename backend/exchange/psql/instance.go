package psql

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"log"
	"time"
)

type Instance struct {
	Exchange   string
	Ticker     string
	Resolution int64
	ch         []exchange.Candle
}

func New(exchange string) *Instance {
	if p == nil {
		log.Panicln("set a database...")
	}
	return &Instance{Exchange: exchange}
}

func (in *Instance) Val(index int) exchange.Candle {
	//TODO implement me
	panic("implement me")
}

func (in *Instance) OHCLV(exchange string, ticker string, resolution int64, start time.Time, end time.Time) ([]exchange.Candle, error) {
	//check if index, if no add the exchange
	//Check if DatabaseEntry exists
	//Download if needed
}

func (in *Instance) Name() string {
	//TODO implement me
	panic("implement me")
}
