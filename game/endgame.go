package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//CanIMoveWOCheck — is there any move that would not end up in a check?
func (s *State) CanIMoveWOCheck(who Color) bool {
	for oac, loacp := range AMFT {
		if s.Board.GPos(Pos(oac)).Fig.Color == who {
			for _, oacp := range loacp {
				m := Move{oac, oacp, s, Queen}
				if _, err := m.After(); err == nil {
					return true
				}
			}
		}
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
