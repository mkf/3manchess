package game

type State struct {
	Board [8][24][2]byte //[color,figure_lowercase] //divided:BlackGrayWhite
	MovesNext byte
	Castling [2]byte //[color,figure_lowercase]
	EnPassant [2][2]uint8 //[previousplayer,currentplayer]  [number,letter]
	HalfmoveClock uint8
	FullmoveNumber uint16
}

//func (s State) String() string {
//}
