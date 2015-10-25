package game

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

var COLORS = [3]Color{White, Gray, Black}
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, King, Queen, Bishop, Knight, Rook}

var BOARDFORNEWGAME Board

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
