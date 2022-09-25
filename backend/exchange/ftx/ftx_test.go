package ftx

import (
	"fmt"
	"testing"
	"time"
)

func TestKek(t *testing.T) {
	f := New()
	res, err := f.OHCLV("BTC-PERP", 3600*24, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, v := range res {
		fmt.Println(v)
	}

}
