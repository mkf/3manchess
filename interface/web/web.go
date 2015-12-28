package web

import (
	"encoding/json"
	"github.com/ArchieT/3manchess/game"
	"github.com/ArchieT/3manchess/player"
	"github.com/ArchieT/3manchess/simple"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Entity struct {
	*game.State
	AwaitingMove bool
	Hurry        bool
	Happened     []*game.Move
}

//type ErrorList []error

//func (el ErrorList) Error() string {
//	return fmt.Sprint(el)
//}

func (p *WebPlayer) HandlerGen() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		var ost *game.State
		for {
			var ent Entity
			ent.Happened = make([]*game.Move, 0, 3)
			for {
				select {
				case i := <-p.happened:
					ent.Happened = append(ent.Happened, i)
				case ost = <-p.awaitingmove:
					ent.AwaitingMove = true
				default:
					break
				}
			}
			ent.State = p.gp.State

			msg, err := json.Marshal(ent)

			if err != nil {
				p.ErrorChan <- err
				panic(err)
			}

			if err = websocket.Message.Send(ws, msg); err != nil {
				p.ErrorChan <- err
				for i := range ent.Happened {
					p.happened <- ent.Happened[i]
				}
				if ent.AwaitingMove {
					p.awaitingmove <- ost
				}
			}

			for ent.AwaitingMove {
				var reply string

				if err = websocket.Message.Receive(ws, &reply); err != nil {
					p.ErrorChan <- err
					continue
				}
				log.Println("Received back from client: " + reply)
				var fromto game.FromTo
				if err := json.Unmarshal([]byte(reply), &fromto); err != nil {
					p.ErrorChan <- err
					continue
				}
				p.givemethefromto <- fromto

			}

		}
	})
}

type WebPlayer struct {
	Name            string
	port            string
	errchan         chan error
	ErrorChan       chan<- error
	hurry           chan bool
	HurryChan       chan<- bool
	awaitingmove    chan *game.State
	givemethefromto chan game.FromTo
	happened        chan *game.Move
	gp              *player.Gameplay
	wsl             []*websocket.Conn
	waiting         bool
}

func (p *WebPlayer) Initialize(gp *player.Gameplay) {
	errchan := make(chan error)
	p.errchan = errchan
	p.happened = make(chan *game.Move, 100)
	http.Handle("/"+p.Name, p.HandlerGen())
	if err := http.ListenAndServe(p.port, nil); err != nil {
		errchan <- err
	}
	for len(p.wsl) == 0 {
	}
}

func (p *WebPlayer) String() string {
	return p.Name
}

func (p *WebPlayer) ErrorChannel() chan<- error {
	return p.ErrorChan
}

func (p *WebPlayer) HurryChannel() chan<- bool {
	return p.HurryChan
}

func (p *WebPlayer) HeyItsYourMove(s *game.State, hurryi <-chan bool) *game.Move {
	hurry := simple.MergeBool(hurryi, p.hurry)
	go func() {
		for {
			<-hurry
		}
	}()
	for i := range p.wsl {
		go func(i int) {
			var err error
			for {
				var reply string

				if err = websocket.Message.Receive(p.wsl[i], &reply); err != nil {
					//p.ErrorChan <- err
					continue
				}
				log.Println("Received back from client: " + reply)
				var fromto game.FromTo
				if err := json.Unmarshal([]byte(reply), &fromto); err != nil {
					p.ErrorChan <- err
					continue
				}
				p.givemethefromto <- fromto
			}
		}(i)
	}
	ourfromto := <-p.givemethefromto
	ourmove := ourfromto.Move(s)
	p.HeyWeAreWaitingForYou(false)
	return &ourmove
}

func (p *WebPlayer) HeySituationChanges(m *game.Move, aft *game.State) { p.happened <- m }

func (p *WebPlayer) HeyYouLost(*game.State) {}
func (p *WebPlayer) HeyYouWon(*game.State)  {}
func (p *WebPlayer) HeyYouDrew(*game.State) {}

func (p *WebPlayer) AreWeWaitingForYou() bool     { return p.waiting }
func (p *WebPlayer) HeyWeAreWaitingForYou(b bool) { p.waiting = b }

type GivingUpError string

func (g GivingUpError) Error() string   { return string(g) }
func (g GivingUpError) IGaveUp() string { return string(g) }
