package main

import "github.com/ArchieT/3manchess/client"
import "github.com/ArchieT/3manchess/multi"
import "github.com/ArchieT/3manchess/client/remote"
import "github.com/ArchieT/3manchess/ai/constsitval"
import "github.com/ArchieT/3manchess/ai"
import "github.com/ArchieT/3manchess/game"
import "flag"
import "os"
import "log"
import "time"

func main() {
	bu := flag.String("baseurl", os.Getenv("CHESSBASEURL"), "3manchess/multi base URL")
	login := flag.String("login", "remote", "login")
	passwd := flag.String("passwd", "remote", "passwd")
	//	name := flag.String("name", "", "if you want to sign up")
	botid := flag.Int64("botid", -1, "botid")
	pwhite := flag.Bool("white", false, "white")
	pgray := flag.Bool("gray", false, "gray")
	pblack := flag.Bool("black", false, "black")
	gameplayid := flag.Int64("gameid", -1, "gameid")
	var pcolor game.Color
	if *pblack {
		pcolor = game.Black
		if *pgray || *pwhite {
			panic("black and gray or white")
		}
	} else if *pgray {
		pcolor = game.Gray
		if *pwhite {
			panic(*pwhite)
		}
	} else if *pwhite {
		pcolor = game.White
	}

	flag.Parse()
	c := client.NewClient(nil, *bu)
	//	var u, p int64
	var u int64
	var a []byte
	/*
		var err error
		if len(*name) > 0 {
			l, _, err := c.SignUp(multi.SignUpPost{*login, *passwd, *name})
			if err == nil {
				u, p, a = l.User, l.Player, l.Auth
			}
		}
		if len(*name) == 0 || err != nil {
			log.Println(err)
	*/
	gpg, _, err := c.Play(*gameplayid)
	if err != nil {
		log.Fatal(err)
	}
	gpst, _, err := c.State(gpg.State)
	if err != nil {
		log.Fatal(err)
	}
	if pcolor == game.ZeroColor {
		pcolor = gpst.MovesNext
	}
	ll, _, err := c.LogIn(multi.LoggingIn{*login, *passwd})
	if err != nil {
		log.Fatalln(err)
	}
	u, a = ll.ID, ll.AuthKey
	//ppp := c.UserInfo(u)
	//p = ppp.Player
	//}
	//log.Println(u, p, a)
	log.Println(u, a)
	log.Println("botid", *botid)
	bbb, _, err := c.BotKey(multi.BotKeyGetting{*botid, multi.Authorization{u, a}})
	if err != nil {
		log.Fatal(err)
	}
	binf, _, err := c.BotInfo(*botid)
	log.Println(binf)
	if err != nil {
		log.Fatal(err)
	}
	bbp, bba := bbb.ID, bbb.AuthKey
	echn := make(chan error)
	endchn := make(chan bool)
	go func() {
		for u := range echn {
			log.Println(u)
		}
	}()
	var aii ai.Player
	if string(binf.WhoAmI[:11]) == "constsitval" {
		aii = new(constsitval.AIPlayer)
	} else {
		log.Fatal(binf)
	}
	//load aiconf
	err = aii.SetConf(binf.Settings)
	if err != nil {
		log.Fatal(err)
	}
	yg, err := remote.New(
		c,
		aii,
		pcolor,
		*gameplayid,
		multi.Authorization{bbp, bba},
		func(g *remote.G) (int64, error) {
			log.Println("AFTFUNCC")
			for {
				log.Println("st for aftfunc", *g.state)
				a, _, err := g.C().After(
					g.gameid,
					[3]*int64{nil, nil, nil},
				)
				log.Println("aftfunc kitchen:", a, err)
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
		log.Fatal(err)
	}
	log.Println("end", <-endchn)
}
