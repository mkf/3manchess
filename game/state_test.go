package game

import "testing"

func TestFromData_newgame(t *testing.T) {
	ns := NewState()
	nd := ns.Data()
	var xs State
	xs.FromData(nd)
	if !xs.Equal(&ns) {
		t.Fatal(ns, ns.Board, nd, xs, xs.Board)
	}
}
