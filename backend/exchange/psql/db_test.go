package psql

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	SetPSQL("localhost", "postgres", "metapine", "admin", 5432)

	f, err := New("ftx")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	ch, err := f.OHCLV("btc-perp", 3600*4, time.Unix(0, 0), exchange.T2023)

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(len(ch))

	ch, err = ftx.New().OHCLV("btc-perp", 3600*4, time.Unix(0, 0), exchange.T2023)

	for i, v := range ch[1:] {
		if v.StartTime.Unix()-ch[i].StartTime.Unix() != 3600*4 {
			fmt.Println(i)
		}
	}

}
