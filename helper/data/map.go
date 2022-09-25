package data

import "fmt"

type DataMap[T any] struct {
	mp    map[uint]T
	index uint
	limit uint
}

func Map[T any](f []T) *DataMap[T] {
	d := DataMap[T]{mp: make(map[uint]T, len(f)+20)}
	for i, v := range f {
		d.mp[uint(i)] = v
	}
	d.index = uint(len(f) - 1)
	return &d
}

func (d *DataMap[T]) V(index uint) T {
	v, ok := d.mp[d.index-uint(index)]
	if !ok {
		fmt.Println("fuck")
	}
	return v
}

func (d *DataMap[T]) SetLimit(limit uint) {
	d.limit = limit
}

func (d *DataMap[T]) Append(v T) {
	d.mp[d.index+1] = v
	delete(d.mp, d.index-d.limit)
	d.index++
}
