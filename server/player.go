package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "fmt"

//import "log"

type ServerPlayer struct {
	ID      int64
	errchan chan error
	wwme    bool
	move    chan game.Move
	sitret  chan *game.State
}

func (sp *ServerPlayer) GiveMove(m game.Move) <-chan *game.State {
	sp.move <- m
	return sp.sitret
}

func (sp *ServerPlayer) Initialize(*player.Gameplay) { sp.sitret = make(chan *game.State) }
func (sp *ServerPlayer) ErrorChannel() chan<- error  { return sp.errchan }
func (sp *ServerPlayer) HurryChannel() chan<- bool {
	c := make(chan bool)
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
func (sp *ServerPlayer) HeySituationChanges(_ game.Move, s *game.State) {
	sp.sitret <- s
}
func (sp *ServerPlayer) HeyItsYourMove(_ *game.State, _ <-chan bool) game.Move { return <-sp.move }

func MoveIt(m *game.Move, ids [3]int64) *game.State {
	w := ServerPlayer{ID: ids[0]}
	g := ServerPlayer{ID: ids[1]}
	b := ServerPlayer{ID: ids[2]}
	p := player.Gameplay{
		Players: map[game.Color]player.Player{game.White: &w, game.Gray: &g, game.Black: &b},
		State:   m.Before}
	w.Initialize(&p)
	g.Initialize(&p)
	b.Initialize(&p)
	r := map[game.Color]*ServerPlayer{game.White: &w, game.Gray: &g, game.Black: &b}
	go func() {
		p.Turn()
	}()
	return <-r[m.Before.MovesNext].GiveMove(*m) //TODO: avoid deadlock if somethings not OK with (m *game.Move)
}
