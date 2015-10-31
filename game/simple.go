package game

func sign(u int8) int8 {
	switch {
	case u == 0:
		return u
	case u < 0:
		return int8(-1)
	case u > 0:
		return int8(1)
	default:
		return int8(127)
	}
}

func abs(u int8) int8 {
	if u < 0 {
		return -u
	}
	return u
}

func min(i int8, j int8) int8 {
	if i < j {
		return i
	}
	return j
}

func max(i int8, j int8) int8 {
	if i > j {
		return i
	}
	return j
}
