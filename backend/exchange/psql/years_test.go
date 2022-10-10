package psql

import (
	"fmt"
	"testing"
)

func TestYears(t *testing.T) {
	var d DataArr

	fmt.Println(d.AddYear(12))
	fmt.Println(d.AddYear(5))
	d.AddYear(10)
	fmt.Println(d.AddYear(100))
	d.AddYear(20)
	fmt.Println(d.AddYear(121))
	d.AddYear(99)
	for _, v := range d.Years {
		fmt.Println(v.Y, v.Months)
	}
}
