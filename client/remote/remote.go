package remote

import "github.com/ArchieT/3manchess/client"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/multi"
import "github.com/ArchieT/3manchess/server"
import "fmt"

//G is a remote gameplay, with one local player accessing it
type G struct {
	c       *client.Client
	gameid  int64
	state   *game.State
	our     player.Player
	color   game.Color
	auth    multi.Authorization
	errchan chan error
	af      AfterFunc
}

func (g *G) C() *client.Client {
	return g.c
}

func (g *G) GState() *game.State {
	return g.state
}

func (g *G) GGameID() int64 {
	return g.gameid
}

//New returns a new remote gameplay (G)
func New(c *client.Client, our player.Player, color game.Color, gameid int64, auth multi.Authorization, af AfterFunc,
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
		c:       c,
		gameid:  gameid,
		state:   sd,
		our:     our,
		color:   color,
		auth:    auth,
		errchan: make(chan error),
		af:      af,
	}
	go func(ec chan<- error) {
		for i := range g.errchan {
			ec <- i
		}
	}(errch)
	go g.Procedure(end, errch)
	return &g, nil
}

func (g *G) askplayer() (oftp game.FromToProm, err error) {
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
		return
	}
	g.gameid = maak.AfterGameKey
	md, _, err := g.c.Service.Move(maak.MoveKey)
	if err != nil {
		return
	}
	err = g.getstate()
	oftp = md.FromToProm()
	return
}

type ErrorNotAfter struct {
	server.MoveData
	MoveKey int64
	Before  int64
}

func (e *ErrorNotAfter) Error() string {
	if e != nil {
		return fmt.Sprint("ErrorNotAfter", *e)
	}
	return ""
}

func (g *G) movechk() (oftp game.FromToProm, err error) {
	var i int64
	i, err = g.af(g)
	if err == nil {
		var md *server.MoveData
		md, _, err = g.c.Service.Move(i)
		if err != nil {
			return
		}
		if md.BeforeGame == g.gameid {
			g.gameid = md.AfterGame
			oftp = md.FromToProm()
			err = g.getstate()
		} else {
			e := ErrorNotAfter{*md, i, g.gameid}
			err = &e
		}
	}
	return
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
	g.state = sd
	return err
}

type AfterFunc func(g *G) (moveid int64, err error)

func (g *G) Turn() (breaking bool, err error) {
	bef := g.state
	var oftp game.FromToProm
	var do func() (game.FromToProm, error)
	if g.state.MovesNext == g.color {
		do = g.askplayer
	} else {
		do = g.movechk
	}
	if oftp, err = do(); err != nil {
		return
	}
	breaking = !g.state.PlayersAlive.Give(g.color)
	g.our.HeySituationChanges(oftp.Move(bef), g.state)
	return
}

func (g *G) Procedure(end chan<- bool, errch chan<- error) {
	if g.state.PlayersAlive.Give(g.color) {
		for {
			br, err := g.Turn()
			if br {
				break
			}
			if err != nil {
				errch <- err
				continue
			}
		}
	}
	end <- false
}
