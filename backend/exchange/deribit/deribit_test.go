package deribit

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"net/http"
	"testing"
)

func TestOHCLV(t *testing.T) {
	d := &Deribit{client: http.DefaultClient}

	ch, err := d.OHCLV("ETH-PERPETUAL", 3600*6, exchange.T2020, exchange.T2023)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(len(ch))
	st := ch[0].StartTime.Unix()
	for i, v := range ch[1:] {
		if v.StartTime.Unix()-st != 3600*6 {
			fmt.Println(i)
		}
		st = v.StartTime.Unix()
	}
}
