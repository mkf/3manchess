//Package board provides a straightforward representation of a board, not containing
//the PawnCenter bool in its structs, which could be misleading when ommited.
package board

import "github.com/ArchieT/3manchess/game"

type Piece struct {
	game.FigType
	game.Color
}

type Square struct {
	NotEmpty bool
	Piece
}

type Pos game.Pos

type Board [6][24]Square

func FromGameBoard(gb *game.Board) *Board {
	var oac game.ACP
	var newb Board
	var gsq game.Square
	for oac.OK() {
		gsq = gb.GPos(oac)
		newb[oac[0]][oac[1]] = Square{gsq.NotEmpty, Piece{gsq.FigType, gsq.Color()}}
		oac.P()
	}
	return &newb
}

func (b *Board) GPos(p Pos) *Square {
	if err := p.Correct(); err != nil {
		panic(err)
	}
	return &((*b)[p[0]][p[1]])
}

func (b *Board) Equal(gb *game.Board) bool {
	var oac game.ACP
	var gs *game.Square
	var os *Square
	for oac.OK() {
		gs = gb.GPos(game.Pos(oac))
		os = b.GPos(Pos(oac))
		if gs.Piece.FigType != os.Fig.FigType || gs.Piece.Color != os.Fig.Color || gs.NotEmpty != os.NotEmpty {
			return false
		}
		oac.P()
	}
	return true
}
