package game

type State struct {
	Board [8][24][2]byte //[color,figure_lowercase] //divided:BlackGrayWhite
	MovesNext byte //W G B
	Castling [2]byte //[color,figure_lowercase]
	EnPassant [2][2]uint8 //[previousplayer,currentplayer]  [number,letter]
	HalfmoveClock uint8
	FullmoveNumber uint16
}

//func (s *State) String() string {
//}

type Move struct {
	From [2]uint8
	To [2]uint8
	What byte
	Before *State
}

//func (m *Move) String() string {
//}

func (m *Move) After() *State {
	white := false
	gray := false
	black := false
	switch i:=m.Before.MovesNext {
	}
}
