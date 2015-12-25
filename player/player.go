package human

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	Initialize(*Gameplay) <-chan error
	HeyItsYourMove(*game.State, <-chan bool) *game.Move //that channel is for signalling to hurry up
	HeySituationChanges(*game.Move, *game.State)
	HeyYouLost(*game.State)
	HeyYouWonOrDrew(*game.State)
	String() string
}

//Gameplay is a list of players and the current gamestate pointer
type Gameplay struct {
	White *Player
	Gray  *Player
	Black *Player
	*game.State
}

type GivingUpError interface {
	error
	IGaveUp() string
}

//NewGame returns a new Gameplay
func NewGame(w *Player, g *Player, b *Player) *Gameplay {
	ns := game.NewState()
	gp := Gameplay{w, g, b, &ns}
	return &gp
}
