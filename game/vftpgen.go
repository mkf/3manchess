package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//VFTPGen : generates all valid FromToProm, given the game state
func VFTPGen(gamestate *State) <-chan FromToProm {
	all_valid := make(chan FromToProm)
	go func() {
		var oac ACFT
		for oac.OK() {
			ft := FromTo(oac)
			move := Move{ft.From(), ft.To(), gamestate, Queen}
			if _, err := move.After(); err == nil {
				all_valid <- FromToProm{ft, Queen}
				fig := (*gamestate).Board.GPos(ft.From()).Fig
				if fig.FigType == Pawn && fig.PawnCenter && ft.From()[0] == 1 {
					all_valid <- FromToProm{ft, Rook}
					all_valid <- FromToProm{ft, Bishop}
					all_valid <- FromToProm{ft, Knight}
				}
			}
			oac.P()
		}
		close(all_valid)
	}()
	return all_valid
}
