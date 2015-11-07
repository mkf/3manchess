package game

type FigType byte //piece type

var Pawn = FigType('p') //pawn typw
var Rook = FigType('r') //rook type
var Knight = FigType('n') //knight type
var Bishop = FigType('b') //bishop type
var Queen = FigType('q') //queen type
var King = FigType('k') //king type

type PawnCenter bool //whether the pawn had already passed through the center

type Fig struct { //a struct describing a single piece: it's type, it's color, and, in case of a pawn, whether is had already passed through the center
	FigType
	Color
	PawnCenter
}

type Square struct { //a struct describing a single square: whether it is empty, and what is on it
	Fig
	NotEmpty bool
}

func (s Square) Empty() bool { //return !s.NotEmpty
	return !s.NotEmpty
}

func (s Square) Color() Color { //return s.Fig.Color
	return s.Fig.Color
}

func (s Square) What() FigType { //return s.Fig.FigType
	return s.Fig.FigType
}

type Pos [2]int8 //coordinates

type Board [6][24]Square //board array

type Color byte //color type

func (c Color) UInt8() uint8 {  //returns 0 for white, 1 for gray, 2 for black
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

func ColorUint8(u uint8) Color {  //returns White for 0, Gray for 1, Black for 2
	switch u {
	case 0:
		return White
	case 1:
		return Gray
	case 2:
		return Black
	default:
		panic(u)
	}
}

var White = Color('W') //white color
var Gray = Color('G') //gray color
var Black = Color('B') //black color

var COLORS = [3]Color{White, Gray, Black}  //array of colors, ordered
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook} //first rank, from 0mod8 to 7mod8

var BOARDFORNEWGAME Board //newgame board

func boardinit() {
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
}
