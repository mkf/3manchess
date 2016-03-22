package main

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import (
	"fmt"
	ai "github.com/ArchieT/3manchess/ai/sitvalues"
	"github.com/ArchieT/3manchess/game"
	deveng "github.com/ArchieT/3manchess/interface/devengchan"
	"github.com/ArchieT/3manchess/player"
	//	"os"
	"time"
)

func main() {
	fmt.Println("3manchess experimental engine")
	fmt.Println("Type 'h' for human, 'a' for AI")
	var ww, gw, bw rune
	fmt.Scanf("%c%c%c", &ww, &gw, &bw)
	var awhite, agrey, ablack ai.AIPlayer
	fmt.Println("enter three precisions for white, gray, black bot (even if not active) in scientific notation, e.g. 1234.456e-78")
	var wwp, gwp, bwp float64
	fmt.Scanf("%e %e %e", &wwp, &gwp, &bwp)
	awhite.Conf.Precision = wwp
	agrey.Conf.Precision = gwp
	ablack.Conf.Precision = bwp
	var white, grey, black deveng.Developer
	white.Name = "Whitey"
	grey.Name = "Greyey"
	black.Name = "Blackey"
	//players := map[game.Color]player.Player{game.White: player.Player(&white), game.Gray: player.Player(&grey), game.Black: player.Player(&black)}
	players := make(map[game.Color]player.Player)
	switch ww {
	case 'h':
		players[game.White] = player.Player(&white)
	case 'a':
		players[game.White] = player.Player(&awhite)
	}
	switch gw {
	case 'h':
		players[game.Gray] = player.Player(&grey)
	case 'a':
		players[game.Gray] = player.Player(&agrey)
	}
	switch bw {
	case 'h':
		players[game.Black] = player.Player(&black)
	case 'a':
		players[game.Black] = player.Player(&ablack)
	}
	end := make(chan bool)
	gp := player.NewGame(players, end)
	fmt.Println("NEW GAME:", gp)
	sitchacol := func(col *deveng.Developer) {
		for {
			fmt.Println(col, ", situation changed: ", <-col.SituationCh)
		}
	}
	go sitchacol(&white)
	go sitchacol(&grey)
	go sitchacol(&black)
	hurchacol := func(col *deveng.Developer) {
		for {
			fmt.Println(col, ", hurryup! ", <-col.HurryChan)
		}
	}
	go hurchacol(&white)
	go hurchacol(&grey)
	go hurchacol(&black)
	winchacol := func(col *deveng.Developer) {
		switch <-col.Result {
		case deveng.WIN:
			fmt.Println(col, "WON")
		case deveng.LOSE:
			fmt.Println(col, "LOST")
		case deveng.DRAW:
			fmt.Println(col, "DREW")
		case deveng.UNDEFRESULT:
			panic(fmt.Sprintln(col, "RESULT UNDEFINED"))
		}
	}
	go winchacol(&white)
	go winchacol(&grey)
	go winchacol(&black)

	go func() {
		var ff, ft, tf, tt int8
		for {
			select {
			case saskin := <-white.AskinForMove:
				fmt.Println(white, "is being asked for a move (ff ft tf tt): ")
				fmt.Scanf("%d %d %d %d", &ff, &ft, &tf, &tt)
				saskfto := game.FromTo{game.Pos{ff, ft}, game.Pos{tf, tt}}
				saskmov := saskfto.Move(saskin)
				white.HereRMoves <- saskmov
			case saskin := <-grey.AskinForMove:
				fmt.Println(grey, "is being asked for a move (ff ft tf tt): ")
				fmt.Scanf("%d %d %d %d", &ff, &ft, &tf, &tt)
				saskfto := game.FromTo{game.Pos{ff, ft}, game.Pos{tf, tt}}
				saskmov := saskfto.Move(saskin)
				grey.HereRMoves <- saskmov
			case saskin := <-black.AskinForMove:
				fmt.Println(black, "is being asked for a move (ff ft tf tt): ")
				fmt.Scanf("%d %d %d %d", &ff, &ft, &tf, &tt)
				saskfto := game.FromTo{game.Pos{ff, ft}, game.Pos{tf, tt}}
				saskmov := saskfto.Move(saskin)
				black.HereRMoves <- saskmov
			case <-time.After(time.Second * 20):
				var saskinrune rune = 0
				fmt.Print("hurry[y]?")
				fmt.Scanf("%c", &saskinrune)
				if saskinrune == 'y' {
					gp.HurryUpWhoever()
				}
			}
		}
	}()
	_ = <-end
}
