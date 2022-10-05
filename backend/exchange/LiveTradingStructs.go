package exchange

import "time"

type SetOrderInput struct {
	Side       bool
	Ticker     string
	Price      float64
	Size       float64
	ReduceOnly bool
}

type Order struct {
	Id         int64     `json:"id,omitempty"`
	Market     string    `json:"market,omitempty"`
	Side       string    `json:"side,omitempty"`
	Size       float64   `json:"size,omitempty"`
	Price      float64   `json:"price,omitempty"`
	ReduceOnly bool      `json:"reduce_only,omitempty"`
	Created    time.Time `json:"created,omitempty"`
	FilledSize float64   `json:"filledSize,omitempty"`
}

type TriggerOrder struct {
	Id         int64     `json:"id,omitempty"`
	Ticker     string    `json:"ticker,omitempty"`
	Side       string    `json:"side,omitempty"`
	Size       float64   `json:"size,omitempty"`
	Price      float64   `json:"price,omitempty"`
	ReduceOnly bool      `json:"reduce_only,omitempty"`
	Created    time.Time `json:"created,omitempty"`
}

type Position struct {
	Side             string  `json:"side"`
	Future           string  `json:"future"`
	NotionalSize     float64 `json:"cost"`
	PositionSize     float64 `json:"size"`
	UPNL             float64 `json:"unrealizedPnl"`
	PNL              float64 `json:"realizedPnl"`
	EntryPrice       float64 `json:"entryPrice"`
	LiquidationPrice float64 `json:"estimatedLiquidationPrice"`
	AvgOpen          float64 `json:"recentAverageOpenPrice"`
	BreakEven        float64 `json:"recentBreakEvenPrice"`
}

type FundingPayments struct {
	Id              int64     `json:"id,omitempty"`
	Future          string    `json:"future,omitempty"`
	Payment         float64   `json:"payment,omitempty"`
	PaymentCurrency string    `json:"-"`
	Time            time.Time `json:"time,omitempty"`
}

type Fill struct {
	Fee          float64   `json:"fee,omitempty"`
	FeeCurrency  string    `json:"feeCurrency,omitempty"`
	Future       string    `json:"market,omitempty"`
	Id           int64     `json:"id,omitempty"`
	OrderId      int       `json:"orderId,omitempty"`
	Price        float64   `json:"price,omitempty"`
	BaseCurrency string    `json:"baseCurrency,omitempty"`
	Side         string    `json:"side,omitempty"`
	Size         float64   `json:"size,omitempty"`
	Time         time.Time `json:"time,omitempty"`
}
