package deveng

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "fmt"

type Developer struct {
	Name      string
	errchan   chan error
	ErrorChan chan<- error
	gp        *player.Gameplay
}

func (p *Developer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	p.errchan = errchan
	fmt.Printf("%s initialized with Gameplay:\n", p)
	fmt.Println(gp)
	fmt.Println("")
	p.gp = gp
	p.ErrorChan = errchan
}

func (p *Developer) String() string {
	return p.Name
}

func (p *Developer) ErrorChannel() <-chan error {
	return p.ErrorChan
}

func (p *Developer) HeyItsYourMove(s *game.State, hurry <-chan bool) *game.Move {
	fmt.Printf("%s, it's your move\n", p)
	fmt.Println(s)
	fmt.Println("from_rank from_file to_rank to_file")
	fmt.Printf("%s:", p)
	var fr, ff, tr, tf int8
	_, err := fmt.Scanf("%d %d %d %d", &fr, &ff, &tr, &tf)
	if err != nil {
		p.errchan <- err
	}
	fromto := game.FromTo{game.Pos{fr, ff}, game.Pos{tr, tf}}
	move := fromto.Move(s)
	return &move
}

func (p *Developer) HeySituationChanges(m *game.Move, aft *game.State) {
	fmt.Printf("%s, situation changed: \n", p)
	fmt.Println(m)
	fmt.Println(aft)
}

func (p *Developer) HeyYouLost(*game.State) {
	fmt.Printf("%s has lost\n", p)
}

func (p *Developer) HeyYouWonOrDrew(*game.State) {
	fmt.Printf("%s has won/drew\n", p)
}

type GivingUpError string

func (g GivingUpError) Error() string {
	return string(g)
}

func (g GivingUpError) IGaveUp() string {
	return string(g)
}
