package data

type array[T any] struct {
	arr           []T
	limit         int
	allocateLimit int
}

func Array[T any](f []T) *array[T] {
	d := new(array[T])
	d.arr = f
	return d
}

func (d *array[T]) Append(v T) {
	d.arr = append(d.arr, v)
	if (len(d.arr)) == d.allocateLimit {
		newArr := make([]T, d.limit, d.allocateLimit)
		copy(newArr, d.arr[len(d.arr)-d.limit:])
		d.arr = newArr
	}
}

func (d *array[T]) V(index int) T {
	return d.arr[len(d.arr)-1-index]
}

func (d *array[T]) SetLimit(limit int) {
	limit++
	d.limit = limit
	if limit < 10 {
		d.allocateLimit = limit + 10
	} else {
		d.allocateLimit = limit * 2
	}
	if len(d.arr) < limit {
		limit = len(d.arr)
	}

	newArr := make([]T, limit, d.allocateLimit)
	copy(newArr, d.arr[len(d.arr)-limit:])
	d.arr = newArr
}

func (d *array[T]) SetValue(index int, val T) {
	d.arr[len(d.arr)-1-index] = val
}

func (d *array[T]) Data() []T {
	return d.arr
}
