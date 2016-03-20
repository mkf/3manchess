package movedet

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/movedet/board"

//var simplyillegal []testentity

var first = game.Move{From: game.Pos{1, 0}, To: game.Pos{3, 0}, Before: &game.NEWGAME, PawnPromotion: 0}
var temp, _ = first.After()
var g = []game.Move{first, {From: game.Pos{1, 8}, To: game.Pos{3, 8}, Before: temp, PawnPromotion: 0}}
var b = []game.Move{{From: game.Pos{3, 0}, To: game.Pos{4, 0}, Before: temp, PawnPromotion: 0}}

//simplyillegal = []testentity{
//	{Pos{3,0},Pos{5,0},g[0].After()}
//}

func TestGood(t *testing.T) {
	for _, pair := range g {
		temp, err := pair.After()

		t.Log("pair", pair)
		t.Log("pair.Before", pair.Before)
		var v *game.Move
		var w *game.State
		if temp != nil {
			t.Log("temp.Board", *temp.Board)
			v, w, err = WhatMove(pair.Before, board.FromGameBoard(temp.Board))
		} else {
			t.Log("temp is nil")
			if err == nil {
				t.Error("Temp is nil and err is nil")
			}
		}
		if err != nil {
			t.Error("For", pair, "got an error", err, "and value", v, w)
		} else if (v.From != pair.From) || (v.To != pair.To) {
			t.Error("For", pair, "got", v, w, "expected", temp)
		} else {
			t.Log("For", pair, "got", v, w, "   GOOD", "  expected ", temp)
		}
	}
}

func TestSimplyIllegal(t *testing.T) {
	for _, pair := range b {
		temp, terr := pair.After()
		t.Log("pair", pair)
		t.Log("pair.Before", pair.Before)
		var v *game.Move
		var w *game.State
		if temp == nil {
			t.Log("temp is nil")
			if terr == nil {
				t.Error("Temp is nil and err is nil")
			}
		} else {
			t.Log("temp IS NOT NIL!!")
			t.Log("temp.Board", *temp.Board)
			v, w, terr = WhatMove(pair.Before, board.FromGameBoard(temp.Board))
		}
		if terr == nil {
			t.Error("For", pair, "there is no error, values are", v, w)
		} else {
			switch _, ok := terr.(game.IllegalMoveError); ok {
			case true:
				t.Log("For", pair, "got an IllegalMoveError", terr)
			case false:
				t.Error("For", pair, "got an error", terr, "of type other than IllegalMoveError")
			}
		}
	}
}
