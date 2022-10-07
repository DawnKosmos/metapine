package ftx

import (
	"fmt"
	"testing"
)

func TestKek(t *testing.T) {
	/*
		f := New()
		res, err := f.OHCLV("SOL-PERP", 3600*1, time.Unix(0, 0), time.Now())
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		var ii int64 = 3600 * 1
		var oc exchange.Candle
		fmt.Println(len(res))
		for i, v := range res {
			if i == 0 {
				oc = v
			} else {
				if ii != v.StartTime.Unix()-oc.StartTime.Unix() {
					fmt.Println(oc.StartTime, v.StartTime)
				}
				oc = v
			}
		}
	*/

	fmt.Println(checkResolution(3600 * 24 * 7))
	fmt.Println(checkResolution(3600 * 24))
	fmt.Println(checkResolution(3600 * 8))
	fmt.Println(checkResolution(3600 * 3))

}
