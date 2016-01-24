package player

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	Initialize(*Gameplay)
	ErrorChannel() chan<- error
	HurryChannel() chan<- bool
	HeyItsYourMove(*game.State, <-chan bool) *game.Move //that channel is for signalling to hurry up
	HeySituationChanges(*game.Move, *game.State)
	HeyYouLost(*game.State)
	HeyYouWon(*game.State)
	HeyYouDrew(*game.State)
	AreWeWaitingForYou() bool
	HeyWeWaitingForYou(bool)
	String() string
	Map() map[string]interface{}
	FromMap(map[string]interface{})
	Data() PlayerData
	FromData(PlayerData)
}

type PlayerGen interface {
	Start() error
	GenPlayer(name string) (Player, error)
	String() string
}

type PlayerData struct {
	WhoAmI      string
	Name        string
	Precision   float64
	Coefficient float64
	Promotion   int8
}

//Gameplay is a list of players and the current gamestate pointer
type Gameplay struct {
	Players map[game.Color]Player
	*game.State
}

func (gp *Gameplay) HurryUpWhoever() {
	for _, color := range game.COLORS {
		if gp.Players[color].AreWeWaitingForYou() {
			gp.Players[color].HurryChannel() <- true
		}
	}
}

type SituationChange struct {
	*game.Move
	After *game.State
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
			gp.Players[listem[0]].HeyYouWon(gp.State)
			break
		}
		if len(listem) == 0 {
			for _, ci := range game.COLORS {
				gp.Players[ci].HeyYouDrew(gp.State)
			}
			break
		}
		gp.Players[gp.State.MovesNext].HeyWeWaitingForYou(true)
		move = gp.Players[gp.State.MovesNext].HeyItsYourMove(gp.State, hurry)
		after, err = move.After()
		if err != nil {
			gp.Players[gp.State.MovesNext].ErrorChannel() <- err
			continue
		}
		gp.State = after
		for _, ci := range game.COLORS {
			gp.Players[ci].HeySituationChanges(move, after)
		}
	}
	end <- false
}
