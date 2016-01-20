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
	color game.Color
	Name  string
}

func NewReality() *Reality {
	return new(Reality)
}

func (re *Reality) MakePlayers(who ...game.Color) map[game.Color]*RealPlayer {
	ourm := make(map[game.Color]*RealPlayer)
	for _, c := range who {
		ourm[c] = RealPlayer{re, c, c.String()}
	}
	return ourm
}
