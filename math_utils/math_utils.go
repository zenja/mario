package math_utils

func Min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
