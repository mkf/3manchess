package sitvalues

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/simple"
import "github.com/ArchieT/3manchess/player"
import "sync"
import "sync/atomic"
import "fmt"

import "log"

const DEFFIXPREC float64 = 0.0002

type AIPlayer struct {
	errchan           chan error
	ErrorChan         chan<- error
	hurry             chan bool
	HurryChan         chan<- bool
	FixedPrecision    float64
	curfixprec        float64
	OwnedToThreatened float64
	gp                *player.Gameplay
	waiting           bool
}

func (a *AIPlayer) Initialize(gp *player.Gameplay) {
	if a.FixedPrecision == 0.0 {
		a.FixedPrecision = DEFFIXPREC
	}
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

func (a *AIPlayer) Worker(chance float64, give chan<- float64, state *game.State, whoarewe game.Color) {
	log.Println(chance)
	state.EvalDeath()
	if !(state.PlayersAlive.Give(whoarewe)) {
		give <- a.SitValue(state) * chance
		return
	}
	if chance < a.curfixprec {
		give <- a.SitValue(state) * chance
		return
	}
	var wg sync.WaitGroup
	var i, j, k, l int8
	possib := make(chan *game.State, 2050)
	var ourft game.FromTo

	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					wg.Add(1)
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						sv := ourft.Move(state)
						if v, err := sv.After(); err == nil {
							possib <- v
						}
						wg.Done()
					}(i, j, k, l)
				}
			}
		}
	}
	wg.Wait()
	var newchance float64
	newchance = chance / float64(len(possib))
	for m := range possib {
		wg.Add(1)
		go func(mgiving *game.State) {
			a.Worker(newchance, give, mgiving, whoarewe)
			wg.Done()
		}(m)
	}
	wg.Wait()
}

func (a *AIPlayer) Think(s *game.State, hurry <-chan bool) *game.Move {
	a.curfixprec = a.FixedPrecision
	hurryup := simple.MergeBool(hurry, a.hurry)
	log.Println("Started thinking")
	for i := len(hurryup); i > 0; i-- {
		<-hurryup
	}
	log.Println("ALONG")
	thoughts := make(map[game.FromTo]*float64)
	var i, j, k, l int8
	var ourft game.FromTo
	countem := new(uint32)
	atomic.StoreUint32(countem, 0)
	var wg, gwg sync.WaitGroup
	wg.Add(1)
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						sv := ourft.Move(s)
						if v, err := sv.After(); err == nil {
							gwg.Add(1)
							go func(n game.FromTo) {
								log.Println("Yeah")
								atomic.AddUint32(countem, 1)
								var newchance float64
								wg.Wait()
								log.Println("GotChance")
								newchance = 1.0 / float64(*countem)
								ourchan := make(chan float64, 100)
								makefloat := new(float64)
								thoughts[n] = makefloat
								go func(ch <-chan float64, ou *float64) {
									*ou += <-ch
								}(ourchan, makefloat)
								a.Worker(newchance, ourchan, v, s.MovesNext)
								gwg.Done()
							}(ourft)
						}
					}(i, j, k, l)
				}
			}
		}
	}
	wg.Done()
	go func() {
		for {
			<-hurryup
			a.curfixprec *= 2
		}
	}()
	gwg.Wait()
	a.HeyWeWaitingForYou(false)
	var max float64
	for i := range thoughts {
		if *(thoughts[i]) > max {
			max = *(thoughts[i])
		}
	}
	ourfts := make([]game.FromTo, 0, 10)
	for i := range thoughts {
		if *(thoughts[i]) == max {
			ourfts = append(ourfts, i)
		}
	}
	if len(ourfts) == 0 {
		panic("len(ourfts)==0 !!!!")
	}
	ormov := ourfts[9].Move(s)
	return &ormov
}

func (a *AIPlayer) HeyItsYourMove(s *game.State, hurryup <-chan bool) *game.Move {
	return a.Think(s, hurryup)
}

func (a *AIPlayer) HeySituationChanges(_ *game.Move, _ *game.State) {}
func (a *AIPlayer) HeyYouLost(_ *game.State)                        {}
func (a *AIPlayer) HeyYouWon(_ *game.State)                         {}
func (a *AIPlayer) HeyYouDrew(_ *game.State)                        {}

func (a *AIPlayer) String() string {
	return fmt.Sprintf("%s%e", "SVBotPrec", a.FixedPrecision)
}

func (a *AIPlayer) AreWeWaitingForYou() bool {
	return a.waiting
}

func (a *AIPlayer) HeyWeWaitingForYou(b bool) {
	a.waiting = b
}
