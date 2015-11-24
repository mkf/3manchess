package ai

import "github.com/ArchieT/3manchess/game"
import "time"
import "sync"

type AIsettings struct {
	Time time.Duration
	//thinking [6][24][6][24]int64
	//mutex RWMutex
}

func NewAI(t time.Duration) AIsettings {
	var our AIsettings
	our.Time = t
	return our
}

//func Worker(thinking *[6][24][6][24]int64, mutex *sync.RWMutex, state game.State) {
//}

func Worker(thought *float64, chance float64, mutex *sync.RWMutex, state game.State, whoarewe game.Color) {
	var wg1 sync.WaitGroup
	var i, j, k, l int8
	possib := make([]*game.State, 0, 30)
	var possibmutex sync.RWMutex
	var ourft game.FromTo

	if state.EvalDeath(); !(state.PlayersAlive.Give(whoarewe)) {
		mutex.Lock()
		*thought -= chance
		mutex.Unlock()
		return
	}

	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					wg1.Add(1)
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						if v, err := ourft.Move(&state).After(); err == nil {
							possibmutex.Lock()
							possib = append(possib, v)
							possibmutex.Unlock()
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
	for _, m := range possib {
		go Worker(thought, newchance, mutex, *m, whoarewe)
	}
}

func (a *AIsettings) Think(s *game.State) { //game.Move {
	//var thinking [6][24][6][24]float64
	thinking := make(map[game.FromTo]*float64)
	states := make(map[game.FromTo]*game.State)
	var mutex, possibmutex, statesmutex sync.RWMutex
	var i, j, k, l int8
	var wg1 sync.WaitGroup
	var ourft game.FromTo
	var possib uint32
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					wg1.Add(1)
					go func(i, j, k, l int8) {
						ourft = game.FromTo{game.Pos{i, j}, game.Pos{k, l}}
						if v, err := ourft.Move(s).After(); err == nil {
							possibmutex.Lock()
							possib++
							possibmutex.Unlock()
							mutex.Lock()
							zerofloat := float64(0)
							thinking[ourft] = &zerofloat
							mutex.Unlock()
							statesmutex.Lock()
							states[ourft] = v
							statesmutex.Unlock()
						}
						wg1.Done()
					}(i, j, k, l)
				}
			}
		}
	}
	wg1.Wait()
	for n := range thinking {
		go func(n game.FromTo) {
			var newchance float64
			newchance = 1.0 / float64(possib)
			Worker((thinking[n]), newchance, mutex, *(states[n]), s.MovesNext)
		}(n)
	}
}
