package game

import "testing"

func TestSimpleGenNoPanic(t *testing.T) {
	first := Move{Pos{1, 0}, Pos{3, 0}, &NEWGAME}
	temp, err := first.After()
	t.Log("first.After ", temp)
	if err != nil {
		t.Log("Error first.After ", err)
	}
	second := Move{From: Pos{3, 0}, To: Pos{4, 0}, Before: temp}
	stemp, err := second.After()
	t.Log("second.After ", stemp)
	if err != nil {
		t.Log("Error second.After ", err)
	}
	third := Move{From: Pos{3, 8}, To: Pos{4, 8}, Before: temp}
	ttemp, err := third.After()
	t.Log("third.After ", ttemp)
	if err != nil {
		t.Log("Error third.After ", err)
	}

}
