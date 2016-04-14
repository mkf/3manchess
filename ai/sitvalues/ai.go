package sitvalues

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/simple"
import "github.com/ArchieT/3manchess/player"
import "sync"
import "sync/atomic"
import "fmt"
import "github.com/ArchieT/3manchess/ai"
import "encoding/json"

const DEFFIXPREC float64 = 0.0002

const DEFPAWNPROMOTION = game.Queen

const DEFOWN2THRTHD = 4.0

const WhoAmI string = "3manchess-ai_sitvalues"

type AIPlayer struct {
	Name       string
	errchan    chan error
	hurry      chan bool
	Conf       AIConfig
	curfixprec float64
	gp         *player.Gameplay
	waiting    bool
}

func (a *AIPlayer) Config() ai.Config {
	return a.Conf
}

func (a *AIPlayer) SetConf(b []byte) error {
	ac := new(AIConfig)
	e := json.Unmarshal(b, ac)
	if e == nil {
		a.Conf = *ac
	}
	return e
}

type AIConfig struct {
	Precision         float64
	OwnedToThreatened float64
	PawnPromotion     game.FigType //it will be possible to set it to 0 for automatic choice (not yet implemented)
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
	if a.Conf.Precision == 0.0 {
		a.Conf.Precision = DEFFIXPREC
	}
	if a.Conf.PawnPromotion == game.ZeroFigType {
		a.Conf.PawnPromotion = DEFPAWNPROMOTION
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

//Worker is the routine behind Think; exported just in case
func (a *AIPlayer) Worker(chance float64, give chan<- float64, state *game.State, whoarewe game.Color) {
	state.EvalDeath()
	if !(state.PlayersAlive.Give(whoarewe)) { //if we are dead
		give <- a.SitValue(state, whoarewe) * chance
		return
	}
	if chance < a.curfixprec { //if we are too deep
		give <- a.SitValue(state, whoarewe) * chance
		return
	}
	var wg sync.WaitGroup
	possib := make(chan *game.State, 2050)
	for ofrom, loto := range game.AMFT {
		for _, oto := range loto {
			wg.Add(1)
			go func(ourft game.FromTo) {
				sv := ourft.Move(state)
				sv.PawnPromotion = a.Conf.PawnPromotion
				if v, err := sv.After(); err == nil {
					possib <- v
				}
				wg.Done()
			}(game.FromTo{ofrom, oto})
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

//Think is the function generating the Move; atm it does not return anything, but will return game.Move
func (a *AIPlayer) Think(s *game.State, hurry <-chan bool) game.Move {
	a.curfixprec = a.Conf.Precision
	hurryup := simple.MergeBool(hurry, a.hurry)
	for i := len(hurryup); i > 0; i-- {
		<-hurryup
	}
	thoughts := make(map[game.FromTo]*float64)
	countem := new(uint32)
	atomic.StoreUint32(countem, 0)
	var wg, gwg sync.WaitGroup
	var tmx sync.Mutex
	wg.Add(1)
	for ofrom, loto := range game.AMFT {
		for _, oto := range loto {
			go func(ourft game.FromTo) {
				sv := ourft.Move(s)
				sv.PawnPromotion = a.Conf.PawnPromotion
				if v, err := sv.After(); err == nil {
					gwg.Add(1)
					go func(n game.FromTo) {
						atomic.AddUint32(countem, 1)
						var newchance float64
						wg.Wait()
						newchance = 1.0 / float64(*countem)
						ourchan := make(chan float64, 100)
						makefloat := new(float64)
						tmx.Lock()
						thoughts[n] = makefloat
						tmx.Unlock()
						go func(ch <-chan float64, ou *float64) {
							*ou += <-ch
						}(ourchan, makefloat)
						a.Worker(newchance, ourchan, v, s.MovesNext)
						gwg.Done()
					}(ourft)
				}
			}(game.FromTo{ofrom, oto})
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
	ormov.PawnPromotion = a.Conf.PawnPromotion
	return ormov
}

func (a *AIPlayer) HeyItsYourMove(s *game.State, hurryup <-chan bool) game.Move {
	return a.Think(s, hurryup)
}

func (a *AIPlayer) HeySituationChanges(_ game.Move, _ *game.State) {}
func (a *AIPlayer) HeyYouLost(_ *game.State)                       {}
func (a *AIPlayer) HeyYouWon(_ *game.State)                        {}
func (a *AIPlayer) HeyYouDrew(_ *game.State)                       {}

func (a *AIPlayer) String() string {
	return fmt.Sprintf("%s%e", "SVBotPrec", a.Conf.Precision) //TODO: print whoami and conf
}

func (a *AIPlayer) AreWeWaitingForYou() bool {
	return a.waiting
}

func (a *AIPlayer) HeyWeWaitingForYou(b bool) {
	a.waiting = b
}
