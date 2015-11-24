package game

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	//Debug is the only log that is always being put out on Stderr
	Debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
	//BodyTrace is the not saved by default log of body.go without errors
	BodyTrace = log.New(ioutil.Discard, "BODY:TRACE ", log.LstdFlags)
	//EndgameTrace is the not saved by default log of endgame.go without errors
	EndgameTrace = log.New(ioutil.Discard, "ENDGAME:TRACE ", log.LstdFlags)
	//MoveTrace is the not saved by default log of move.go without errors
	MoveTrace = log.New(os.Stderr, "MOVE:TRACE ", log.LstdFlags) //enabled
	//PossibTrace is the not saved by default log of possibilities.go without errors
	PossibTrace = log.New(ioutil.Discard, "POSSIB:TRACE ", log.LstdFlags)
	//SimplTrace is the not saved by default log of simple.go without errors
	SimplTrace = log.New(ioutil.Discard, "SIMPL:TRACE ", log.LstdFlags)
	//StateTrace is the not saved by default log of state.go without errors
	StateTrace = log.New(ioutil.Discard, "STATE:TRACE ", log.LstdFlags)
)
