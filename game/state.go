package game

func (b *Board) Straight(from Pos, to Pos, m MoatsState) (bool, bool) { //(whether it can, whether it can capture/check)
	var cantech, canmoat, canfig bool
	capcheck := true
	if from == to {
		panic("Same square!")
	}
	if from[0] == to[0] {
		cantech = true
		if from[0] == 0 {
			var mshort, mlong, capcheckshort bool
			if from[1]/8 == to[1]/8 {
				capcheckshort = true
				canmoat = true
				mshort = true
				if m[0] && m[1] && m[2] {
					mlong = true
				}
			} else {
				capcheckshort = false
				fromto := [2]int8{from[1] / 8, to[1] / 8}
				switch fromto {
				case [2]int8{0, 1}, [2]int8{1, 0}:
					mshort = m[1]
					mlong = m[0] && m[2]
				case [2]int8{1, 2}, [2]int8{2, 1}:
					mshort = m[2]
					mlong = m[0] && m[1]
				case [2]int8{2, 0}, [2]int8{0, 2}:
					mshort = m[0]
					mlong = m[1] && m[2]
				}
			}
			var direcshort int8
			//if to[0]
		} else {
			canmoat = true
			canfig = true
			for i := from[1] + 1; ((i-from[1])%24 < (to[1]-from[1])%24) && canfig; i = (i + 1) % 24 {
				go func() {
					if canfig && !((*b)[from[0]][i].Empty()) {
						canfig = false
					}
				}()
			}
			canfiga := true
			for i := from[1] - 1; ((i-from[1])%24 > (to[1]-from[1])%24) && canfiga; i = (i - 1) % 24 {
				go func() {
					if canfiga && !((*b)[from[0]][i].Empty()) {
						canfiga = false
					}
				}()
			}
			canfig = canfig || canfiga
		}
	} else if from[1] == to[1] {
		cantech = true
		canmoat = true
		canfig = true
		sgn := sign(to[0] - from[0])
		for i := from[0] + sgn; (sgn*i < to[0]) && canfig; i += sgn {
			go func() {
				if canfig && !((*b)[i][from[1]].Empty()) {
					canfig = false
				}
			}()
		}
	} else if ((from[1] - 12) % 24) == to[1] {
		cantech = true
		canmoat = true
		canfig = true
		for i, j := from[0], to[0]; canfig && (i < 6 && j < 6); i, j = i+1, j+1 {
			go func() {
				go func() {
					if canfig && !((*b)[i][from[1]].Empty()) {
						canfig = false
					}
				}()
				go func() {
					if canfig && !((*b)[j][to[1]].Empty()) {
						canfig = false
					}
				}()
			}()
		}
	}
	return false, false
}

//func (b Board) Diagonal

type MoatsState [3]bool //Black-White, White-Gray, Gray-Black

var DEFMOATSSTATE = MoatsState{true, true, true}

type Castling [3][2]bool

func forcastlingconv(c Color, b byte) (uint8, uint8) {
	col := c.UInt8()
	var ct uint8
	switch b {
	case 'k', 'K':
		ct = 0
	case 'q', 'Q':
		ct = 1
	}
	return col, ct
}

func (cs Castling) Give(c Color, b byte) bool {
	col, ct := forcastlingconv(c, b)
	return cs[col][ct]
}

func (cs Castling) Change(c Color, b byte, w bool) Castling {
	cso := cs
	col, ct := forcastlingconv(c, b)
	cso[col][ct] = w
	return cso
}

type EnPassant []Pos

type HalfmoveClock uint8

type FullmoveNumber uint16

type State struct {
	*Board //[color,figure_lowercase] //[0,0]
	MoatsState
	MovesNext Color //W G B
	Castling        //0W 1G 2B  //0K 1Q
	EnPassant       //[previousplayer,currentplayer]  [number,letter]
	HalfmoveClock
	FullmoveNumber
}

var DEFENPASSANT = make(EnPassant, 0, 2)

var DEFCASTLING = [3][2]bool{
	{true, true},
	{true, true},
	{true, true},
}

var NEWGAME State

func init() {
	boardinit()
	NEWGAME = State{&BOARDFORNEWGAME, DEFMOATSSTATE, Color('W'), DEFCASTLING, DEFENPASSANT, HalfmoveClock(0), FullmoveNumber(1)}
}

//func (s *State) String() string {   // returns FEN
//}

//func ParsBoard3FEN([]byte) *[8][24][2]byte {
//}

//func Pars3FEN([]byte) *State {
//}


