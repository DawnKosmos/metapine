package size

type SizeBase int

const (
	Dollar SizeBase = iota
	Account
)

type Size struct {
	feeType SizeBase
	val     float64
}

func New(v SizeBase, val float64) *Size {
	return &Size{
		feeType: v,
		val:     val,
	}
}

func DefaultSize() *Size {
	return &Size{
		feeType: Account,
		val:     1,
	}
}

func (s *Size) Amount(price float64) float64 {
	switch s.feeType {
	case Dollar:
		return s.val / price
	case Account:
		return s.val
	}
	return 1
}
