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

func WhatMove(bef *game.State, aft *game.Board) (*game.Move, *game.State, error) {
	//yep, it's right! all over, again!
	//but now... with concurrency!
	var i, j int8
	appeared := make([]changeempty, 0, 4)
	disappeared := make([]changeempty, 0, 4)
	replaced := make([]changereplace, 0, 2)
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			prev = bef.Board[i][j]
			next = (*aft)[i][j]
			if prev.Empty() && next.NotEmpty {
				appeared = append(appeared, changeempty{next.Fig, game.Pos{i, j}})
			} else if next.Empty() && prev.NotEmpty {
				disappeared = append(disappeared, changeempty{next.Fig, game.Pos{i, j}})
			} else if next.NotEmpty && prev.NotEmpty {
				replaced = append(replaced, changereplace{prev.Fig, next.Fig, game.Pos{i, j}})
			} else {
				panic([2]Board{bef.Board, *aft})
			}
		}
	}
}
