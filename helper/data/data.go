package data

type Dater[T any] interface {
	V(index int) T
	SetLimit(limit int)
	Append(v T)
	SetValue(index int, val T)
	Data() []T
}
