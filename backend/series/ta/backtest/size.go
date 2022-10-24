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

func DefaultSize() *Size {
	return &Size{
		Type: AccountSize,
		Val:  100,
	}
}
