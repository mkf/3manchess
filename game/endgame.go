package game

//CanIMoveWOCheck — is there any move that would not end up in a check?
func (s *State) CanIMoveWOCheck(who Color) bool {
	var oac, oacp ACP
	var from Pos
	for oac.OK() {
		from = Pos(oac)
		if s.Board.GPos(from).Fig.Color != who {
			continue
		}
		oacp = ACP{0, 0}
		for oacp.OK() {
			to := Pos(oacp)
			if s.AnyPiece(from, to) {
				m := Move{from, to, s, Queen}
				_, err := m.After()
				if err == nil {
					return true
				}
			}
			oacp.P()
		}
		oac.P()
	}
	return false
}

//Check type contains the important thing, that is If/Bool(), and a descriptive field From [Pos]
type Check struct {
	If   bool
	From Pos
}

//Bool returns the If fields, who knows if it will work with bool(Check)
func (c Check) Bool() bool {
	return c.If
}

//AmIInCheck — Am I in check right now?
func (s *State) AmIInCheck(who Color) Check {
	return s.Board.CheckChecking(who, s.PlayersAlive)
}
