package player

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	Initialize(*Gameplay)
	ErrorChannel() chan<- error
	HeyItsYourMove(*game.State, <-chan bool) *game.Move //that channel is for signalling to hurry up
	HeySituationChanges(*game.Move, *game.State)
	HeyYouLost(*game.State)
	HeyYouWonOrDrew(*game.State)
	String() string
}

//Gameplay is a list of players and the current gamestate pointer
type Gameplay struct {
	Players map[game.Color]Player
	*game.State
}

type GivingUpError interface {
	error
	IGaveUp() string
}

//NewGame returns a new Gameplay
func NewGame(ourplayers map[game.Color]Player, end chan<- bool) *Gameplay {
	ns := game.NewState()
	gp := Gameplay{ourplayers, &ns}
	for _, ci := range game.COLORS {
		ourplayers[ci].Initialize(&gp)
	}
	go gp.Procedure(end)
	return &gp
}

func (gp *Gameplay) Procedure(end chan<- bool) {
	var move *game.Move
	var after *game.State
	var hurry chan bool
	var listem []game.Color
	var err error
	for {
		hurry = make(chan bool)
		gp.State.EvalDeath()
		for _, ci := range game.COLORS {
			if !gp.State.PlayersAlive.Give(ci) {
				gp.Players[ci].HeyYouLost(gp.State)
			}
		}
		listem = gp.State.PlayersAlive.ListEm()
		if len(listem) == 1 {
			gp.Players[listem[0]].HeyYouWonOrDrew(gp.State)
			break
		}
		if len(listem) == 0 {
			for _, ci := range game.COLORS {
				gp.Players[ci].HeyYouWonOrDrew(gp.State)
			}
			break
		}
		move = gp.Players[gp.State.MovesNext].HeyItsYourMove(gp.State, hurry)
		after, err = move.After()
		if err != nil {
			gp.Players[gp.State.MovesNext].ErrorChannel() <- err
			continue
		}
		gp.State = after
	}
	end <- false
}
