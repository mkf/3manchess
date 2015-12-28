package ai

import "testing"
import "github.com/ArchieT/3manchess/game"
import "time"
import "log"

func TestSimpleGenNoPanic(t *testing.T) {
	var a AIPlayer
	hurry := make(chan bool)
	newgame := game.NewState()
	go func() {
		time.Sleep(time.Minute)
		hurry <- true
	}()
	move := a.HeyItsYourMove(&newgame, hurry)
	log.Println(move)
}
