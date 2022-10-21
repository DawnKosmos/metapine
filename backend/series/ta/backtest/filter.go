package backtest

type Filter func(sf []SafeFloat) bool

func Greater(index1 int, index2 int) Filter {
	return func(sf []SafeFloat) bool {
		if sf[index1].Safe && sf[index2].Safe {
			return sf[index1].Value > sf[index2].Value
		}
		return false
	}
}

func GreaterAs(index int, val float64) Filter {
	return func(sf []SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value > val
		}
		return false
	}
}

func SmallerAs(index int, val float64) Filter {
	return func(sf []SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value < val
		}
		return false
	}
}

func Between(index int, lowerBound, upperBound float64) Filter {
	return func(sf []SafeFloat) bool {
		if sf[index].Safe {
			return sf[index].Value > lowerBound && sf[index].Value < upperBound
		}
		return false
	}
}
