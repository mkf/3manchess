package sitvalues

import "github.com/ArchieT/3manchess/game"

const DEATH float64 = 100000

const OPDIES float64 = 15000

var VALUES = map[game.FigType]int32{
	game.Pawn:   1,
	game.Knight: 3,
	game.Bishop: 6,
	game.Rook:   5,
	game.Queen:  10,
	game.King:   2400,
}

func SitValue(s *game.State) float64 {
	who := s.MovesNext.Next().Next()
	nasze, _ := s.Board.FriendsNAllies(who, s.PlayersAlive)
	myatakujem := s.Board.WeAreThreateningTypes(who, s.PlayersAlive, s.EnPassant)
	nasatakujo := s.Board.WeAreThreatened(who, s.PlayersAlive, s.EnPassant)
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
	zyjacy := float64(ostatecznie)
	if !s.PlayersAlive.Give(who) {
		return -DEATH
	}
	if !s.PlayersAlive.Give(who.Next()) {
		zyjacy += OPDIES
	}
	if !s.PlayersAlive.Give(who.Next().Next()) {
		zyjacy += OPDIES
	}
	return zyjacy
}
