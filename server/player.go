package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "fmt"

//import "log"

type ServerPlayer struct {
	ID      int64
	errchan chan error
	wwme    bool
	move    chan *game.Move
	sitret  chan player.SituationChange
}

func (sp *ServerPlayer) GiveMove(m *game.Move) <-chan player.SituationChange {
	sp.move <- m
	return sp.sitret
}

func (sp *ServerPlayer) Initialize(*player.Gameplay) { sp.sitret = make(chan player.SituationChange) }
func (sp *ServerPlayer) ErrorChannel() chan<- error  { return sp.errchan }
func (sp *ServerPlayer) HurryChannel() chan<- error {
	c := make(chan error)
	go func() {
		for {
			<-c
		}
	}()
	return c
}
func (sp *ServerPlayer) HeyWeWaitingForYou(b bool) { sp.wwme = b }
func (sp *ServerPlayer) AreWeWaitingForYou() bool  { return sp.wwme }
func (sp *ServerPlayer) String() string            { return fmt.Sprint(sp.ID) }
func (sp *ServerPlayer) HeyYouWon(*game.State)     {}
func (sp *ServerPlayer) HeyYouDrew(*game.State)    {}
func (sp *ServerPlayer) HeyYouLost(*game.State)    {}
func (sp *ServerPlayer) HeySituationChanges(m *game.Move, s *game.State) {
	sp.sitret <- player.SituationChange{m, s}
}
func (sp *ServerPlayer) HeyItsYourMove(_ *game.State, _ <-chan bool) *game.Move { return <-sp.move }
