package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import (
	//	"io/ioutil"
	"log"
	"os"
)

var (
	//Debug is the only log that is always being put out on Stderr
	Debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
)
