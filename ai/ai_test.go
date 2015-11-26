package ai

import "testing"
import "github.com/ArchieT/3manchess/game"

func TestSimpleGenNoPanic(t *testing.T) {
	var a AIsettings
	a.Think(&game.NEWGAME)
}
