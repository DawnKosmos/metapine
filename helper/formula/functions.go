package formula

func Change[T calc](old, new T) T {
	return new - old
}

func MedianOf3[T calc](a, b, c T) T {
	if a > b {
		if a < c {
			return a
		} else if b > c {
			return b
		} else {
			return c
		}
	} else {
		if a > c {
			return a
		} else if b < c {
			return b
		} else {
			return c
		}
	}
}

func MinInt(f ...int) (val int, position int) {
	val = f[0]
	for i, v := range f {
		if v < val {
			position = i
			val = v
		}
	}
	return
}

func MaxInt(f ...int) (val int, position int) {
	val = f[0]
	for i, v := range f {
		if v > val {
			position = i
			val = v
		}
	}
	return
}
