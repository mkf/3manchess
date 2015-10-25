package game

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
