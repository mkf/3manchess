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

//VFTPGen : generates all valid FromToProm, given the game state
func VFTPGen(gamestate *State) <-chan FromToProm {
	all_valid := make(chan FromToProm)
	go func() {
		var y1, x1, y2, x2 int8
		for y1 = 0; y1 < 6; y1++ {
			for x1 = 0; x1 < 24; x1++ {
				for y2 = 0; y2 < 6; y2++ {
					for x2 = 0; x2 < 24; x2++ {
						p1 := Pos{y1, x1}
						p2 := Pos{y2, x2}
						fromto := FromTo{p1, p2}
						move := Move{p1, p2, gamestate, Queen}
						if _, err := move.After(); err == nil {
							all_valid <- FromToProm{fromto, Queen}
							fig := (*gamestate).Board[y1][x1].Fig
							if fig.FigType == Pawn && fig.PawnCenter && y1 == 1 {
								all_valid <- FromToProm{fromto, Rook}
								all_valid <- FromToProm{fromto, Bishop}
								all_valid <- FromToProm{fromto, Knight}
							}
						}
					}
				}
			}
		}
		close(all_valid)
	}()
	return all_valid
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
