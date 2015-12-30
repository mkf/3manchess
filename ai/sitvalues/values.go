package sitvalues

import "github.com/ArchieT/3manchess/game"

var VALUES = map[game.FigType]int32{
	game.Pawn:   1,
	game.Knight: 3,
	game.Bishop: 6,
	game.Rook:   5,
	game.Queen:  10,
	game.King:   2400,
}

func SitValue(s *game.State) int32 {
	nasze := s.Board.FriendsNAllies(s.MovesNext, s.PlayersAlive)
	myatakujem := s.Board.WeAreThreateningTypes(s.MovesNext, s.PlayersAlive, s.EnPassant)
	nasatakujo := s.Board.WeAreThreatened(s.MovesNext, s.PlayersAlive, s.EnPassant)
	var own, myich, oninas, ostatecznie int32
	for _, o := range nasze {
		own += VALUES[(*s.Board)[o[0]][o[1]].FigType]
	}
	for n := range myatakujem {
		myich += VALUES[n]
	}
	for n := range nasatakujo {
		oninas += VALUES[n]
	}
	ostatecznie = own + myich - oninas
	return ostatecznie
}
