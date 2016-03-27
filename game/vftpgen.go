package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//VFTPGen : generates all valid FromToProm, given the game state
func VFTPGen(gamestate *State) <-chan FromToProm {
	allValid := make(chan FromToProm)
	go func() {
		var oac ACFT
		for oac.OK() {
			ft := FromTo(oac)
			move := Move{ft.From(), ft.To(), gamestate, Queen}
			if _, err := move.After(); err == nil {
				fig := (*gamestate).Board.GPos(ft.From()).Fig
				if fig.FigType == Pawn && fig.PawnCenter && ft.From()[0] == 1 {
					allValid <- FromToProm{ft, Queen}
					allValid <- FromToProm{ft, Rook}
					allValid <- FromToProm{ft, Bishop}
					allValid <- FromToProm{ft, Knight}
				} else {
					allValid <- FromToProm{ft, ZeroFigType}
				}
			}
			oac.P()
		}
		close(allValid)
	}()
	return allValid
}
