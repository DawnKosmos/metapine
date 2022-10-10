package psql

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	SetPSQL("localhost", "postgres", "metapine", "admin", 5432)
	err := p.Ping()
	fmt.Println(err)

	fmt.Println(getDbName("ftx", 1800))
	fmt.Println(getDbName("ftx", 7200))

}
