package main

import "flag"
import "os"
import "github.com/ArchieT/3manchess/client"
import "fmt"

func main() {
	bu := flag.String("baseurl", os.Getenv("CHESSBASEURL"), "3manchess/multi base URL")
	gamid := flag.Int64("gameid", 1, "gameid")
	flag.Parse()
	cl := client.NewClient(nil, *bu)
	gid := *gamid
	u(gid, cl)
}

func u(ga int64, c *client.Client) {
	fmt.Println(ga)
	a, _, e := c.Play(ga)
	fmt.Println(a)
	if e != nil {
		return
	}
	b, _, e := c.State(a.State)
	fmt.Println(b, e)
	de, _, e := c.Before(ga)
	d := *de
	fmt.Println(d)
	if e != nil {
		return
	}
	for _, r := range d {
		u(r.MoveData.BeforeGame, c)
	}
}
