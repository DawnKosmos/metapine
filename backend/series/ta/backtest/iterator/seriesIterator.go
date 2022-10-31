package iterator

import "github.com/DawnKosmos/metapine/backend/series/ta"

type seriesIterator struct {
	index int
	val   *ta.Series
	ss    []ta.Series
}

func (s *seriesIterator) Next() bool {
	return s.index < len(s.ss)
}

func (s *seriesIterator) Iterate() {
	s.index++
	if s.index < len(s.ss) {
		*s.val = s.ss[s.index]
	}
}

func (s *seriesIterator) Reset() {
	s.index = 0
	*s.val = s.ss[0]
}
