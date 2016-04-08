package constsitval

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"
import "github.com/ArchieT/3manchess/game"
import "time"

//func TestHeyItsYourMove_depth_eq_0(t *testing.T) {NewgameAI(t, AIConfig{Depth: 0, OwnedToThreatened: DEFOWN2THRTHD})}
func TestHeyItsYourMove_newgame(t *testing.T) {
	NewgameAI(t, AIConfig{Depth: DEFFIXDEPTH, OwnedToThreatened: DEFOWN2THRTHD})
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

var plat1448 []game.FromTo = []game.FromTo{
	{game.Pos{1, 0}, game.Pos{2, 0}},
	{game.Pos{1, 8}, game.Pos{2, 8}},
	{game.Pos{1, 16}, game.Pos{2, 16}},
	{game.Pos{2, 0}, game.Pos{3, 0}},
	{game.Pos{2, 8}, game.Pos{3, 8}},
	{game.Pos{2, 16}, game.Pos{3, 16}},
	{game.Pos{0, 0}, game.Pos{2, 0}},
	{game.Pos{0, 9}, game.Pos{2, 10}},
	{game.Pos{0, 16}, game.Pos{2, 16}},
}

func TestHeyItsYourMove_plat1448(t *testing.T) {
	var a AIPlayer
	a.Name = "Bot testowy"
	a.Conf = AIConfig{Depth: DEFFIXDEPTH, OwnedToThreatened: DEFOWN2THRTHD};
	hurry := make(chan bool)
	newState := game.NewState()
	var s *game.State
	s = &newState
	var err error
	for _, ft := range plat1448 {
		mov := ft.Move(s)
		s, err = mov.EvalAfter()
		if err != nil {
			t.Error(err, s, mov, ft)
		}
		t.Log(ft, mov, s)
	}
	go func() {
		time.Sleep(time.Minute)
		hurry <- true
	}()
	move := a.HeyItsYourMove(s, hurry)
	t.Log(move)
}
