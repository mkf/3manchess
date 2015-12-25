package deveng

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "fmt"

type Developer string

func (p *Developer) Initialize(*player.Gameplay) <-chan error {
	errchan := make(chan error)

	return errchan
}

func (p *Developer) HeyItsYourMove(s *game.State, hurry <-chan bool) *game.Move {

	return
}

func (p *Developer) HeySituationChanges(m *game.Move, aft *game.State) {
	return
}

func (p *Developer) HeyYouLost(*game.State) {
	fmt.Println("%s has lost", *p)
}

func (p *Developer) HeyYouWonOrDrew(*game.State) {
	fmt.Println("%s has won/drew", *p)
}

type GivingUpError string

func (g GivingUpError) Error() string {
	return string(g)
}

func (g GivingUpError) IGaveUp() string {
	return string(g)
}
