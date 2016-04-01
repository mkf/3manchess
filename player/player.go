package player

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "sync"
import "log" //debug

//Player is either AI or a human via some UI
type Player interface {
	Initialize(*Gameplay)
	ErrorChannel() chan<- error
	HurryChannel() chan<- bool
	HeyItsYourMove(*game.State, <-chan bool) game.Move //that channel is for signalling to hurry up
	HeySituationChanges(game.Move, *game.State)
	HeyYouLost(*game.State)
	HeyYouWon(*game.State)
	HeyYouDrew(*game.State)
	AreWeWaitingForYou() bool
	HeyWeWaitingForYou(bool)
	String() string
}

type PlayerGen interface {
	Start() error
	GenPlayer(name string) (Player, error)
	String() string
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
	game.Move
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

//GiveResult does *not* run EvalDeath!
func (gp *Gameplay) GiveResult() (breaking bool) {
	for _, ci := range game.COLORS {
		if !gp.State.PlayersAlive.Give(ci) {
			gp.Players[ci].HeyYouLost(gp.State)
		}
	}
	listem := gp.State.PlayersAlive.ListEm()
	switch len(listem) {
	case 1:
		gp.Players[listem[0]].HeyYouWon(gp.State)
		breaking = true
	case 0:
		for _, ci := range game.COLORS { //TODO: Draw only if alive before
			gp.Players[ci].HeyYouDrew(gp.State)
		}
		breaking = true
	}
	return
}

func (gp *Gameplay) Lifes() (end bool) {
	log.Println("LifesFunc")
	gp.State.EvalDeath()
	return gp.GiveResult()
}

func (gp *Gameplay) Turn() (breaking bool) {
	log.Println("TURNFUNC START")
	gp.Players[gp.State.MovesNext].HeyWeWaitingForYou(true)
	hurry := make(chan bool)
	log.Println("TURNFUNC ASKING", gp.Players, gp.State.MovesNext)
	move := gp.Players[gp.State.MovesNext].HeyItsYourMove(gp.State, hurry)
	log.Println("TURNFUNC AFTERING", move)
	after, err := move.After()
	log.Println("TURNFUNC AFTERED", after, err)
	if err != nil {
		gp.Players[gp.State.MovesNext].ErrorChannel() <- err
		log.Println(err)
		return gp.Turn()
	}
	log.Println("TURNFUNC is gonna EVALDEATH")
	after.EvalDeath()
	log.Println("TURNFUNC EVALDEATHED")
	gp.State = after
	var wg sync.WaitGroup
	for _, ci := range game.COLORS {
		wg.Add(1)
		log.Println("NOTIFYING", ci, "ABOUT", move, after)
		gp.Players[ci].HeySituationChanges(move, after)
		wg.Done()
	}
	wg.Wait()
	return gp.GiveResult()
}

func (gp *Gameplay) Procedure(end chan<- bool) {
	log.Println("Procedure")
	if !gp.Lifes() {
		log.Println("Given")
		for !gp.Turn() {
			log.Println("Turning...")
		}
	}
	log.Println("NotTurningAnymore")
	end <- false
}
