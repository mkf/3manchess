package reality

import "github.com/ArchieT/3manchess/interface/reality/machine"
import "github.com/ArchieT/3manchess/interface/reality/camget"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "log"

type Reality struct {
	camget.View
	machine.Machine
	BlackIsOnWhitesRight bool
}

//func
