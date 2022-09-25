package data

type DataFloat struct {
	arr   []float64
	limit int
	allo  int
}

func New(f []float64) *DataFloat {
	d := new(DataFloat)
	d.arr = f
	return d
}

func (d *DataFloat) Append(v float64) {
	d.arr = append(d.arr, v)
	if (len(d.arr)) == d.allo {
		newArr := make([]float64, d.limit, d.allo)
		copy(newArr, d.arr[len(d.arr)-d.limit:])
		d.arr = newArr
	}
}

func (d *DataFloat) V(index uint) float64 {
	return d.arr[len(d.arr)-1-int(index)]
}

func (d *DataFloat) SetLimit(limit int) {
	limit++
	d.limit = limit
	if limit < 10 {
		d.allo = limit + 10
	} else {
		d.allo = limit * 2
	}
	if len(d.arr) < limit {
		limit = len(d.arr)
	}

	newArr := make([]float64, limit, d.allo)
	copy(newArr, d.arr[len(d.arr)-int(limit):])
	d.arr = newArr
}
