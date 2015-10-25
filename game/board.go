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

type FigType byte

var Pawn = FigType('p')
var Rook = FigType('r')
var Knight = FigType('n')
var Bishop = FigType('b')
var Queen = FigType('q')
var King = FigType('k')

type Fig struct {
	FigType
	Color
}

type Square struct {
	Fig
	NotEmpty bool
}

func (s Square) Empty() bool {
	return !s.NotEmpty
}

func (s Square) Color() Color {
	return s.Fig.Color
}

func (s Square) What() FigType {
	return s.Fig.FigType
}

type Pos [2]int8

type Board [6][24]Square

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

type Color byte

func (c Color) UInt8() uint8 {
	switch c {
	case 'W', 'w':
		return 0
	case 'G', 'g':
		return 1
	case 'B', 'b':
		return 2
	}
	panic(c)
}

var White = Color('W')
var Gray = Color('G')
var Black = Color('B')

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

var COLORS = [3]Color{White, Gray, Black}
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, King, Queen, Bishop, Knight, Rook}

var DEFENPASSANT = make(EnPassant, 0, 2)

var DEFCASTLING = [3][2]bool{
	{true, true},
	{true, true},
	{true, true},
}

var BOARDFORNEWGAME Board

var NEWGAME State

func init() {
	for ci, c := range COLORS {
		for fi, f := range FIRSTRANKNEWGAME {
			a := ci*8 + fi
			BOARDFORNEWGAME[0][a].FigType = f
			BOARDFORNEWGAME[0][a].Fig.Color = c
			BOARDFORNEWGAME[0][a].NotEmpty = true
			BOARDFORNEWGAME[1][a].Fig.Color = c
			BOARDFORNEWGAME[1][a].Fig.FigType = Pawn
			BOARDFORNEWGAME[1][a].NotEmpty = true
			for l := 2; l < 6; l++ {
				BOARDFORNEWGAME[l][a].Fig.Color = 0
				BOARDFORNEWGAME[l][a].Fig.FigType = 0
				BOARDFORNEWGAME[l][a].NotEmpty = false
			}
		}
	}
	NEWGAME = State{&BOARDFORNEWGAME, DEFMOATSSTATE, Color('W'), DEFCASTLING, DEFENPASSANT, HalfmoveClock(0), FullmoveNumber(1)}
}

//func (s *State) String() string {   // returns FEN
//}

//func ParsBoard3FEN([]byte) *[8][24][2]byte {
//}

//func Pars3FEN([]byte) *State {
//}

type Move struct {
	From        Pos
	To          Pos
	What        Fig
	AlreadyHere Fig
	Before      *State
}

//func (m *Move) String() string {
//}

func (m *Move) After() *State {
	var movesnext Color
	if m.What.Color != m.Before.MovesNext {
		panic(m)
	}
	switch m.Before.MovesNext {
	case White:
		movesnext = Gray
	case Gray:
		movesnext = Black
	case Black:
		movesnext = White
	}
	if m.What.Color == m.AlreadyHere.Color {
		panic(m)
	}

}
