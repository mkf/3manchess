package game

import "strconv"

//FigType : type of a figure
type FigType byte //piece type
var (
	//Pawn FigType  const
	Pawn = FigType('p')
	//Rook FigType   const
	Rook = FigType('r')
	//Knight FigType   const
	Knight = FigType('n')
	//Bishop FigType   const
	Bishop = FigType('b')
	//Queen FigType   const
	Queen = FigType('q')
	//King FigType    const
	King = FigType('k')
)

func (f *Fig) String() string {
	r := f.Rune()
	otoczka := "."
	if f.PawnCenter {
		otoczka = "!"
	}
	return otoczka + string(r)
}

//PawnCenter : whether the pawn had already passed through the center
type PawnCenter bool

//Byte returns "Y" or "N"
func (pc PawnCenter) Byte() byte {
	return ynbool(bool(pc))
}

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

//EMPTYOURSTR is the string that is the value of Square.String() if Square.Empty()
var EMPTYOURSTR = "__"

func (s Square) String() string {
	if s.NotEmpty {
		//ourbyte = []byte{byte(s.Fig.Color), byte(s.Fig.FigType), s.Fig.PawnCenter.Byte()}
		return s.Fig.String()
	}
	return EMPTYOURSTR
}

//Pos : coordinates
type Pos [2]int8

//Pos.String : give a nice [0,0] string
func (p Pos) String() string {
	return "[" + strconv.Itoa(int(p[0])) + "," + strconv.Itoa(int(p[1])) + "]"
}

//Board : board array
type Board [6][24]Square

//GPos gets the position's square pointer
func (b *Board) GPos(p Pos) *Square {
	if err := p.Correct(); err != nil {
		panic(err)
	}
	return &((*b)[p[0]][p[1]])
}

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
	//case 0:
	//	return 127 //Bug(ArchieT): sometimes c==byte(0)
	default:
		//panic(c)
		panic(strconv.Itoa(int(uint8(byte(c)))))
	}
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
		panic(byte(c))
	}
}

var (
	//White color const
	White = Color('W') //white color
	//Gray color const
	Gray = Color('G') //gray color
	//Black color const
	Black = Color('B') //black color
)

//COLORS : array of colors, ordered	const
var COLORS = [3]Color{White, Gray, Black}

//FIRSTRANKNEWGAME : first rank, from 0mod8 to 7mod8
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

//BOARDFORNEWGAME — a newgame board
var BOARDFORNEWGAME Board //newgame board

func boardinit() { //initialize BOARDFORNEWGAME module pseudoconstant
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

//NewBoard is a replacement for NEWBOARD
func NewBoard() Board {
	return BOARDFORNEWGAME
}
