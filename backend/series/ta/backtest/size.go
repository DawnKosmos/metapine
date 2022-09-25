package backtest

type SizeBase int

const (
	Dollar SizeBase = iota
	CurrencyBase
	AccountSize
)

type Size struct {
	Type SizeBase
	Val  float64
}
