package iterator

import "github.com/DawnKosmos/metapine/backend/series/ta"

type funcIterator struct {
	index int
	val   *func(src ta.Series, l int) ta.Series
	ss    []func(src ta.Series, l int) ta.Series
}

func (s *funcIterator) Next() bool {
	return s.index < len(s.ss)
}

func (s *funcIterator) Iterate() {
	s.index++
	if s.index < len(s.ss) {
		*s.val = s.ss[s.index]
	}

}

func (s *funcIterator) Reset() {
	s.index = 0
	*s.val = s.ss[0]
}
