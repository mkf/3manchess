package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"

var plat129 []FromTo = []FromTo{
	{Pos{1, 0}, Pos{2, 0}},
	{Pos{1, 8}, Pos{2, 8}},
	{Pos{1, 16}, Pos{2, 16}},
	{Pos{2, 0}, Pos{3, 0}},
	{Pos{2, 8}, Pos{3, 8}},
	{Pos{2, 16}, Pos{3, 16}},
	{Pos{0, 0}, Pos{2, 0}},
	{Pos{0, 8}, Pos{2, 8}},
	{Pos{0, 22}, Pos{2, 23}},
	{Pos{2, 0}, Pos{2, 8}},
	{Pos{0, 9}, Pos{2, 10}},
	{Pos{0, 16}, Pos{2, 16}},
	{Pos{2, 8}, Pos{2, 15}},
/*	{Pos{0, 14}, Pos{2, 13}},
	{Pos{2, 16}, Pos{2, 15}},
	{Pos{0, 1}, Pos{2, 0}},
	{Pos{0, 11}, Pos{0, 8}},
	{Pos{2, 15}, Pos{2, 7}},
	{Pos{0, 6}, Pos{2, 5}},
	{Pos{2, 10}, Pos{1, 8}},
	{Pos{2, 7}, Pos{3, 7}},
	{Pos{2, 0}, Pos{3, 2}},
	{Pos{1, 8}, Pos{3, 7}},
	{Pos{2, 23}, Pos{1, 1}},
	{Pos{2, 5}, Pos{3, 7}},
	{Pos{0, 8}, Pos{2, 8}},
	{Pos{1, 1}, Pos{2, 3}},*/
}

func TestSimpleGenNoPanic(t *testing.T) {
	newState := NewState()
	first := Move{Pos{1, 0}, Pos{3, 0}, &newState, 0}
	temp, err := first.After()
	t.Log("first.After ", temp)
	if err != nil {
		t.Error("Error first.After ", err)
	}
	second := Move{From: Pos{3, 0}, To: Pos{4, 0}, Before: temp, PawnPromotion: 0}
	stemp, err := second.After()
	t.Log("second.After ", stemp)
	if err != nil {
		t.Log("Error second.After ", err)
	}
	if err == nil {
		t.Error("Error second.After IS NIL!")
	}
	third := Move{From: Pos{3, 8}, To: Pos{4, 8}, Before: temp, PawnPromotion: 0}
	ttemp, err := third.After()
	t.Log("third.After ", ttemp)
	if err != nil {
		t.Log("Error third.After ", err)
	}
	if err == nil {
		t.Error("Error third.After IS NIL!")
	}
	forth := Move{From: Pos{1, 8}, To: Pos{3, 8}, Before: temp, PawnPromotion: 0}
	ftemp, err := forth.After()
	t.Log("forth.After ", ftemp)
	if err != nil {
		t.Error("Error forth.After ", err)
	}
}

func TestAfter_pawnCrossCenter(t *testing.T) {
	var i, j int8
	newState := NewState()
	move := Move{Pos{1, i * 8}, Pos{3, i * 8}, &newState, 0}
	statePointer, err := move.After()
	for i = 1; i < 3; i++ {
		move := Move{Pos{1, i * 8}, Pos{3, i * 8}, statePointer, 0}
		statePointer, err = move.After()
		if err != nil {
			t.Error("Unexpected error (1st loop, i =", i, "): ", err)
		}
	}
	t.Log("state after 3 moves: ", statePointer)
	for i = 3; i < 5; i++ {
		for j = 0; j < 3; j++ {
			move := Move{Pos{i, j * 8}, Pos{i + 1, j * 8}, statePointer, 0}
			statePointer, err = move.After()
			if err != nil {
				t.Error("Unexpected error (2nd loop, i =", i, ", j =", j, "): ", err)
			}
		}
		t.Log("state after", (i-1)*3, "moves: ", statePointer)
	}
	for i = 0; i < 3; i++ {
		move := Move{Pos{5, i * 8}, Pos{5, (i*8 + 12) % 24}, statePointer, 0}
		statePointer, err = move.After()
		if err != nil {
			t.Error("Unexpected error (3rd loop, i =", i, "): ", err)
		}
	}
}

func TestEvalAfter_plat129(t *testing.T) {
	newState := NewState()
	var s *State
	s = &newState
	var err error
	var mov Move
	for _, ft := range plat129 {
		if s == nil {
			t.Error("Move considered invalid:", mov)
		}
		mov = ft.Move(s)
		s, err = mov.EvalAfter()
		if err == nil {
			t.Log(ft, mov, s)
		}
	}
	if err == nil {
		t.Error("Invalid move accepted. State afterwards:", s)
	}
}

func TestAfter_pawnCapture(t *testing.T) {
	newState := NewState()
	var s *State
	s = &newState
	var mov Move
	var ft FromTo
	var err error
	for i := 0; i < 10; i++ {
		ft = plat129[i]
		mov = ft.Move(s)
		s, err = mov.EvalAfter()
		if err != nil {
			t.Error(err, s, mov, ft)
		}
		t.Log(ft, mov, s)
	}
	s.Board[2][10].NotEmpty = true
	s.Board[2][10].Fig.FigType = Pawn
	s.Board[2][10].Fig.Color = White
	s.Board[2][10].Fig.PawnCenter = true
	t.Log(s)
	ft = FromTo{Pos{1, 9}, Pos{2, 10}}
	mov = ft.Move(s)
	if s, err = mov.After(); err != nil {
		t.Error(err, s, mov, ft)
	}
	t.Log(ft, mov, s)
	s.Board[2][18].NotEmpty = true
	s.Board[2][18].Fig.FigType = Pawn
	s.Board[2][18].Fig.Color = White
	s.Board[2][18].Fig.PawnCenter = true
	t.Log(s)
	ft = FromTo{Pos{1, 19}, Pos{2, 18}}
	mov = ft.Move(s)
	if s, err = mov.After(); err != nil {
		t.Error(err, s, mov, ft)
	}
}
