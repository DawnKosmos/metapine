package psql

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	SetPSQL("localhost", "postgres", "metapine", "admin", 5432)
	err := p.Ping()
	fmt.Println(err)

	_, err = p.qq.ReadMinuteManager(ctx, 10)

	fmt.Println(ErrNoRows == err, err)
}
