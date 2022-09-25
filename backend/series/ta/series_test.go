package ta

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestTypeTest(t *testing.T) {
	rand.Seed(31121)
	var f []float64
	for i := 0; i < 100; i++ {
		f = append(f, float64(rand.Intn(100)))
	}

	src := empty(f, 1000, 3600)
	rsi := Rsi(src, 10)
	mm := Smaller(rsi, 50)
	f = rsi.Data()
	ff := mm.Data()

	for i, v := range f {
		fmt.Println(v, ff[i])
	}
}
