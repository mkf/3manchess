package game

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	//Debug is the only log that is always being put out on Stderr
	Debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
)
