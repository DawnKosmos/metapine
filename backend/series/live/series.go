package live

import (
	"github.com/DawnKosmos/metapine/backend/series/ta"
)

type Series interface {
	Val(i int) float64 //Val(0) returns the actual value, Val(1) the last etc.
	//ValArr(first, last int64) []float64 //ValArr returns an array of numbers ValArr(0,4) returns the latest 5 candles
	Updater
	Data() []float64
	ResolutionStartTime
}

type Condition interface {
	Val(i int) bool //Val(0) returns the actual value, Val(1) the last etc.
	//	ValArr(first, last int64) []bool //ValArr returns an array of numbers ValArr(0,4) returns the latest 5 candles
	Updater
	Data() []bool
	ResolutionStartTime
}

type ResolutionStartTime interface {
	StartTime() int64
	Resolution() int64
}

type ERS struct {
	ta ta.ResolutionStartTime //
}

func (e *ERS) StartTime() int64 {
	return e.ta.StartTime()
}

func (e *ERS) Resolution() int64 {
	return e.ta.Resolution()
}
