package game

type MoatsState [3]bool //Black-White, White-Gray, Gray-Black  //true: bridged. Originally, true meant still active, i.e. non-bridged!!!

//var DEFMOATSSTATE = MoatsState{true, true, true}
var DEFMOATSSTATE = MoatsState{false, false, false} //are they bridged?

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

func (cs Castling) OffRook(c Color, b byte) Castling {
	return cs.Change(c, b, false)
}

func (cs Castling) OffKing(c Color) Castling {
	return cs.OffRook(c, 'K').OffRook(c, 'Q')
}

type EnPassant [2]Pos

func (e EnPassant) Appeared(p Pos) EnPassant {
	ep := e
	ep[0] = ep[1]
	ep[1] = p
	return ep
}
func (e EnPassant) Nothing() EnPassant {
	return EnPassant{e[1], Pos{127, 127}}
}

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

func (s *State) AnyPiece(from, to Pos) bool {
	return s.Board.AnyPiece(from, to, s.MoatsState, s.Castling, s.EnPassant)
}

var DEFENPASSANT = EnPassant{Pos{127, 127}, Pos{127, 127}}

var DEFCASTLING = [3][2]bool{
	{true, true},
	{true, true},
	{true, true},
}

var FALSECASTLING = [3][2]bool{
	{false, false},
	{false, false},
	{false, false},
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
