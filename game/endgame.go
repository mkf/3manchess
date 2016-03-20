package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//CanIMoveWOCheck — is there any move that would not end up in a check?
func (s *State) CanIMoveWOCheck(who Color) bool {
	var oac ACP
	for oac.OK() {
		if s.Board.GPos(Pos(oac)).Fig.Color != who {
			oac.P()
			continue
		}
		var oacp ACP
		for oacp.OK() {
			if s.AnyPiece(Pos(oac), Pos(oacp)) {
				m := Move{Pos(oac), Pos(oacp), s, Queen}
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
