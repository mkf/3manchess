package game

type Square [2]byte

func (s Square) Empty() bool {
	return s[0] == 0
}

func (s Square) Color() byte {
	return s[0]
}

func (s Square) What() byte {
	return s[1]
}

type Pos [2]uint8

type Board [6][24]Square

func (b *Board) Straight(from Pos, to Pos, m MoatsState) bool {
	var cantech, canmoat, canfig, canbridge bool
	canbridge = true
	if from[0]==to[0] {
		cantech = true
		if from[0]==0 {
			canbridge=false
		} else {
			canmoat = true
			canfig = true
			for i:=from[1]+1;((i-from[1])%24<(to[1]-from[1])%24)&&canfig;i=(i+1)%24 {
				go func() {
					if canfig&&!(*b[from[0]][i].Empty()) {
						canfig = false
					}
				}()
			}
			canfiga = true
			for i:=from[1]-1;((i-from[1])%24>(to[1]-from[1])%24)&&canfiga;i=(i-1)%24 {
				go func() {
					if canfiga&&!(*b[from[0]][i].Empty()) {
						canfiga = false
					}
				}()
			}
			canfig = canfig||canfiga
		}
	} else if from[1]==to[1] {
		cantech = true
		canmoat = true
		canfig = true
		//for i
	} else if ((from[1]-12)%24)==to[1] {
		cantech = true
		canmoat = true
	} else {
		return false
	}
}

//func (b Board) Diagonal

type MoatsState [3]bool //Black-White, White-Gray, Gray-Black

type MovesNext byte

type Castling [3][2]bool

func forcastlingconv(c byte, b byte) (bool,bool) {
	var col uint8
	switch c {
	case 'W', 'w':
		col = 0
	case 'G', 'g':
		col = 1
	case 'B', 'b':
		col = 2
	}
	var ct uint8
	switch b {
	case 'k', 'K':
		ct = 0
	case 'q', 'Q':
		ct = 1
	}
	return col, ct
}

func (cs Castling) Give(c byte, b byte) bool {
	col, ct := forcastlingconv(c, b)
	return cs[col][ct]
}

func (cs Castling) Change(c byte, b byte, w bool) Castling {
	cso = cs
	col, ct := forcastlingconv(c, b)
	cso[col][ct] = w
	return cso
}

type EnPassant []Pos

type HalfmoveClock uint8

type Fullmovenumber uint16

type State struct {
	*Board //[color,figure_lowercase] //[0,0]
	MoatsState
	MovesNext //W G B
	Castling  //0W 1G 2B  //0K 1Q
	EnPassant //[previousplayer,currentplayer]  [number,letter]
	HalfmoveClock
	FullmoveNumber
}

var COLORS = [3]byte{'W', 'G', 'B'}
var FIRSTRANKNEWGAME = [8]byte{'r', 'n', 'b', 'k', 'q', 'b', 'n', 'r'}

var DEFENPASSANT = make([][2]uint8, 0, 2)

var DEFCASTLING = [3][2]bool{
	{{true, true}, {true, true}},
	{{true, true}, {true, true}},
	{{true, true}, {true, true}},
}

var BOARDFORNEWGAME [6][24][2]byte

var NEWGAME State

func init() {
	for ci, c := range COLORS {
		for fi, f := range FIRSTRANKNEWGAME {
			a = ci*8 + fi
			BOARDFORNEWGAME[0][a][1] = f
			BOARDFORNEWGAME[0][a][0] = c
			BOARDFORNEWGAME[1][a][0] = c
			BOARDFORNEWGAME[1][a][1] = 'p'
			for l := 2; l < 6; l++ {
				BOARDFORNEWGAME[l][a][0] = 0
				BOARDFORNEWGAME[l][a][1] = 0
			}
		}
	}
	NEWGAME = State{&BOARDFORNEWGAME, 'W', DEFCASTLING, DEFENPASSANT, uint8(0), uint16(1)}
}

//func (s *State) String() string {   // returns FEN
//}

//func ParsBoard3FEN([]byte) *[8][24][2]byte {
//}

//func Pars3FEN([]byte) *State {
//}

type Move struct {
	From        [2]uint8
	To          [2]uint8
	What        [2]byte
	AlreadyHere [2]byte
	Before      *State
}

//func (m *Move) String() string {
//}

func (m *Move) After() *State {
	var movesnext byte
	if m.What[0] != m.Before.MovesNext {
		panic(m)
	}
	switch m.Before.MovesNext {
	case 'W':
		movesnext = 'G'
	case 'G':
		movesnext = 'B'
	case 'B':
		movesnext = 'W'
	}
	if m.What[0] == m.AlreadyHere[0] {
		panic(m)
	}

}
