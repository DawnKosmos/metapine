package live

/*
Just wonder how to create this:

ftx := db.LiveTicker("index", "BTC", 3600*HOUR, 200) //200 mean 100 candles to start with

ftx := db.LiveTickers("index", {"btc","eth","sol","ftt", 3600*HOUR, 200) //200 mean 200 init candles to start with
ftx.RegisterAccount("key", "name")

	r1 := Sma(Rsi(src1, l1), 2)
	r2 := Sma(Rsi(src2, l2), 2)
	r := AddF(r1, r2, 2)
	b1 := Sma(r, l2)
	b2 := Sma(b1, l2)
	b := SubF(b1, b2, 2)
	cc := Sub(r, b)
	saphir = Sma(cc, 2)

	buy := Crossover(saphir, cc)
	sell := Crossunder(saphir,cc)

	strat := strategy.New(strategy.Parameter{pyramiding: 2, Size: [size.Account,50]}

	strat.Long(ftx.Account("name"),buy, strategy.ONCLOSE, execute.Market)
	start.Short(ftx.Account("name"), sell, strategy.ONCLOSE, execute.Marke)

	ftx.Run
*/
