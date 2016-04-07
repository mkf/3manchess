package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"

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
