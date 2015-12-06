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

//Winner : return a string with a brief description of the result
func Winner(state *State) string {
	var number_of_winners, last_winner, first_winner uint8
	first_winner = 9
	number_of_winners = 0
	for i := 0 ; i < len(state.PlayersAlive); i++ {
		if state.PlayersAlive[i] {
			number_of_winners++
			if first_winner == 9 {
				first_winner = uint8(i)
			}
			last_winner = uint8(i)
		}
	}
	var answer string
	switch number_of_winners {
	case 1:
		answer = ColorUint8(last_winner).String() + " wins"
	case 2:
		answer = ColorUint8(first_winner).String() + " and " + ColorUint8(last_winner).String() + " tie"
	default:
		panic("Game isn't finished.")
	}
	return answer
}
