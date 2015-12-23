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

func Worker(thought *float64, chance float64, mutex *sync.RWMutex, state game.State) {

}

func (a *AIsettings) Think(s *game.State) game.Move {
	var thinking [6][24][6][24]float64
	var mutex sync.RWMutex
	var i, j, k, l int8
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					go func(i, j, k, l int8) {
						if !(s.AnyPiece(game.Pos{i, j}, game.Pos{k, l})) {
							mutex.Lock()
							thinking[i][j][k][l] = -1
							mutex.Unlock()
						}
					}(i, j, k, l)
				}
			}
		}
	}

}
