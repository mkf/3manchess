package devengchan

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"

//import "github.com/ArchieT/3manchess/simple"
import "fmt"
import "log"

type ResultCode int8

const (
	UNDEFRESULT ResultCode = 0
	WIN         ResultCode = -2
	DRAW        ResultCode = 1
	LOSE        ResultCode = 2
)

type Developer struct {
	Name         string
	errchan      chan error
	ErrorChan    chan<- error
	HurryChan    chan bool
	hurry        chan bool
	gp           *player.Gameplay
	waiting      bool
	askformove   chan<- *game.State
	AskinForMove <-chan *game.State
	heremoves    <-chan *game.Move
	HereRMoves   chan<- *game.Move
	Result       <-chan ResultCode
	sendresult   chan<- ResultCode
	SituationCh  <-chan player.SituationChange
	sitchan      chan<- player.SituationChange
}

func (p *Developer) AskingForMove() <-chan *game.State {
	return p.AskinForMove
}

func (p *Developer) HereAreMoves() chan<- *game.Move {
	return p.HereRMoves
}

func (p *Developer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	p.errchan = errchan
	hurry := make(chan bool)
	p.hurry = hurry
	fmt.Printf("%s initialized with Gameplay:\n", p)
	fmt.Println(gp)
	fmt.Println("")
	p.gp = gp
	p.ErrorChan = errchan
	p.HurryChan = hurry
	sres := make(chan ResultCode)
	p.Result = sres
	p.sendresult = sres
	afm := make(chan *game.State)
	p.askformove = afm
	p.AskinForMove = afm
	hrm := make(chan *game.Move)
	p.heremoves = hrm
	p.HereRMoves = hrm
	sch := make(chan player.SituationChange)
	p.SituationCh = sch
	p.sitchan = sch
	go p.logger()
}

func (p *Developer) logger() {
	var err error
	for {
		err = <-p.errchan
		log.Println(err)
	}
}

func (p *Developer) String() string { return p.Name }

func (p *Developer) ErrorChannel() chan<- error { return p.ErrorChan }

func (p *Developer) HurryChannel() chan<- bool { return p.HurryChan }

func (p *Developer) HeyItsYourMove(s *game.State, hurryi <-chan bool) *game.Move {
	go func() {
		for {
			p.hurry <- <-hurryi
		}
	}()
	p.askformove <- s
	move := <-p.heremoves
	if move.Before != s {
		return p.HeyItsYourMove(s, hurryi)
	}
	p.HeyWeWaitingForYou(false)
	return move
}

func (p *Developer) HeySituationChanges(m *game.Move, aft *game.State) {
	p.sitchan <- player.SituationChange{m, aft}
}

func (p *Developer) HeyYouLost(_ *game.State) { p.sendresult <- LOSE }

func (p *Developer) HeyYouWon(_ *game.State) { p.sendresult <- WIN }

func (p *Developer) HeyYouDrew(_ *game.State) { p.sendresult <- DRAW }

func (p *Developer) HeyWeWaitingForYou(b bool) { p.waiting = b }

func (p *Developer) AreWeWaitingForYou() bool { return p.waiting }

type GivingUpError string

func (g GivingUpError) Error() string { return string(g) }

func (g GivingUpError) IGaveUp() string { return string(g) }
