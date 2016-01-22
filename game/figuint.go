package game

func (f *Fig) Uint8() uint8 {
	if f.FigType == 0 || f.Color == 0 {
		return 0
	}
	c := f.Color - 1
	t := f.FigType - 1
	var p uint8
	if f.PawnCenter && f.FigType == Pawn {
		p = 1
	} else {
		p = 0
	}
	return c*7 + t + p
}

func FigUint8(i uint8) Fig {
	var f Fig
	f.Color = Color(i/7 + 1)
	t := i % 7
	f.PawnCenter = t == 6
	if f.PawnCenter {
		f.FigType = Pawn
	} else {
		f.FigType = FigType(t + 1)
	}
	return f
}

func BoardUint(s *([6][24]uint8)) *Board {
	var b Board
	var t uint8
	var oac ACP
	for oac.OK() {
		t = (*s)[oac[0]][oac[1]]
		if t == 0 {
			b[oac[0]][oac[1]] = Square{NotEmpty: false, Fig{Color: ZeroColor, FigType: ZeroFigType}}
		} else {
			b[oac[0]][oac[1]] = Square{NotEmpty: true, FigUint8(t)}
		}
		oac.P()
	}
	return &b
}
