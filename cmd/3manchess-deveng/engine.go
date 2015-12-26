package main

import (
	"fmt"
	"github.com/ArchieT/3manchess/game"
	"github.com/ArchieT/3manchess/interface/deveng"
	"github.com/ArchieT/3manchess/player"
	//	"os"
)

func main() {
	fmt.Println("3manchess experimental engine")
	var white, grey, black deveng.Developer
	white.Name = "Whitey"
	grey.Name = "Greyey"
	black.Name = "Blackey"
	players := map[game.Color]player.Player{game.White: player.Player(&white), game.Gray: player.Player(&grey), game.Black: player.Player(&black)}
	end := make(chan bool)
	player.NewGame(players, end)
	_ = <-end
}
