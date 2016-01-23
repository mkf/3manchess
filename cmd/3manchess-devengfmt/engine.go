package main

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import (
	"fmt"
	"github.com/ArchieT/3manchess/game"
	deveng "github.com/ArchieT/3manchess/interface/devengfmt"
	"github.com/ArchieT/3manchess/player"
	//	"os"
)

func main() {
	fmt.Println("3manchess experimental engine")
	var white, grey, black deveng.Developer
	white.Name = "Whitey"
	grey.Name = "Greyey"
	black.Name = "Blackey"
	players := make(map[game.Color]player.Player)
	players[game.White] = player.Player(&white)
	players[game.Gray] = player.Player(&grey)
	players[game.Black] = player.Player(&black)
	end := make(chan bool)
	gp := player.NewGame(players, end)
	fmt.Println("NEW GAME:", gp)
	_ = <-end
}
