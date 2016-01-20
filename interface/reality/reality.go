package reality

import "github.com/ArchieT/3manchess/interface/reality/machine"
import "github.com/ArchieT/3manchess/interface/reality/camget"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"

//import "log"

type Reality struct {
	camget.View
	machine.Machine
	BlackIsOnWhitesRight bool
}

type RealPlayer struct {
	*Reality
	color     game.Color
	Name      string
	errchan   chan error
	ErrorChan chan<- error
	HurryChan chan<- bool
	hurry     chan bool
	gp        *player.Gameplay
}

func NewReality() *Reality {
	return new(Reality)
}

func (re *Reality) MakePlayers(who ...game.Color) map[game.Color]*RealPlayer {
	ourm := make(map[game.Color]*RealPlayer)
	var ourp *RealPlayer
	for _, c := range who {
		ourp = new(RealPlayer)
		ourp.Reality = re
		ourp.Name = c.String() + "_real"
		ourm[c] = ourp
	}
	return ourm
}

func (p *RealPlayer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	p.errchan = errchan
	hurry := make(chan bool)
	p.hurry = hurry
	p.gp = gp
	p.ErrorChan = errchan
	p.HurryChan = hurry
}

func (p *RealPlayer) String() string { return p.Name }

func (p *RealPlayer) ErrorChannel() chan<- error { return p.ErrorChan }

func (p *RealPlayer) HurryChannel() chan<- bool { return p.HurryChan }

func (p *RealPlayer) HeyItsYourMove(s *game.State, hurryi <-chan bool) *game.Move {
	go func() {
		for {
			p.hurry <- <-hurryi
		}
	}()
	return nil
}
