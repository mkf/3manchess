//Package board provides a straightforward representation of a board, not containing
//the PawnCenter bool in its structs, which could be misleading if omitted.
package board

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"

type Piece struct {
	game.FigType
	game.Color
}

func (p *Piece) Equal(gf *game.Fig) bool {
	return p.Color == gf.Color && p.FigType == gf.FigType
}

type Square struct {
	NotEmpty bool
	Piece
}

func (s *Square) Equal(gs *game.Square) bool {
	return s.Piece.Equal(&(gs.Fig)) && s.NotEmpty == gs.NotEmpty
}

func (s *Square) Empty() bool {
	return !s.NotEmpty
}

type Pos game.Pos

func (p *Pos) Correct() error {
	return game.Pos(*p).Correct()
}

type Board [6][24]Square

func FromGameBoard(gb *game.Board) *Board {
	var newb Board
	var gsq *game.Square
	for _, oac := range game.ALLPOS {
		gsq = gb.GPos(oac)
		newb[oac[0]][oac[1]] = Square{gsq.NotEmpty, Piece{gsq.FigType, gsq.Color()}}
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
	var gs *game.Square
	var os *Square
	for _, oac := range game.ALLPOS {
		gs = gb.GPos(oac)
		os = b.GPos(Pos(oac))
		if !os.Equal(gs) {
			return false
		}
	}
	return true
}
