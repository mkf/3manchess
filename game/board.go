package game

import "strconv"

//FigType : type of a figure
type FigType byte //piece type

//Pawn FigType  const
var Pawn = FigType('p')

//Rook FigType   const
var Rook = FigType('r')

//Knight FigType   const
var Knight = FigType('n')

//Bishop FigType   const
var Bishop = FigType('b')

//Queen FigType   const
var Queen = FigType('q')

//King FigType    const
var King = FigType('k')

//PawnCenter : whether the pawn had already passed through the center
type PawnCenter bool

//Fig : a struct describing a single piece: it's type, it's color, and, in case of a pawn, whether is had already passed through the center
type Fig struct {
	FigType
	Color
	PawnCenter
}

//Square : a struct describing a single square: whether it is empty, and what is on it
type Square struct {
	Fig
	NotEmpty bool
}

//Empty : return !s.NotEmpty
func (s Square) Empty() bool {
	return !s.NotEmpty
}

//Color : legacy : return s.Fig.Color
func (s Square) Color() Color {
	return s.Fig.Color
}

//What : legacy : return s.Fig.FigType
func (s Square) What() FigType {
	return s.Fig.FigType
}

//Pos : coordinates
type Pos [2]int8

//Pos.String : give a nice [0,0] string
func (p Pos) String() string {
	return "[" + strconv.Itoa(int(p[0])) + "," + strconv.Itoa(int(p[1])) + "]"
}

//Board : board array
type Board [6][24]Square

//Color : color type
type Color byte

//UInt8 : returns 0 for white, 1 for gray, 2 for black
func (c Color) UInt8() uint8 {
	switch c {
	case Color('W'), Color('w'):
		return 0
	case Color('G'), Color('g'):
		return 1
	case Color('B'), Color('b'):
		return 2
	}
	//panic(c)
	panic("replacementpanic2")
}

//ColorUint8 : returns White for 0, Gray for 1, Black for 2
func ColorUint8(u uint8) Color {
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

//Next returns the next color: White, Gray, Black,  White, etc.
func (c Color) Next() Color {
	return ColorUint8((c.UInt8() + 1) % 3)
}

//String returns string "White"/"Gray"/"Black"
func (c Color) String() string {
	switch c.UInt8() {
	case 0:
		return "White"
	case 1:
		return "Gray"
	case 2:
		return "Black"
	default:
		panic(c.UInt8())
	}
}

//White color const
var White = Color('W') //white color
//Gray color const
var Gray = Color('G') //gray color
//Black color const
var Black = Color('B') //black color

//COLORS : array of colors, ordered	const
var COLORS = [3]Color{White, Gray, Black}

//FIRSTRANKNEWGAME : first rank, from 0mod8 to 7mod8
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

//BOARDFORNEWGAME — a newgame board
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
