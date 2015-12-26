package game

//CanIMoveWOCheck — is there any move that would not end up in a check?
func (s *State) CanIMoveWOCheck(who Color) bool {
	var i, j, k, l int8
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			if tojefig := (*(s.Board))[i][j].Fig; tojefig.Color != who {
				continue
			}
			from := Pos{i, j}
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					to := Pos{k, l}
					if s.AnyPiece(from, to) {
						m := Move{from, to, s}
						_, err := m.After()
						if err == nil {
							EndgameTrace.Println("Yes, you can move without check!")
							return true
						}
					}
				}
			}
		}
	}
	EndgameTrace.Println("No, you cannot move without check!")
	return false
}

//AmIInCheck — Am I in check right now?
func (s *State) AmIInCheck(who Color) bool {
	return s.Board.CheckChecking(who, s.PlayersAlive)
}