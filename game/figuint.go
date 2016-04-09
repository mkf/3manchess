package game

import "strconv"

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//Uint8 returns   [[ _ P C C C T T T ]]
func (f *Fig) Uint8() uint8 {
	return (bool2uint8(bool(f.PawnCenter)) << 6) + (uint8(f.Color) << 3) + uint8(f.FigType)
}

func bool2uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

//Uint8 returns an uint8 repr of a Square
func (s *Square) Uint8() uint8 {
	if s.Empty() {
		return 0
	}
	return s.Fig.Uint8()
}

//MarshalJSON makes *Square fulfill the Marshaler interface with sq.Uint8()
func (s *Square) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(s.Uint8()))), nil
}

//UnmarshalJSON makes *Square fulfill the Unmarshaler interface wt=with FromUint8(u)
func (s *Square) UnmarshalJSON(b []byte) error {
	i, err := strconv.Atoi(string(b))
	s.FromUint8(uint8(i))
	return err
}

//FromUint8 changes the value to SqUint8
func (s *Square) FromUint8(u uint8) {
	*s = SqUint8(u)
}

//SqUint8 reproduces a Square from an uint8 repr
func SqUint8(i uint8) Square {
	if i == 0 {
		return Square{Fig: Fig{FigType: 0, Color: 0, PawnCenter: false}, NotEmpty: false}
	}
	return Square{Fig: FigUint8(i), NotEmpty: true}
}

//FigUint8 reproduces a Fig from an uint8 repr
func FigUint8(i uint8) Fig {
	var f Fig
	f.PawnCenter = PawnCenter((i >> 7) > 0)
	f.Color = Color((i >> 3) & 7)
	f.FigType = FigType(i & 7)
	return f
}

//BoardUint reproduces a Board from 2d array repr
func BoardUint(s *([6][24]uint8)) *Board {
	var b Board
	var t uint8
	for _, pos := range ALLPOS {
		t = (*s)[pos[0]][pos[1]]
		b[pos[0]][pos[1]] = SqUint8(t)
	}
	return &b
}

func byteoac(oac Pos) uint8 { return (24 * uint8(oac[0])) + uint8(oac[1]) }

//BoardByte reproduces a Board from byte slice repr
func BoardByte(s []byte) *Board {
	var b Board
	var t uint8
	if len(s) != 24*6 {
		panic(len(s))
	}
	for _, pos := range ALLPOS {
		t = s[byteoac(pos)]
		b[pos[0]][pos[1]] = SqUint8(t)
	}
	return &b
}

//Byte returns all 6 concatenated ranks, where each rank is 24 squares, each represented by Square.Uint8
func (b *Board) Byte() [144]byte {
	var d [144]byte
	for _, pos := range ALLPOS {
		d[byteoac(pos)] = b.GPos(pos).Uint8()
	}
	return d
}

//BBArray puts the selected rank's repr into arr
func (b *Board) BBArray(rank int8, arr *[24]uint8) {
	var i int8
	for i = 0; i < 24; i++ {
		(*arr)[i] = (*b)[rank][i].Uint8()
	}
}
