package human

import "github.com/ArchieT/3manchess/game"

//Player is either AI or a human via some UI
type Player interface {
	HeyItsYourMove(*game.Move) *game.Move
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
	HeyItsYourMove(*game.Move) *game.Move
	HeySituationChanges(*game.Move)
	HeyYouLost(*game.State)
	HeyYouWonOrDrew(*game.State)
	AreYouGivinUp(*game.State) bool
}

func NewGame(w *Player, g *Player, b *Player) {
	return Gameplay{w, g, b, game.NEWGAME}
}
