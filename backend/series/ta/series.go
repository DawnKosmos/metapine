package ta

//A Series  is an Interface that Every Indicator needs to implement to Communicate its Value to other Indicators. OHCLV Data is also an Indicator
type Series interface {
	Data() []float64
	ResolutionStartTime
}

//A Condition is a Interface that Every Condition Series needs to implement. Such as And, Greater, OR, IFF
type Condition interface {
	Data() []bool
	ResolutionStartTime
}

//ErrorResolutionStartime is needed to sync the Indicators in a fast way
type ResolutionStartTime interface {
	StartTime() int64
	Resolution() int64
}

//ERS implements the ErrorResolutionStartTime and can be implemented in a Series and Condition
type ERS[T any] struct {
	st   int64
	res  int64
	data []T
}

func (e *ERS[T]) StartTime() int64 {
	return e.st
}

func (e *ERS[T]) Resolution() int64 {
	return e.res
}

func (e *ERS[T]) Data() []T {
	return e.data
}
