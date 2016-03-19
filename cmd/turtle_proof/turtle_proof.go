package main

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/ai/constsitval"
import "time"
import "log"

func main() {
	var a constsitval.AIPlayer
	a.Name = "Bot testowy"
	a.Conf = constsitval.AIConfig{
		Depth:             constsitval.DEFFIXDEPTH,
		OwnedToThreatened: constsitval.DEFOWN2THRTHD,
	}
	hurry := make(chan bool)
	newgame := game.NewState()
	go func() {
		time.Sleep(time.Minute)
		hurry <- true
	}()
	move := a.HeyItsYourMove(&newgame, hurry)
	log.Println(move)
}
