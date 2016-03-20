package constsitval

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"
import "github.com/ArchieT/3manchess/game"
import "time"

func TestHeyItsYourMove_depth_eq_0(t *testing.T) {
	NewgameAI(t, AIConfig{Depth: 0, OwnedToThreatened: DEFOWN2THRTHD})
}

func NewgameAI(t *testing.T, acf AIConfig) {
	var a AIPlayer
	a.Name = "Bot testowy"
	a.Conf = acf
	hurry := make(chan bool)
	newgame := game.NewState()
	go func() {
		time.Sleep(time.Minute)
		hurry <- true
	}()
	move := a.HeyItsYourMove(&newgame, hurry)
	t.Log(move)
}
