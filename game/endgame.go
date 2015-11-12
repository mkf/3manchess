package game

func (s *State) CanIMoveWOCheck() bool {
	var i, j, k, l int8
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			from = Pos{i, j}
			for k = 0; k < 6; k++ {
				for l = 0; l < 24; l++ {
					to = Pos{k, l}
					if s.AnyPiece(from, to) {
						_, err := Move{from, to, s}.After()
						if err == nil {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
