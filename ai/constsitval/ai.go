package constsitval

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/simple"
import "github.com/ArchieT/3manchess/player"

//import "sync"
//import "sync/atomic"
import "fmt"
import "github.com/ArchieT/3manchess/ai"
import "encoding/json"

const MAXDEPTHCONSIDERED int8 = 8 // should be renamed to MINDEPTHNOTCONSIDERED

const DEFFIXDEPTH int8 = 1

const DEFOWN2THRTHD = 4.0

const WhoAmI string = "3manchess-ai_constsitval"

type AIPlayer struct {
	Name    string
	errchan chan error
	hurry   chan bool
	Conf    AIConfig
	gp      *player.Gameplay
	waiting bool
}

func (a *AIPlayer) Config() ai.Config {
	return a.Conf
}

type AIConfig struct {
	Depth             int8
	OwnedToThreatened float64
}

func (c AIConfig) Byte() []byte {
	o, e := json.Marshal(c)
	if e != nil {
		panic(e)
	}
	return o
}

func (c AIConfig) String() string {
	return string(c.Byte())
}

func (a *AIPlayer) Initialize(gp *player.Gameplay) {
	if a.Conf.Depth == 0 {
		a.Conf.Depth = DEFFIXDEPTH
	}
	if a.Conf.OwnedToThreatened == 0.0 {
		a.Conf.OwnedToThreatened = DEFOWN2THRTHD
	}
	a.gp = gp
	a.errchan = make(chan error)
	a.hurry = make(chan bool)

	go func() {
		for b := range a.errchan {
			panic(b)
		}
	}()
}

func (a *AIPlayer) HurryChannel() chan<- bool {
	return a.hurry
}

func (a *AIPlayer) ErrorChannel() chan<- error {
	return a.errchan
}

//Worker is the routine behind the Think; exported just in case
func (a *AIPlayer) Worker(s *game.State, whoarewe game.Color, depth int8) []float64 {
	minmax_slice := make([]float64, depth+1) // (depth < -1) causes a panic here
	if depth < 0 {                           // negative depth may be considered error in the future
		minmax_slice[0] = a.SitValue(s)
		return minmax_slice
	}
	mythoughts := make(map[int][]float64)
	index := 0 // index is for mythoughts map
	var bestsitval float64
	for state := range game.ASAOMGen(s, whoarewe) {
		mythoughts[index] = append(mythoughts[index], a.SitValue(state)) // fills in first element of mythoughts[index]
		if int(depth) > 0 {
			bestsitval = -1000000
			for mymove := range game.VFTPGen(state) {
				move_to_apply := game.Move{mymove.FromTo[0], mymove.FromTo[1], state, mymove.PawnPromotion}
				newstate, _ := move_to_apply.After()
				newthought := append([]float64{mythoughts[index][0]}, a.Worker(newstate, whoarewe, depth)...) // new slice of size (depth+1)
				if newthought[depth] > bestsitval {                                                           // we have found (so far) the best response to opponents' moves (state after 2 ops' moves)
					bestsitval = newthought[depth]
					mythoughts[index] = newthought
				}
			}
		}
		index++
	}
	bestsitval = 1000000
	for i := 0; i < index; i++ {
		if mythoughts[i][depth] < bestsitval { // we need to find the best opponents' moves to test our strategy
			minmax_slice = mythoughts[i]
		}
	}
	return minmax_slice
}

//Think is the function generating the Move
func (a *AIPlayer) Think(s *game.State, hurry <-chan bool) game.Move {
	hurryup := simple.MergeBool(hurry, a.hurry)
	for i := len(hurryup); i > 0; i-- {
		<-hurryup
	}
	thoughts := make(map[game.FromToProm][]float64) // so "bloated" for future use of hurry channel (multithreading)
	var bestmove game.FromToProm
	var bestsitval float64
	bestsitval = -1000000
	for move := range game.VFTPGen(s) {
		move_to_apply := game.Move{move.FromTo[0], move.FromTo[1], s, move.PawnPromotion}
		newstate, _ := move_to_apply.After()
		thoughts[move] = a.Worker(newstate, s.MovesNext, a.Conf.Depth)
		if thoughts[move][a.Conf.Depth] > bestsitval {
			bestmove = move
			bestsitval = thoughts[move][a.Conf.Depth]
		}
	}
	return game.Move{bestmove.FromTo[0], bestmove.FromTo[1], s, bestmove.PawnPromotion}
}

func (a *AIPlayer) HeyItsYourMove(s *game.State, hurryup <-chan bool) game.Move {
	return a.Think(s, hurryup)
}

func (a *AIPlayer) HeySituationChanges(_ *game.Move, _ *game.State) {}
func (a *AIPlayer) HeyYouLost(_ *game.State)                        {}
func (a *AIPlayer) HeyYouWon(_ *game.State)                         {}
func (a *AIPlayer) HeyYouDrew(_ *game.State)                        {}

func (a *AIPlayer) String() string {
	return fmt.Sprintf("%v%v", "SVBotDepth", a.Conf.Depth) //TODO: print whoami and conf
}

func (a *AIPlayer) AreWeWaitingForYou() bool {
	return a.waiting
}

func (a *AIPlayer) HeyWeWaitingForYou(b bool) {
	a.waiting = b
}
