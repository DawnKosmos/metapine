package ta

import "github.com/DawnKosmos/metapine/backend/exchange"

type source struct {
	ERS[float64]
}

func Source(o Chart, op func(candle exchange.Candle) float64, name string) Series {
	s := new(source)
	s.st = o.StartTime()
	s.res = o.Resolution()
	s.name = name
	d := make([]float64, 0, len(o.Data()))
	for _, c := range o.Data() {
		d = append(d, op(c))
	}
	s.data = d
	return s
}

func Open(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return e.Open
	}
	return Source(c, fn, "Open")
}

func Close(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return e.Close
	}
	return Source(c, fn, "Close")
}

func High(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return e.High
	}
	return Source(c, fn, "High")
}

func Low(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return e.Low
	}
	return Source(c, fn, "Low")
}

func Volume(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return e.Volume
	}
	return Source(c, fn, "Volume")
}

func HL2(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return (e.High + e.Low) / 2
	}
	return Source(c, fn, "HL2")
}

func OHCL4(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return (e.Open + e.Close + e.High + e.Low) / 4
	}
	return Source(c, fn, "OHCL4")
}

func OC2(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return (e.Open + e.Close) / 2
	}
	return Source(c, fn, "OC2")
}

func HCL3(c Chart) Series {
	fn := func(e exchange.Candle) float64 {
		return (e.Close + e.High + e.Low) / 3
	}
	return Source(c, fn, "HCL3")
}
