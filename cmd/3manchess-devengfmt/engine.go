package main

import (
	"fmt"
	//	"github.com/ArchieT/3manchess/ai"
	"github.com/ArchieT/3manchess/game"
	deveng "github.com/ArchieT/3manchess/interface/devengfmt"
	"github.com/ArchieT/3manchess/player"
	//	"os"
)

func main() {
	fmt.Println("3manchess experimental engine")
	/*
		fmt.Println("Type 'h' for human, 'a' for AI")
		var ww, gw, bw rune
		fmt.Scanf("%c%c%c", &ww, &gw, &bw)
		var awhite, agrey, ablack ai.AIPlayer
		fmt.Println("enter three precisions for white, gray, black bot (even if not active) in scientific notation, e.g. 1234.456e-78")
		var wwp, gwp, bwp float64
		fmt.Scanf("%e %e %e", &wwp, &gwp, &bwp)
		awhite.FixedPrecision = wwp
		agrey.FixedPrecision = gwp
		ablack.FixedPrecision = bwp
	*/
	var white, grey, black deveng.Developer
	white.Name = "Whitey"
	grey.Name = "Greyey"
	black.Name = "Blackey"
	//players := map[game.Color]player.Player{game.White: player.Player(&white), game.Gray: player.Player(&grey), game.Black: player.Player(&black)}
	players := make(map[game.Color]player.Player)
	/*
		switch ww {
		case 'h':
	*/
	players[game.White] = player.Player(&white)
	/*
		case 'a':
			players[game.White] = player.Player(&awhite)
		}
		switch gw {
		case 'h':
	*/
	players[game.Gray] = player.Player(&grey)
	/*
		case 'a':
			players[game.Gray] = player.Player(&agrey)
		}
		switch bw {
		case 'h':
	*/
	players[game.Black] = player.Player(&black)
	/*
		case 'a':
			players[game.Black] = player.Player(&ablack)
		}
	*/
	end := make(chan bool)
	gp := player.NewGame(players, end)
	fmt.Println("NEW GAME:", gp)
	/*
		go func() {
			for {
				var heyhurry bool
				fmt.Scanf("%t", &heyhurry)
				for _, icol := range game.COLORS {
					go func() {
						players[icol].HurryChannel() <- heyhurry
						recover()
					}()
				}
			}
		}()
	*/
	_ = <-end
}
