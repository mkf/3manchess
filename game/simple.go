package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

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

func absu(i int8) uint8 {
	if i < 0 {
		return uint8(-i)
	}
	return uint8(i)
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

func ynbool(b bool) byte {
	if b {
		return 'Y'
	}
	return 'N'
}
