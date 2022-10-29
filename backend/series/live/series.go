package live

import (
	"github.com/DawnKosmos/metapine/helper/data"
)

type Series interface {
	Val(index int) float64 //Val(0) returns the actual value, Val(1) the last etc.
	Data() []float64
	Updater
}

type Condition interface {
	Val(index int) bool //Val(0) returns the actual value, Val(1) the last etc.
	Data() []bool       //Needed for Initialisation
	Updater
}

type Updater interface {
	ResolutionStartTime
	OnTick(NewTick bool) //Updates the latest Tick, When Update(true) adds a Tick
	SetLimit(i int)      //Sets the Limit that needs to be allocated for the indicator to work
	ExecuteLimit()       //Gets called once
	GetUpdateGroup() *UpdateGroup
}

type ResolutionStartTime interface {
	StartTime() int64
	Resolution() int64
}

type URS[T any] struct {
	st, res int64
	data    data.Dater[T]
	ug      *UpdateGroup
	limit   int
	recent  float64
}

func (e *URS[T]) StartTime() int64 {
	return e.st
}

func (e *URS[T]) Resolution() int64 {
	return e.res
}

func (e *URS[T]) Data() []T {
	return e.data.Data()
}

func (e *URS[T]) Val(i int) T {
	return e.data.V(i)
}

func (e *URS[T]) SetLimit(limit int) {
	if limit > e.limit {
		limit = e.limit
	}
}

func (e *URS[T]) ExecuteLimit() {
	e.data.SetLimit(e.limit)
}

func (e *URS[T]) GetUpdateGroup() *UpdateGroup {
	return e.ug
}
