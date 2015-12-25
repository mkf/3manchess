package human

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	Initialize() <-chan error
	HeyItsYourMove(*game.Move, <-chan bool) *game.Move //that channel is for signalling to hurry up
	HeySituationChanges(*game.Move, *game.State)
	HeyYouLost(*game.State)
	HeyYouWonOrDrew(*game.State)
}

//Gameplay is a list of players and the current gamestate pointer
type Gameplay struct {
	White *Player
	Gray  *Player
	Black *Player
	*game.State
}

//NewGame returns a new Gameplay
func NewGame(w *Player, g *Player, b *Player) Gameplay {
	ns := game.NewState()
	return Gameplay{w, g, b, &ns}
}
