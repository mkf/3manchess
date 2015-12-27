package ai

import "github.com/ArchieT/3manchess/game"

//import "github.com/ArchieT/3manchess/player"
import "sync"
import "sync/atomic"

//AIPlayer is a struct of AI settings; Think is it's method.
type AIPlayer struct {
	errchan        chan error
	ErrorChan      chan<- error
	hurry          chan bool
	HurryChan      chan<- bool
	FixedPrecision float64
}

func (a *AIPlayer) HurryChannel() chan<- bool {
	return a.HurryChan
}

func (a *AIPlayer) ErrorChannel() chan<- error {
	go func() {
		for b := range a.errchan {
			panic(b)
		}
	}()
	return a.ErrorChan
}

//func Worker(thinking *[6][24][6][24]int64, mutex *sync.RWMutex, state game.State) {
//}

//Worker is the routine behind the Think; exported just in case
func (a *AIPlayer) Worker(chance float64, give chan<- float64, state *game.State, whoarewe game.Color) {
	if state.EvalDeath(); !(state.PlayersAlive.Give(whoarewe)) {
		give <- -chance
		return
	}
	if chance < a.FixedPrecision {
		return
	}
	var wg1 sync.WaitGroup //whole thread waitgroup
	var i, j, k, l int8    //pos int8s
	possib := make(chan *game.State, 2050)
	var ourft game.FromTo

	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					wg1.Add(1)
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						sv := ourft.Move(state)
						if v, err := sv.After(); err == nil {
							possib <- v
						}
						wg1.Done()
					}(i, j, k, l)
				}
			}
		}
	}
	wg1.Wait()
	var newchance float64
	newchance = chance / float64(len(possib))
	for m := range possib {
		go a.Worker(newchance, give, m, whoarewe)
	}
}

//Think is the function generating the Move; atm it does not return anything, but will return game.Move
func (a *AIPlayer) Think(s *game.State, hurry <-chan bool) *game.Move {
	//var thinking [6][24][6][24]float64
	hurryup := merge(hurry, a.hurry)
	if len(hurryup) > 0 {
		for _ := range hurryup {
		}
	}
	thoughts := make(map[game.FromTo]*float64)
	var i, j, k, l int8
	var ourft game.FromTo
	countem := new(uint32)
	atomic.StoreUint32(countem, 0)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						sv := ourft.Move(s)
						if v, err := sv.After(); err == nil {
							go func(n game.FromTo) {
								atomic.AddUint32(countem, 1)
								var newchance float64
								wg1.Wait()
								newchance = 1.0 / float64(*countem)
								ourchan := make(chan float64, 100)
								makefloat := new(float64)
								thoughts[n] = makefloat
								a.Worker(newchance, ourchan, v, s.MovesNext)
								go func(ch <-chan float64, ou *float64) {
									*ou += <-ch
								}(ourchan, makefloat)
							}(ourft)
						}
					}(i, j, k, l)
				}
			}
		}
	}
	wg1.Done()
	_ = <-hurryup
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
		panic("len(ourfts)==0 !!!")
	}
	ormov := ourfts[0].Move(s)
	return &ormov
}

//HeyItsYourMove works as specified in github.com/ArchieT/3manchess/player; here it just calls Think
func (a *AIPlayer) HeyItsYourMove(m *game.Move, s *game.State, hurryup <-chan bool) *game.Move {
	return a.Think(s, hurryup)
}

func (a *AIPlayer) HeySituationChanges(_ *game.Move, _ *game.State) {}

func (a *AIPlayer) HeyYouLost(_ *game.State) {}

func (a *AIPlayer) HeyYouWonOrDrew(_ *game.State) {}

func (a *AIPlayer) String() string {
	return a.FixedPrecision.String()
}
