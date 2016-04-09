package movedet

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/movedet/board"

type changefromempty struct {
	what  board.Piece
	where board.Pos
}

type changetoempty struct {
	what  game.Fig
	where game.Pos
}

type changereplace struct {
	before game.Fig
	after  board.Piece
	where  board.Pos
}

//IllegalMoveDetected : error dtruct containing a nice description and a codename string
type IllegalMoveDetected struct {
	Description string
	Codename    string
}

func (i IllegalMoveDetected) Error() string {
	return i.Description
}

//WhatMove : return a move that happened between the before state and the board after
func WhatMove(bef *game.State, aft *board.Board) (*game.Move, *game.State, error) {
	//yep, it's right! all over, again!
	//but now... with concurrency!
	var om *game.Move
	var os *game.State
	var ourmove game.Move
	//var whatafter game.State
	appeared := make([]changefromempty, 0, 2)
	disappeared := make([]changetoempty, 0, 2)
	replaced := make([]changereplace, 0, 1)
	for _, oac := range game.ALLPOS {
		prev := bef.Board.GPos(game.Pos(oac))
		next := aft.GPos(board.Pos(oac))
		if prev.Empty() && next.NotEmpty {
			appeared = append(appeared, changefromempty{next.Piece, board.Pos(oac)})
		} else if next.Empty() && prev.NotEmpty {
			disappeared = append(disappeared, changetoempty{prev.Fig, game.Pos(oac)})
		} else if next.NotEmpty && prev.NotEmpty && !(next.Piece.Equal(&prev.Fig)) {
			replaced = append(replaced, changereplace{prev.Fig, next.Piece, board.Pos(oac)})
		} else {
			//panic([2]game.Board{*bef.Board, *aft})
			//panic("replacementpanic1")
		}
	}
	if len(replaced) > 1 {
		return om, os, IllegalMoveDetected{"Too many replaced pieces!", "TooManyReplaced"}
	}
	if len(appeared) > len(disappeared) {
		return om, os, IllegalMoveDetected{"More appeared than disappeared!", "MoreAppearedThanDisappeared"}
	}
	if (len(appeared) == 2) && (len(disappeared) == 2) && (len(replaced) == 0) {
		var aking changefromempty
		var dking changetoempty
		if appeared[0].what.FigType == game.King {
			aking = appeared[0]
		} else if appeared[1].what.FigType == game.King {
			aking = appeared[1]
		} else {
			return om, os, IllegalMoveDetected{"It ain't no castling!", "NotACastling"}
		}
		if disappeared[0].what.FigType == game.King {
			dking = disappeared[0]
		} else if disappeared[1].what.FigType == game.King {
			dking = disappeared[1]
		} else {
			return om, os, IllegalMoveDetected{"It ain't no castling, though there was a king!", "NotACastlingButKing"}
		}
		ourmove = game.Move{From: dking.where, To: game.Pos(aking.where), Before: bef}
		whatafter, err := ourmove.After()
		if err != nil {
			return &ourmove, whatafter, err
		}
		if !aft.Equal(whatafter.Board) {
			panic("Legal move, yet the effect is different from what we've got on input???")
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) > 2 {
		return om, os, IllegalMoveDetected{"More than 2 disappeared!", "MoreThanTwoDisappeared"}
	}
	if len(appeared) == 1 && len(disappeared) == 2 && len(replaced) == 0 && appeared[0].what.FigType == game.Pawn {
		for _, j := range disappeared {
			ourmove = game.Move{From: j.where, To: game.Pos(appeared[0].where), Before: bef, PawnPromotion: appeared[0].what.FigType}
			whatafter, err := ourmove.After()
			if err == nil {
				if !aft.Equal(whatafter.Board) {
					panic("Legal move, yet the effect is different from what we've got on input???")
				}
			}
			return &ourmove, whatafter, err
		}
	}
	if (len(appeared) > 0 || len(disappeared) > 1) && len(replaced) > 0 {
		return om, os, IllegalMoveDetected{"Both (dis)appeared and been replaced!", "BothAppearDisappearAndReplace"}
	}
	if len(disappeared) == 1 && len(replaced) == 0 && len(appeared) == 0 {
		return om, os, IllegalMoveDetected{"One disappearance with no reason!", "NoReasonDisappear"}
	}
	if len(disappeared) == 1 && len(replaced) == 1 {
		ourmove = game.Move{From: disappeared[0].where, To: game.Pos(replaced[0].where), Before: bef, PawnPromotion: replaced[0].after.FigType}
		whatafter, err := ourmove.After()
		if err == nil {
			if !aft.Equal(whatafter.Board) {
				panic("Legal move, yet the effect is different from what we've got on input???")
			}
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) == 1 && len(appeared) == 1 {
		ourmove = game.Move{From: disappeared[0].where, To: game.Pos(appeared[0].where), Before: bef, PawnPromotion: appeared[0].what.FigType}
		whatafter, err := ourmove.After()
		if err == nil {
			if !aft.Equal(whatafter.Board) {
				panic("Legal move, yet the effect is different from what we've got on input???")
			}
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) == 0 && len(appeared) == 0 && len(replaced) == 0 {
		return om, bef, IllegalMoveDetected{"The board remains unchanged", "Unchanged"}
	}
	if len(disappeared) == 0 && len(appeared) == 0 && len(replaced) == 1 {
		return om, os, IllegalMoveDetected{"One replacement with no reason (Note: promotion should be done in the same move)!", "NoReasonReplace"}
	}
	panic("None of the cases???")
}
