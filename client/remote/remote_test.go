package remote

import "testing"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/client"
import "flag"
import "os"

import "github.com/ArchieT/3manchess/ai/constsitval"
import "github.com/ArchieT/3manchess/multi"
import "time"
import "fmt"

var c *client.Client

var nso game.State
var ns *game.StateData

func init() {
	nso = game.NewState()
	ns = nso.Data()
	bu := flag.String("baseurl", os.Getenv("CHESSBASEURL"), "3manchess/multi base URL")
	flag.Parse()
	//t.Log("baseurl", bu)
	c = client.NewClient(nil, *bu)
}

func TestNew_ai3(t *testing.T) {
	var u, p int64
	var a []byte
	if l, _, err := c.SignUp(multi.SignUpPost{"remotetest", "remotetest", "remotetest"}); err == nil {
		u, p, a = l.User, l.Player, l.Auth
	} else {
		t.Log(err)
		ll, _, err := c.LogIn(multi.LoggingIn{"remotetest", "remotetest"})
		if err != nil {
			t.Fatal(err)
		}
		u, a = ll.ID, ll.AuthKey
		ppp, _, err := c.UserInfo(u)
		if err != nil {
			t.Fatal(err)
		}
		p = ppp.Player
	}
	t.Log(u, p, a)
	var mgpp multi.GameplayPost
	mgpp.State = nso
	var botsconf [3]constsitval.AIConfig
	botsconf[0].OwnedToThreatened = 4.0
	botsconf[1].OwnedToThreatened = 5.0
	botsconf[2].OwnedToThreatened = 6.0
	var botsau [3]multi.NewBotGive
	for bno, bco := range botsconf {
		var bbi, bbp int64
		var bba []byte
		if bb, _, err := c.NewBot(multi.NewBotPost{[]byte("constsitval-demotesting"), multi.Authorization{u, a}, fmt.Sprint("bot", bno), bco.Byte()}); err == nil {
			bbi, bbp, bba = bb.Botid, bb.PlayerID, bb.AuthKey
		} else {
			t.Log(err)
			bbi = int64(bno + 1) //yeah
			bbb, _, err := c.BotKey(multi.BotKeyGetting{bbi, multi.Authorization{u, a}})
			if err != nil {
				t.Fatal(err)
			}
			binfo, _, err := c.BotInfo(bbi)
			t.Log(binfo)
			if err != nil {
				t.Fatal(err)
			}
			bbp, bba = bbb.ID, bbb.AuthKey
		}
		botsau[bno] = multi.NewBotGive{bbi, bbp, bba}
	}
	mgpp.White = &botsau[0].PlayerID
	mgpp.Gray = &botsau[1].PlayerID
	mgpp.Black = &botsau[2].PlayerID
	gpg, _, err := c.AddGame(mgpp)
	t.Log(*gpg)
	if err != nil {
		t.Fatal(err)
	}
	echn := make(chan error)
	endchn := make(chan bool)
	go func() {
		for u := range echn {
			t.Log(u)
		}
	}()
	var nboty [3]*G
	for bno := range game.COLORS {
		var aii constsitval.AIPlayer
		aii.Conf = botsconf[bno]
		yg, err := New(
			c,
			&aii,
			game.COLORS[bno],
			gpg.Key,
			multi.Authorization{
				botsau[bno].PlayerID,
				botsau[bno].AuthKey,
			},
			func(g *G) (int64, error) {
				t.Log("AFTFUNCC")
				for {
					t.Log("st for aftfunc", *g.state)
					a, _, err := g.C().After(
						g.gameid,
						[3]*int64{nil, nil, nil},
					)
					t.Log("aftfunc kitchen:", a, err)
					if len(*a) > 0 {
						return (*a)[0].Key, err
					}
					if err != nil {
						return -1, err
					}
					time.Sleep(3 * time.Second)
				}
				return -1, err
			},
			endchn,
			echn,
		)
		if err != nil {
			t.Log(err)
		}
		nboty[bno] = yg
	}
	t.Log("end", <-endchn)
}
