package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"
import "log"

func TestASAOMGen_newgame(t *testing.T) {
	newstate := NewState()
	for s := range ASAOMGen(&newstate, Black) {
		log.Println(s)
	}
}
