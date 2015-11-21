package game

import "testing"

func TestSimpleGenNoPanic(t *testing.T) {
	first := Move{Pos{1, 0}, Pos{3, 0}, &game.NEWGAME}
	temp, err := first.After()
	t.Log("first.After ", temp)
	if err != nil {
		t.Error("Error first.After ", err)
	}
	second := Move{From: Pos{3, 0}, To: Pos{4, 0}, Before: temp}
	stemp, err := second.After()
	t.Log("second.After ", stemp)
	if err != nil {
		t.Error("Error second.After ", err)
	}
}
