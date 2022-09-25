package data

import (
	"fmt"
	"testing"
	"time"
)

func TestData(t *testing.T) {

	v := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	d := Map(v)
	d.SetLimit(10)
	var sum float64
	fmt.Println(d.V(0), d.V(9))
	tNow := time.Now()
	for i := 10; i < 1000000000; i++ {
		d.Append(float64(i))
		sum += d.V(10) + d.V(0)
	}
	tPast := time.Now()
	fmt.Println(len(d.mp), sum, tPast.UnixMicro()-tNow.UnixMicro())
}
