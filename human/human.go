package human

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	HeyItsYourMove(*game.Move, <-chan bool) *game.Move //that channel is for signalling to hurry up
}

//Gameplay is a list of players and the current gamestate pointer
type Gameplay struct {
	White *Player
	Gray  *Player
	Black *Player
	*game.State
}

//Human implements Player and some more, UI-oriented options
type Human interface {
	HeyItsYourMove(*game.Move, <-chan bool) *game.Move
	HeySituationChanges(*game.Move)
	HeyYouLost(*game.State)
	HeyYouWonOrDrew(*game.State)
	AreYouGivinUp(*game.State) bool
}

//NewGame returns a new Gameplay
func NewGame(w *Player, g *Player, b *Player) Gameplay {
	ns := game.NewState()
	return Gameplay{w, g, b, &ns}
}
