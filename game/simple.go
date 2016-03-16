package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//ACFT — all combinations fromto
type ACFT FromTo

//P adds one to our FromTo
func (a *ACFT) P() {
	(*a)[0][0]++
	if (*a)[0][0] != 6 {
		return
	}
	(*a)[0][0] = 0
	(*a)[0][1]++
	if (*a)[0][1] != 24 {
		return
	}
	(*a)[0][1] = 0
	(*a)[1][0]++
	if (*a)[1][0] != 6 {
		return
	}
	(*a)[1][0] = 0
	(*a)[1][1]++
}

//OK checks whether it is correct
func (a *ACFT) OK() bool {
	return (*a)[1][1] != 24
}

//G checks whether it is correct and returns it
func (a *ACFT) G() (bool, FromTo) {
	return a.OK(), FromTo(*a)
}

//ACP — all combinations pos
type ACP Pos

//P add one to our Pos
func (a *ACP) P() {
	(*a)[0]++
	if (*a)[0] != 6 {
		return
	}
	(*a)[0] = 0
	(*a)[1]++
}

//OK checks whether it is correct
func (a *ACP) OK() bool {
	return (*a)[1] != 24
}

//G checks whether it is correct and returns it
func (a *ACP) G() (bool, Pos) {
	return a.OK(), Pos(*a)
}

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

func ynbool(b bool) byte {
	if b {
		return 'Y'
	}
	return 'N'
}
