package ta

import (
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"testing"
	"time"
)

func TestTypeTest(t *testing.T) {
	ff := ftx.New()

	ch := NewOHCLV(ff, "btc-perp", time.Unix(0, 0), exchange.T2023, 3600*24)
	//o, h, c, l, v := ChartSources(ch)
	c := Close(ch)
	rsi := Roc(c, 5)
	chD := ch.Data()[5:]

	for i, v := range rsi.Data() {
		fmt.Println(chD[i].StartTime.Format(time.Stamp), chD[i].Open, chD[i].Close, v)
	}

}
