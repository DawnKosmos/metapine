package ta

type emptySeries[v any] struct {
	ERS[v]
}

func empty[T any](b []T, st, res int64) *emptySeries[T] {
	s := new(emptySeries[T])
	s.st = st
	s.res = res
	s.data = b
	return s
}
