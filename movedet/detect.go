package movedet

import "github.com/ArchieT/3manchess/game"

type changeempty struct {
	what  game.Fig
	where game.Pos
}

type changereplace struct {
	before game.Fig
	after  game.Fig
	where  game.Pos
}

type IllegalMoveDetected struct {
	description string
	codename    string
}

func (i IllegalMoveDetected) Error() string {
	return i.description
}

func WhatMove(bef *game.State, aft *game.Board) (*game.Move, *game.State, error) {
	//yep, it's right! all over, again!
	//but now... with concurrency!
	var i, j int8
	var om *game.Move
	var os *game.State
	var ourmove game.Move
	//var whatafter game.State
	appeared := make([]changeempty, 0, 2)
	disappeared := make([]changeempty, 0, 2)
	replaced := make([]changereplace, 0, 1)
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			prev := bef.Board[i][j]
			next := (*aft)[i][j]
			if prev.Empty() && next.NotEmpty {
				appeared = append(appeared, changeempty{next.Fig, game.Pos{i, j}})
			} else if next.Empty() && prev.NotEmpty {
				disappeared = append(disappeared, changeempty{next.Fig, game.Pos{i, j}})
			} else if next.NotEmpty && prev.NotEmpty {
				replaced = append(replaced, changereplace{prev.Fig, next.Fig, game.Pos{i, j}})
			} else {
				panic([2]game.Board{*bef.Board, *aft})
			}
		}
	}
	if len(replaced) > 1 {
		return om, os, IllegalMoveDetected{"Too many replaced pieces!", "TooManyReplaced"}
	}
	if len(appeared) > len(disappeared) {
		return om, os, IllegalMoveDetected{"More appeared than disappeared!", "MoreAppearedThanDisappeared"}
	}
	if (len(appeared) == 2) && (len(disappeared) == 2) && (len(replaced) == 0) {
		var aking, dking changeempty
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
		ourmove = game.Move{dking.where, aking.where, bef}
		whatafter, err := ourmove.After()
		if err != nil {
			return &ourmove, whatafter, err
		}
		if *whatafter.Board != *aft {
			panic("Legal move, yet the effect is different from what we've got on input???")
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) > 2 {
		return om, os, IllegalMoveDetected{"More than 2 disappeared!", "MoreThanTwoDisappeared"}
	}
	if len(appeared) == 1 && len(disappeared) == 2 && len(replaced) == 0 && appeared[0].what.FigType == game.Pawn {
		for _, j := range disappeared {
			ourmove = game.Move{j.where, appeared[0].where, bef}
			whatafter, err := ourmove.After()
			if err == nil {
				if *whatafter.Board != *aft {
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
		ourmove = game.Move{disappeared[0].where, replaced[0].where, bef}
		whatafter, err := ourmove.After()
		if err == nil {
			if *whatafter.Board != *aft {
				panic("Legal move, yet the effect is different from what we've got on input???")
			}
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) == 1 && len(appeared) == 1 {
		ourmove = game.Move{disappeared[0].where, appeared[0].where, bef}
		whatafter, err := ourmove.After()
		if err == nil {
			if *whatafter.Board != *aft {
				panic("Legal move, yet the effect is different from what we've got on input???")
			}
		}
		return &ourmove, whatafter, err
	}
	if len(disappeared) == 0 && len(appeared) == 0 && len(replaced) == 0 {
		return om, bef, IllegalMoveDetected{"The board remains unchanged", "Unchanged"}
	}
	if len(disappeared) == 0 && len(appeared) == 0 && len(replaced) == 1 {
		return om, os, IllegalMoveDetected{"One replacement with no reason!", "NoReasonReplace"}
	}
	panic("None of the cases???")
}
