package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

func (f *Fig) Uint8() uint8 {
	if f.FigType == 0 || f.Color == 0 {
		return 0
	}
	c := uint8(f.Color) - 1
	t := uint8(f.FigType) - 1
	var p uint8
	if f.PawnCenter && f.FigType == Pawn {
		p = 1
	} else {
		p = 0
	}
	return c*7 + t + p
}

func (s *Square) Uint8() uint8 {
	if s.Empty() {
		return 0
	}
	return s.Fig.Uint8()
}

func SqUint8(i uint8) Square {
	if i == 0 {
		return Square{Fig: Fig{FigType: 0, Color: 0, PawnCenter: false}, NotEmpty: false}
	}
	return Square{Fig: FigUint8(i), NotEmpty: true}
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
		b[oac[0]][oac[1]] = SqUint8(t)
		oac.P()
	}
	return &b
}

func BoardByte(s []byte) *Board {
	var b Board
	var t uint8
	var oac ACP
	if len(s) < 24*6 {
		panic(len(s))
	}
	for oac.OK() {
		t = s[24*oac[0]+oac[1]]
		b[oac[0]][oac[1]] = SqUint8(t)
		oac.P()
	}
	return &b
}

func (b *Board) Byte() []byte {
	d := make([]byte, 0, 24*6)
	var oac ACP
	for oac.OK() {
		d = append(d, b.GPos(Pos(oac)).Uint8())
		oac.P()
	}
	return d
}

type BadBoardForTemplate struct {
	Zero  [24]uint8
	One   [24]uint8
	Two   [24]uint8
	Three [24]uint8
	Four  [24]uint8
	Five  [24]uint8
}

func (b *Board) BadBoardForTemplate() *BadBoardForTemplate {
	var bb BadBoardForTemplate
	pl := [6]*[24]uint8{&(bb.Zero), &(bb.One), &(bb.Two), &(bb.Three), &(bb.Four), &(bb.Five)}
	var i int8
	for i = 0; i < 6; i++ {
		b.bbarray(i, pl[i])
	}
	return &bb
}

func (b *Board) bbarray(rank int8, arr *[24]uint8) {
	var i int8
	for i = 0; i < 24; i++ {
		(*arr)[i] = (*b)[rank][i].Uint8()
	}
}
