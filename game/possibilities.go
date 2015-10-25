package game

func (b *Board) Straight(from Pos, to Pos, m MoatsState) (bool, bool) { //(whether it can, whether it can capture/check)
	var cantech, canmoat, canfig bool
	capcheck := true
	if from == to {
		panic("Same square!")
	}
	if from[0] == to[0] {
		cantech = true
		if from[0] == 0 {
			var mshort, mlong, capcheckshort bool
			if from[1]/8 == to[1]/8 {
				capcheckshort = true
				canmoat = true
				mshort = true
				if m[0] && m[1] && m[2] {
					mlong = true
				}
			} else {
				capcheckshort = false
				fromto := [2]int8{from[1] / 8, to[1] / 8}
				switch fromto {
				case [2]int8{0, 1}, [2]int8{1, 0}:
					mshort = m[1]
					mlong = m[0] && m[2]
				case [2]int8{1, 2}, [2]int8{2, 1}:
					mshort = m[2]
					mlong = m[0] && m[1]
				case [2]int8{2, 0}, [2]int8{0, 2}:
					mshort = m[0]
					mlong = m[1] && m[2]
				}
			}
			var direcshort int8
			var fromtominus int8
			if capcheckshort {
				direcshort = sign(to[1]-from[1])
			} else {
				fromtominus = fromto[1]-fromto[0]
				if abs(fromtominus)==2 {
					fromtominus=-fromtominus
				}
				direcshort = sign(fromtominus)
			}
			//zrobić zamienne capfig/capmoat jeśli jedno to drugie było(próbować)
		} else {
			canmoat = true
			canfig = true
			for i := from[1] + 1; ((i-from[1])%24 < (to[1]-from[1])%24) && canfig; i = (i + 1) % 24 {
				go func() {
					if canfig && !((*b)[from[0]][i].Empty()) {
						canfig = false
					}
				}()
			}
			canfiga := true
			for i := from[1] - 1; ((i-from[1])%24 > (to[1]-from[1])%24) && canfiga; i = (i - 1) % 24 {
				go func() {
					if canfiga && !((*b)[from[0]][i].Empty()) {
						canfiga = false
					}
				}()
			}
			canfig = canfig || canfiga
		}
	} else if from[1] == to[1] {
		cantech = true
		canmoat = true
		canfig = true
		sgn := sign(to[0] - from[0])
		for i := from[0] + sgn; (sgn*i < to[0]) && canfig; i += sgn {
			go func() {
				if canfig && !((*b)[i][from[1]].Empty()) {
					canfig = false
				}
			}()
		}
	} else if ((from[1] - 12) % 24) == to[1] {
		cantech = true
		canmoat = true
		canfig = true
		for i, j := from[0], to[0]; canfig && (i < 6 && j < 6); i, j = i+1, j+1 {
			go func() {
				go func() {
					if canfig && !((*b)[i][from[1]].Empty()) {
						canfig = false
					}
				}()
				go func() {
					if canfig && !((*b)[j][to[1]].Empty()) {
						canfig = false
					}
				}()
			}()
		}
	} else {
		cantech = false
	}
	return cantech&&canmoat&&canfig,
}

//func (b Board) Diagonal
