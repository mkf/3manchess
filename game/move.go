package game

type Move struct { //a struct describing a single move with the situation before it
	From        Pos
	To          Pos
	What        Fig
	AlreadyHere Fig
	Before      *State
}

//func (m *Move) String() string {
//}

func (m *Move) After() *State {   //situation after
	var movesnext Color
	if m.What.Color != m.Before.MovesNext {
		panic(m)
	}
	switch m.Before.MovesNext {
	case White:
		movesnext = Gray
	case Gray:
		movesnext = Black
	case Black:
		movesnext = White
	}
	if m.What.Color == m.AlreadyHere.Color {
		panic(m)
	}

}
