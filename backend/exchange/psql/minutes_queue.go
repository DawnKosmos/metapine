package psql

import (
	"database/sql"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"github.com/jackc/pgtype"
	"time"
)

/*
To manage 1 min data efficient. if we import them always a whole month gets downloaded
A own table manages the name and the months used by this database
a month starts with time.Data(YEAR, time.Month, 0,0,0,0,0,time.UTC)
and ends with time.Data(YEAR, time.Month+1, 0,0,0,-1,0,time.UTC)
the current month get treaten differently

*/

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

func (a *DataArr) AddYear(Y int) (arrPosition int) {
	for i, v := range a.Years {
		if v.Y == Y {
			return i
		}
		if Y < v.Y {
			a.Years = append(append(a.Years[:i], NewYear(Y)), a.Years[i:]...)
			return i
		}
	}
	a.Years = append(a.Years, NewYear(Y))
	return len(a.Years) - 1
}

func kek() {

	var d DataArr
	b := pgtype.JSON{}
	b.Set(d)

	k := p.qq.CreateMinuteManager(ctx, gen.CreateMinuteManagerParams{
		IndexID:   sql.NullInt32{},
		Tablename: "",
		Dataarr:   b,
	})
	time.Date(2021, time.January, 0, 0, 0, 0, 0, time.UTC)
}
