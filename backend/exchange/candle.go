package exchange

import (
	"fmt"
	"time"
)

type Candle struct {
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Close     float64   `json:"close"`
	Low       float64   `json:"low"`
	Volume    float64   `json:"volume"`
	StartTime time.Time `json:"startTime"`
}

//New res must me greater than old
func ConvertChartResolution(oldResolution, newResolution int64, ch []Candle) ([]Candle, error) {
	if newResolution == oldResolution {
		return ch, nil
	}

	if oldResolution > newResolution || newResolution%oldResolution != 0 {
		return ch, fmt.Errorf("New Res %v and old %v do not fit", newResolution, oldResolution)
	}

	quotient := int(newResolution / oldResolution)
	var newChart []Candle = make([]Candle, 0, len(ch)/quotient)

	for _, c := range ch {
		if c.StartTime.Unix()%newResolution != 0 {
			ch = ch[1:]
		} else {
			break
		}
	}
	for {
		if len(ch) < quotient {
			break
		}
		newChart = append(newChart, ConvertCandleResolution(ch[:quotient]))
		ch = ch[quotient:]
	}
	if len(ch) != 0 {
		newChart = append(newChart, ConvertCandleResolution(ch))
	}
	return newChart, nil
}

//ConvertResolution converts the a lower resolution into a higher resolution
func ConvertCandleResolution(c []Candle) Candle {
	var out Candle = Candle{c[0].Close, c[0].High, c[0].Low, c[0].Open, c[0].Volume, c[0].StartTime}
	if len(c) == 1 {
		return c[0]
	}
	for _, i := range c[1:] {
		out.Close = i.Close
		out.Volume += i.Volume
		if i.High > out.High {
			out.High = i.High
		}
		if i.Low < out.Low {
			out.Low = i.Low
		}
	}
	return out
}
