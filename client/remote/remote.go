package remote

import "github.com/ArchieT/3manchess/client"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/multi"
import "log"

//G is a remote gameplay, with one local player accessing it
type G struct {
	c         *client.Client
	gameid    int64
	state     *game.State
	our       player.Player
	poll      chan bool
	aftermove chan int64
	ourftpaft chan game.FromToProm
	color     game.Color
	auth      multi.Authorization
	errchan   chan error
}

//New returns a new remote gameplay (G)
func New(c *client.Client, our player.Player, color game.Color, gameid int64, auth multi.Authorization,
	end chan<- bool, errch chan<- error) (*G, error) {
	gd, _, err := c.Service.Play(gameid)
	if err != nil {
		return nil, err
	}
	sd, _, err := c.Service.State(gd.State)
	if err != nil {
		return nil, err
	}
	our.Initialize(nil) //TODO: make player.Player play blind, w/o info about other players (or at least allow nil)
	g := G{
		c:         c,
		gameid:    gameid,
		state:     new(game.State),
		our:       our,
		color:     color,
		auth:      auth,
		poll:      make(chan bool),
		aftermove: make(chan int64),
		ourftpaft: make(chan game.FromToProm),
		errchan:   make(chan error),
	}
	go func(ec chan<- error) {
		for i := range g.errchan {
			ec <- i
		}
	}(errch)
	g.state.FromData(sd)
	go g.Procedure(end, errch)
	return &g, nil
}

func (g *G) Poll() {
	g.poll <- true
}

func (g *G) ForceAfter(moveid int64) {
	g.poll <- false
	g.aftermove <- moveid
}

func (g *G) askplayer() error {
	g.our.HeyWeWaitingForYou(true)
	hurry := make(chan bool)
	move := g.our.HeyItsYourMove(g.state, hurry)
	maak, _, err := g.c.Service.Turn(
		g.gameid,
		multi.TurnPost{
			FromToProm: game.FromToProm{
				FromTo:        game.FromTo{move.From, move.To},
				PawnPromotion: move.PawnPromotion,
			},
			WhoPlayer: g.auth,
		},
	)
	if err != nil {
		return err
	}
	g.gameid = maak.AfterGameKey
	g.getstate()
	md, _, err := g.c.Service.Move(maak.MoveKey)
	if err != nil {
		return err
	}
	g.ourftpaft <- md.FromToProm()
	return err
}

func (g *G) polling() {
	for i := range g.poll {
		if !i {
			break
		}
		var gotit bool
		if gotit {
			g.aftermove <- 1234
		}
	}
	g.poll = make(chan bool)
}

func (g *G) movechk() {
	for i := range g.aftermove {
		md, _, err := g.c.Service.Move(i)
		if err != nil {
			g.errchan <- err
			continue
		}
		if md.BeforeGame == g.gameid {
			g.gameid = md.AfterGame
			g.getstate()
			g.ourftpaft <- md.FromToProm()
		}
	}
}

func (g *G) getstate() error {
	gd, _, err := g.c.Service.Play(g.gameid)
	if err != nil {
		return err
	}
	sd, _, err := g.c.Service.State(gd.State)
	if err != nil {
		return err
	}
	s := new(game.State)
	s.FromData(sd)
	g.state = s
	return err
}

func (g *G) Turn() (breaking bool, err error) {
	bef := g.state
	if g.state.MovesNext == g.color {
		if err = g.askplayer(); err != nil {
			return
		}
	} else {
		go g.polling()
	}
	oftp := <-g.ourftpaft
	breaking = g.state.PlayersAlive.Give(g.color)
	g.our.HeySituationChanges(oftp.Move(bef), g.state)
	return
}

func (g *G) Procedure(end chan<- bool, errch chan<- error) {
	log.Println("Remote::Procedure")
	if g.state.PlayersAlive.Give(g.color) {
		log.Println("Remote::Given")
		for {
			ok, err := g.Turn()
			if !ok {
				break
			}
			if err != nil {
				errch <- err
				continue
			}
			log.Println("Remote::Turning...")
		}
	}
	log.Println("Remote::NotTurningAnymore")
	end <- false
}
