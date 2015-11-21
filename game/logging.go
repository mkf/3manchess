package game

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	//Debug is the only log that is always being put out on Stderr
	Debug        = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
	BodyTrace    = log.New(ioutil.Discard, "BODY:TRACE ", log.LstdFlags)
	EndgameTrace = log.New(ioutil.Discard, "ENDGAME:TRACE ", log.LstdFlags)
	MoveTrace    = log.New(ioutil.Discard, "MOVE:TRACE ", log.LstdFlags)
	PossibTrace  = log.New(ioutil.Discard, "POSSIB:TRACE ", log.LstdFlags)
	SimplTrace   = log.New(ioutil.Discard, "SIMPL:TRACE ", log.LstdFlags)
	StateTrace   = log.New(ioutil.Discard, "STATE:TRACE ", log.LstdFlags)
)
