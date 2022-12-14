package psql

type Month struct {
	M        int  `json:"m,omitempty"`
	Done     bool `json:"done,omitempty"`
	Complete bool `json:"complete,omitempty"`
}

type Year struct {
	Y      int     `json:"y,omitempty"`
	Months []Month `json:"months,omitempty"`
}

func NewYear(y int) Year {
	a := Year{Y: y}
	for i := 1; i < 13; i++ {
		a.Months = append(a.Months, Month{i, false, false})
	}
	return a
}

type DataArr struct {
	Years []Year `json:"years,omitempty"`
}

func (a *DataArr) AddMonth(y int, m int, done, complete bool) {
	pos := a.AddYear(y)
	a.Years[pos].Months[m] = Month{
		M:        m,
		Done:     done,
		Complete: complete,
	}
}

func (a *DataArr) AddYear(Y int) (arrPosition int) {
	for i, v := range a.Years {
		if v.Y == Y {
			return i
		}
		if Y < v.Y {
			old := append([]Year{NewYear(Y)}, a.Years[i:]...)
			a.Years = append(a.Years[:i], old...)
			return i
		}
	}
	a.Years = append(a.Years, NewYear(Y))
	return len(a.Years) - 1
}

func (a *DataArr) Has(y int, m int) bool {
	if len(a.Years) == 0 {
		return false
	}
	for _, v := range a.Years {
		if v.Y == y {
			if v.Months[m].Done {
				return true
			}
		}
	}
	return false
}
