package iterator

type intIterator struct {
	val             *int
	start, end, add int
}

func (i *intIterator) Next() bool {
	return *i.val <= i.end
}

func (i *intIterator) Iterate() {
	*i.val = *i.val + i.add
}

func (i *intIterator) Reset() {
	*i.val = i.start
}
