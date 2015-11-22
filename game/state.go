package game

import "fmt"

//MoatsState :  Black-White, White-Gray, Gray-Black  //true: bridged. Originally, true meant still active, i.e. non-bridged!!!
type MoatsState [3]bool

//var DEFMOATSSTATE = MoatsState{true, true, true}

//DEFMOATSSTATE , ie. nothig is bridged     const
var DEFMOATSSTATE = MoatsState{false, false, false} //are they bridged?

//Castling : White,Gray,Black King-side,Queen-side
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

//Give (color, K/Q)
func (cs Castling) Give(c Color, b byte) bool {
	col, ct := forcastlingconv(c, b)
	return cs[col][ct]
}

//Change (color, K/Q, bool)
func (cs Castling) Change(c Color, b byte, w bool) Castling {
	cso := cs
	col, ct := forcastlingconv(c, b)
	cso[col][ct] = w
	return cso
}

//OffRook : Rook can no longr castle
func (cs Castling) OffRook(c Color, b byte) Castling {
	return cs.Change(c, b, false)
}

//OffKing : No castling anymore for this player
func (cs Castling) OffKing(c Color) Castling {
	return cs.OffRook(c, 'K').OffRook(c, 'Q')
}

//EnPassant type : two positions of enpassant, moving one left on each move
type EnPassant [2]Pos

//Appeared : new EnPassant possibility
func (e EnPassant) Appeared(p Pos) EnPassant {
	ep := e
	ep[0] = ep[1]
	ep[1] = p
	return ep
}

//Nothing : just a move, no new enpassant possibility
func (e EnPassant) Nothing() EnPassant {
	return EnPassant{e[1], Pos{127, 127}}
}

//HalfmoveClock : not used atm, TODO
type HalfmoveClock uint8

//FullmoveNumber : not used atm, TODO
type FullmoveNumber uint16

//PlayersAlive : which players are still active
type PlayersAlive [3]bool

//Give : tell if a player is active by color
func (pa PlayersAlive) Give(who Color) bool {
	return pa[who.UInt8()]
}

//Die : disactivate a player
func (pa PlayersAlive) Die(who Color) PlayersAlive {
	pan := pa
	pan[who.UInt8()] = false
	return pan
}

//DEFPLAYERSALIVE : true,true,true const
var DEFPLAYERSALIVE = [3]bool{true, true, true}

//State : single gamestate
type State struct {
	*Board //[color,figure_lowercase] //[0,0]
	MoatsState
	MovesNext Color //W G B
	Castling        //0W 1G 2B  //0K 1Q
	EnPassant       //[previousplayer,currentplayer]  [number,letter]
	HalfmoveClock
	FullmoveNumber
	PlayersAlive
}

func (s *State) EvalDeath() *State {
	if !(s.CanIMoveWOCheck(s.MovesNext)) {
		s.PlayersAlive.Die(next.MovesNext)
	}
	return s
}

func (s *State) String() string {
	return fmt.Sprintln("Board: ", (*s.Board), s.MoatsState, s.MovesNext, s.Castling, s.EnPassant, s.HalfmoveClock, s.FullmoveNumber, s.PlayersAlive)
}

//AnyPiece : if a piece could move (any piece, whatever stays there)
func (s *State) AnyPiece(from, to Pos) bool {
	return s.Board.AnyPiece(from, to, s.MoatsState, s.Castling, s.EnPassant)
}

//DEFENPASSANT : empty enpassant   const
var DEFENPASSANT = EnPassant{Pos{127, 127}, Pos{127, 127}}

//DEFCASTLING : everybody capable of castling everywhere  const
var DEFCASTLING = [3][2]bool{
	{true, true},
	{true, true},
	{true, true},
}

//FALSECASTLING : nobody can castle anymore   const
var FALSECASTLING = [3][2]bool{
	{false, false},
	{false, false},
	{false, false},
}

//NEWGAME : gamestate of a new game   const
var NEWGAME State

func init() {
	boardinit()
	NEWGAME = State{&BOARDFORNEWGAME, DEFMOATSSTATE, Color('W'), DEFCASTLING, DEFENPASSANT, HalfmoveClock(0), FullmoveNumber(1), DEFPLAYERSALIVE}
}

//func (s *State) String() string {   // returns FEN
//}

//func ParsBoard3FEN([]byte) *[8][24][2]byte {
//}

//func Pars3FEN([]byte) *State {
//}
