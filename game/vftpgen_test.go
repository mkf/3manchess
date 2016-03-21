package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "testing"

func TestVFTPGen_newgame(t *testing.T) {
	newstate := NewState()
	for r := range VFTPGen(&newstate) {
		t.Log(r)
	}
}
