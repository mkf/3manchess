package sitvalues

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "github.com/ArchieT/3manchess/game"

const DEATH float64 = -100000

const OPDIES float64 = 15000

var VALUES = map[game.FigType]int32{
	game.Pawn:   1,
	game.Knight: 3,
	game.Bishop: 5,
	game.Rook:   6,
	game.Queen:  10,
	game.King:   3,
}

func (a *AIPlayer) SitValue(s *game.State, who game.Color) float64 {
	nasze, ich := s.Board.FriendsNAllies(who, s.PlayersAlive)
	myatakujem := s.Board.WeAreThreateningTypes(who, s.PlayersAlive, s.EnPassant)
	nasatakujo := s.Board.WeAreThreatened(who, s.PlayersAlive, s.EnPassant)
	var own, theirs, myich, oninas int32
	var zyjacy float64
	for _, o := range nasze {
		own += VALUES[(*s.Board)[o[0]][o[1]].FigType]
	}
	for o := range ich {
		theirs += VALUES[(*s.Board)[o[0]][o[1]].FigType]
	}
	for n := range myatakujem {
		myich += VALUES[n]
	}
	for n := range nasatakujo {
		oninas += VALUES[n]
	}
	zyjacy = float64(own-theirs)*(a.Conf.OwnedToThreatened) + float64(myich-oninas)
	if !s.PlayersAlive.Give(who) {
		return DEATH
	}
	for _, player := range game.COLORS {
		if player != who && !s.PlayersAlive.Give(player) {
			zyjacy += OPDIES
		}
	}
	return zyjacy
}
