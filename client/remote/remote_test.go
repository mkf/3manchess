package remote

import "testing"
import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/client"
import "flag"
import "os"
import "github.com/coreos/pkg/flagutil"
import "github.com/ArchieT/3manchess/ai/constsitval"
import "github.com/ArchieT/3manchess/multi"
import "time"

var c *client.Client

var ns *game.StateData

func init() {
	ns = game.NewState().Data()
	flags := flag.NewFlagSet("remotetest", flag.ExitOnError)
	bu := flags.String("baseurl", "http://platinum.edu.pl:8082/", "3manchess/multi base URL")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "REMOTECHESSTEST")
	log.Println("baseurl", bu)
	c = client.NewClient(nil, bu)
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
	var mgpp multi.GameplayPost
	mgpp.Date = time.Now()
	mgpp.State = *ns
	gpg, _, err := c.AddGame(mgpp)
	if err != nil {
		t.Fatal(err)
	}
	//	b1,_,err:=
}
