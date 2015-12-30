package sitvalues

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/simple"
import "github.com/ArchieT/3manchess/player"
import "sync"
import "sync/atomic"
import "fmt"

type AIPlayer struct {
	errchan           chan error
	ErrorChan         chan<- error
	hurry             chan bool
	HurryChan         chan<- bool
	FixedPrecision    float64
	OwnedToThreatened float64
	gp                *player.Gameplay
	waiting           bool
}

func (a *AIPlayer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	a.errchan = errchan
	a.ErrorChan = errchan
	hurry := make(chan bool)
	a.hurry = hurry
	a.HurryChan = hurry
	a.gp = gp
	go func() {
		for b := range a.errchan {
			panic(b)
		}
	}()
}

func (a *AIPlayer) HurryChannel() chan<- bool {
	return a.HurryChan
}

func (a *AIPlayer) ErrorChannel() chan<- error {
	return a.ErrorChan
}
