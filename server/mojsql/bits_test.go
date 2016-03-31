package mojsql

import "testing"

type przyp struct {
	arg arg
	res []bool
}

type arg struct {
	source uint8
	lenght int8
}

var przypadki = []przyp{
	{arg{0, 3}, []bool{false, false, false}},
	{arg{0, 6}, []bool{false, false, false, false, false, false}},
	{arg{127, 6}, []bool{true, true, true, true, true, true}},
	//	{arg{28,3},[]bool{_,_,_}},
}

func compsl(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIntbit_przypadki(t *testing.T) {
	for _, pr := range przypadki {
		t.Log(pr)
		res := intbit(pr.arg.source, pr.arg.lenght)
		t.Log("r:", res)
		if !compsl(res, pr.res) {
			t.Fatal(false)
		}
	}
}
