package devengfmt

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/simple"
import "fmt"
import "log"

const WhoAmI string = "3manchess-devengfmt"

type Developer struct {
	Name      string
	errchan   chan error
	ErrorChan chan<- error
	HurryChan chan<- bool
	hurry     chan bool
	gp        *player.Gameplay
	waiting   bool
}

func (p *Developer) Map() map[string]interface{} {
	return map[string]interface{}{
		"Name":   p.Name,
		"WhoAmI": WhoAmI,
	}
}

func (p *Developer) FromMap(m map[string]interface{}) {
	ok := true
	var t interface{}
	t, ok = m["Name"]
	p.Name = t.(string)
	if !ok {
		panic("Name")
	}
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

func (p *Developer) HeyItsYourMove(s *game.State, hurryi <-chan bool) game.Move {
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
	p.HeyWeWaitingForYou(false)
	return move
}

func (p *Developer) HeySituationChanges(m game.Move, aft *game.State) {
	fmt.Printf("%s, situation changed: \n", p)
	fmt.Println(m)
	fmt.Println(aft)
}

func (p *Developer) HeyYouLost(*game.State) {
	fmt.Printf("%s has lost\n", p)
}

func (p *Developer) HeyYouWon(*game.State)  { fmt.Printf("%s has won\n", p) }
func (p *Developer) HeyYouDrew(*game.State) { fmt.Printf("%s has drew\n", p) }

func (p *Developer) HeyWeWaitingForYou(b bool) {
	p.waiting = b
}

func (p *Developer) AreWeWaitingForYou() bool {
	return p.waiting
}

type GivingUpError string

func (g GivingUpError) Error() string {
	return string(g)
}

func (g GivingUpError) IGaveUp() string {
	return string(g)
}
