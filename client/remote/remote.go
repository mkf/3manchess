package remote

import "github.com/ArchieT/3manchess/client"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/multi"
import "log"

type G struct {
	c       *client.Client
	gameid  int64
	state   *game.State
	our     player.Player
	color   game.Color
	auth    multi.Authorization
	errchan chan error
}

func New(c *client.Client, our player.Player, color game.Color, gameid int64, auth multi.Authorization) (*G, error) {
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
		c:      c,
		gameid: gameid,
		state:  new(game.State),
		our:    our,
		color:  color,
		auth:   auth,
	}
	g.state.FromData(sd)
	//go procedure
	return &g, nil
}

func (g *G) Turn() (breaking bool, err error) {
	g.our.HeyWeWaitingForYou(true)
	hurry := make(chan bool)
	bef := g.state
	move := g.our.HeyItsYourMove(bef, hurry)
	ftp := game.FromToProm{FromTo: game.FromTo{move.From, move.To}, PawnPromotion: move.PawnPromotion}
	tp := multi.TurnPost{ftp, g.auth}
	maak, _, err := g.c.Service.Turn(g.gameid, tp)
	if err != nil {
		return
	}
	g.gameid = maak.AfterGameKey
	gd, _, err := g.c.Service.Play(g.gameid)
	if err != nil {
		return
	}
	sd, _, err := g.c.Service.State(gd.State)
	if err != nil {
		return
	}
	g.state = new(game.State)
	g.state.FromData(sd)
	breaking = g.state.PlayersAlive.Give(g.color)
	md, _, err := g.c.Service.Move(maak.MoveKey)
	if err != nil {
		return
	}
	mov := game.Move{
		From:          game.Pos{md.FromTo[0], md.FromTo[1]},
		To:            game.Pos{md.FromTo[2], md.FromTo[3]},
		Before:        bef,
		PawnPromotion: game.FigType(md.PawnPromotion),
	}
	g.our.HeySituationChanges(&mov, g.state)
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
