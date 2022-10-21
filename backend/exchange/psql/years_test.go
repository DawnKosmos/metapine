package psql

import (
	"fmt"
	"testing"
)

func TestYears(t *testing.T) {
	fmt.Println(checkResolution(3600 * 4))
}
