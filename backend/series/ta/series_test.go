package ta

import (
	"github.com/DawnKosmos/metapine/backend/exchange"
	"github.com/DawnKosmos/metapine/backend/exchange/ftx"
	"testing"
)

func TestTypeTest(t *testing.T) {
	ff := ftx.New()

	ch := NewOHCLV(ff, "btc-perp", exchange.T2022, exchange.T2023, 3600*4)
	//o, h, c, l, v := ChartSources(ch)
	c := Close(ch)
	rsi := Ema(c, 10)
	Crossover(rsi)

}
