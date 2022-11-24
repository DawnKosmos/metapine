package backtest

import (
	"fmt"
	"time"
)

type FillType int

const (
	LIMIT FillType = iota
	MARKET
	STOP
)

/* Fills describes a Fill done in a Trade */
type Fill struct {
	Side  bool
	Type  FillType
	Price float64
	Size  float64
	Time  time.Time
	Fee   float64
}

func (f Fill) String() string {
	side := "SELL"
	if f.Side {
		side = "BUY"
	}
	return fmt.Sprintf("%s /t Price:%f /t Size:%f %s", side, f.Price, f.Size, f.Time.Format(time.RFC822))
}
