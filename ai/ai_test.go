package ai

import "testing"
import "github.comArchieT/3manchess/game"

func TestSimpleGenNoPanic(t *testing.T) {
	var a AIsettings
	a.Think(&game.NEWGAME)
}
