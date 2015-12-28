package devengchan

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/simple"
import "fmt"
import "log"

type Developer struct {
	Name      string
	errchan   chan error
	ErrorChan chan<- error
	HurryChan chan<- bool
	hurry     chan bool
	gp        *player.Gameplay
}

func (p *Developer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	p.errchan = errchan
	hurry := make(chan bool)
	p.hurry = hurry
	fmt.Printf("%s initialized with Gameplay:\n", p)
	fmt.Println(gp)
	fmt.Println("")
	p.gp = gp
	p.ErrorChan = errchan
	p.HurryChan = hurry
	go p.logger()
}

func (p *Developer) logger() {
	var err error
	for {
		err = <-p.errchan
		log.Println(err)
	}
}

func (p *Developer) String() string {
	return p.Name
}

func (p *Developer) ErrorChannel() chan<- error {
	return p.ErrorChan
}

func (p *Developer) HurryChannel() chan<- bool {
	return p.HurryChan
}

func (p *Developer) HeyItsYourMove(s *game.State, hurryi <-chan bool) *game.Move {
	hurry := simple.MergeBool(hurryi, p.hurry)
	go func() {
		for {
			<-hurry
			fmt.Print("@")
		}
	}()
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
