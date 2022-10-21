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

func Min[T calc](f ...T) (T, int) {
	val := f[0]
	var pos int
	for i, v := range f[1:] {
		if v > val {
			pos = i
			val = v
		}
	}
	return val, pos
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

func Last[T any](a []T) T {
	return a[len(a)-1]
}
