package iterator

import "github.com/DawnKosmos/metapine/backend/series/ta"

type ArrayIterator[T any] struct {
	index int
	val   *T
	ss    []T
}

func (s *ArrayIterator[T]) Next() bool {
	return s.index < len(s.ss)
}

func (s *ArrayIterator[T]) Iterate() {
	s.index++
	if s.index < len(s.ss) {
		*s.val = s.ss[s.index]
	}
}

func (s *ArrayIterator[T]) Reset() {
	s.index = 0
	*s.val = s.ss[0]
}

type IteratorGenerics[T any] interface {
	StructsAdresse() ([]*int, []*ta.Series, []*T)
	Calculation() (buy, sell ta.Condition)
	Parameters() string
}

type IterGeneric[T any] struct {
	it         IteratorGenerics[T]
	registered []iterator
	Parameter  []*int
	Srcs       []*ta.Series
	Generics   []*T
}

func NewGeneric[T any](it IteratorGenerics[T]) *IterGeneric[T] {
	iter := new(IterGeneric[T])
	iter.Parameter, iter.Srcs, iter.Generics = it.StructsAdresse()
	iter.it = it
	return iter
}

func (it *IterGeneric[T]) RegisterInt(position, start, end, add int) {
	if position >= len(it.Parameter) || add == 0 {
		return
	}
	if start > end && add > 0 {
		return
	}
	if start < end && add < 0 {
		return
	}

	*it.Parameter[position] = start
	ii := &intIterator{
		val:   it.Parameter[position],
		start: start,
		end:   end,
		add:   add,
	}
	it.registered = append(it.registered, ii)
}

func (it *IterGeneric[T]) RegisterSeries(position int, src ...ta.Series) {
	if position >= len(it.Parameter) || len(src) == 0 {
		return
	}

	*it.Srcs[position] = src[0]
	ii := &seriesIterator{
		index: position,
		val:   it.Srcs[position],
		ss:    src,
	}
	it.registered = append(it.registered, ii)
}

func (it *IterGeneric[T]) RegisterGenerics(position int, src ...T) {
	*it.Generics[position] = src[0]
	ii := &ArrayIterator[T]{
		index: position,
		val:   it.Generics[position],
		ss:    src,
	}
	it.registered = append(it.registered, ii)
}
