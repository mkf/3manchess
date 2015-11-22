package movedet

import "testing"
import "github.com/ArchieT/3manchess/game"

//var simplyillegal []testentity

var first = game.Move{game.Pos{1, 0}, game.Pos{3, 0}, &game.NEWGAME}
var temp, _ = first.After()
var g = []game.Move{first, game.Move{From: game.Pos{3, 0}, To: game.Pos{4, 0}, Before: temp}}

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
			v, w, err = WhatMove(pair.Before, temp.Board)
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

//func TestSimplyIllegal
