package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

import "fmt"

//FigType : type of a figure, piece type
type FigType uint8

//const FigTypes
const (
	ZeroFigType FigType = iota
	Rook
	Knight
	Bishop
	Queen
	King
	Pawn
)

func (f *Fig) String() string {
	r := f.Rune()
	otoczka := "."
	if f.PawnCenter {
		otoczka = "!"
	}
	return otoczka + string(r)
}

//PawnCenter : whether the pawn had already passed through the center
type PawnCenter bool

//Byte returns "Y" or "N"
func (pc PawnCenter) Byte() byte {
	return ynbool(bool(pc))
}

//Fig : a struct describing a single piece: it's type, it's color, and, in case of a pawn, whether is had already passed through the center
type Fig struct {
	FigType
	Color
	PawnCenter
}

//Square : a struct describing a single square: whether it is empty, and what is on it
type Square struct {
	Fig
	NotEmpty bool
}

//Empty : return !s.NotEmpty
func (s Square) Empty() bool {
	return !s.NotEmpty
}

//Color : legacy : return s.Fig.Color
func (s Square) Color() Color {
	return s.Fig.Color
}

//What : legacy : return s.Fig.FigType
func (s Square) What() FigType {
	return s.Fig.FigType
}

//EMPTYOURSTR is the string that is the value of Square.String() if Square.Empty()
var EMPTYOURSTR = "__"

func (s Square) String() string {
	if s.NotEmpty {
		//ourbyte = []byte{byte(s.Fig.Color), byte(s.Fig.FigType), s.Fig.PawnCenter.Byte()}
		return s.Fig.String()
	}
	return EMPTYOURSTR
}

func (b *Board) String() string {
	var s string = "\n"
	for i := 0; i < 6; i++ {
		s += fmt.Sprintln((*b)[i])
	}
	return s
}

//Pos : coordinates
type Pos [2]int8

//Pos.String : give a nice [0,0] string
func (p Pos) String() string {
	return fmt.Sprintf("[%v,%v]", p[0], p[1])
}
func (p Pos) AddVec(vec Vector) Pos { return vec.addTo(p) }

type Vector interface {
	Rank() int8
	File() int8
	addTo(from Pos) Pos
	Bool() bool
	Units(fromRank int8) <-chan Vector
	EmptiesFrom(from Pos) <-chan Pos
	Moats(from Pos) <-chan Color
	implementsVector()
}

func (v Vector) AddTo(from Pos) Pos { return from.AddVec(v) }

type ZeroVector struct{}

func (v ZeroVector) Rank() int8                   { return 0 }
func (v ZeroVector) File() int8                   { return 0 }
func (v ZeroVector) Bool() bool                   { return true }
func (v ZeroVector) addTo(from Pos) Pos           { return from }
func (v ZeroVector) Units(_ int8) <-chan Vector   { ch := make(chan Vector); close(ch); return ch }
func (v ZeroVector) EmptiesFrom(_ Pos) <-chan Pos { ch := make(chan Pos); close(ch); return ch }
func (v ZeroVector) Moats(_ Pos) <-chan Color     { ch := make(chan Color); close(ch); return ch }
func (v ZeroVector) implementsVector()            {}

type JumpVector interface {
	Vector
	itsAJumpVector()
}

func (v JumpVector) Units(_ int8) <-chan Vector {
	ch := make(chan Vector)
	go func() { ch <- v; close(ch) }()
	return ch
}

type CastlingVector interface {
	JumpVector
	itsACastlingVector()
	empties() []int8
}

func (v CastlingVector) Rank() int8 { return 0 }

const KFM int8 = 4 //King's file modulo Color (mod 8)

func (v CastlingVector) Moats(_ Pos) <-chan Color { ch := make(chan Color); close(ch); return ch }
func (v CastlingVector) EmptiesFrom(from Pos) <-chan Pos {
	ch := make(chan Pos)
	go func() {
		if from[1]%8 != KFM {
			close(ch)
		}
		add := from[1] - KFM
		for toempt := range v.empties() {
			ch <- Pos{0, add + toempt}
		}
	}()
	return ch
}

const QUEENSIDE_CASTLING_EMPTIES [3]int8 = [3]int8{3, 2, 1}
const KINGSIDE_CASTLING_EMPTIES [2]int8 = [2]int8{5, 6}

type QueenSideCastlingVector struct{}

func (v QueenSideCastlingVector) File() int8         { return -2 }
func (v QueenSideCastlingVector) empties() []int8    { return QUEENSIDE_CASTLING_EMPTIES[:] }
func (v QueenSideCastlingVector) isACastlingVector() {}

type KingSideCastlingVector struct{}

func (v KingSideCastlingVector) File() int8         { return 2 }
func (v KingSideCastlingVector) empties() []int8    { return KINGSIDE_CASTLING_EMPTIES[:] }
func (v KingSideCastlingVector) isACastlingVector() {}

type PawnVector interface {
	JumpVector
	isAPawnVector()
	ReqPC() bool //returns needed PawnCenter value
	ReqProm(rank int8) bool
}

type PawnLongJumpVector struct{}

func (v PawnLongJumpVector) Rank() int8  { return 2 }
func (v PawnLongJumpVector) File() int8  { return 0 }
func (v PawnLongJumpVector) ReqPC() bool { return false }
func (v PawnLongJumpVector) addTo(from Pos) Pos {
	if from[0] != 1 {
		panic(struct {
			from Pos
			v    PawnLongJumpVector
		}{from, v})
	}
	return Pos{3, from[1]}
}
func (v PawnLongJumpVector) EnPassantCapField(from Pos) Pos {
	if from[0] != 1 {
		panic(struct {
			from Pos
			v    PawnLongJumpVector
		}{from, v})
	}
	return Pos{2, from[1]}
}
func (v PawnLongJumpVector) Moats(_ Pos) <-chan Color { ch := make(chan Color); close(ch); return ch }
func (v PawnLongJumpVector) isAPawnVector()           {}

//func (v PawnLongJumpVector) Units(_ int8) <-chan Vector {
//	ch := make(chan Vector)
//	go func() { ch <- v; close(ch) }()
//	return ch
//}

func (v PawnLongJumpVector) EmptiesFrom(from Pos) <-chan Pos {
	ch := make(chan Pos)
	go func() {
		ch <- v.addTo(from)
		ch <- v.EnPassantCapField(from)
		close(ch)
	}()
	return ch
}

type KnightVector struct {
	//Towards the center, i.e. inwards
	inward bool

	//Positive file direction (switched upon mirroring)
	plusFile bool

	//One rank closer to the center?
	//(about that one more (twice instead of once) rank or file)
	centerOneCloser bool
}

func (v KnightVector) MoreRank() bool { return v.centerOneCloser == v.inward }
func (v KnightVector) MoreFile() bool { return !v.MoreRank() }

func (v KnightVector) Rank() int8 { return tI8(v.inward, 1, -2) + tI8(v.centerOneCloser, 1, 0) }
func (v KnightVector) File() int8 { return tI8(v.MoreFile(), 2, 1) * tI8(v.plusFile, 1, -1) }
func (v KnightVector) addTo(from Pos) Pos {
	if v.inward && (v.centerOneCloser && from[0] >= 4 || from[0] == 5) {
		if v.centerOneCloser {
			return Pos{
				(5 + 4) - from[0],
				(from[0] + tI8(v.plusFile, 2, -2) + 12) % 24,
			}
		}
		return Pos{
			5,
			(from[0] + tI8(v.plusFile, 2, -2) + 12) % 24,
		}
	}
	return from.AddVector([2]int8{v.Rank(), v.File()})
}
func (v KnightVector) Moat(from Pos) Color {
	to := from.AddVec(v)
	ourxoreq := xoreqFunc(from, to)
	return ourxoreq
}
func (v KnightVector) Moats(from Pos) <-chan Color {
	ch := make(chan Color)
	go func() { ch <- v.Moat(from); close(ch) }()
	return ch
}
func (v KnightVector) Units(_ int8) <-chan KnightVector {
	ch := make(chan KnightVector)
	go func() { ch <- v; close(ch) }()
	return ch
}
func (v KnightVector) EmptiesFrom(_ Pos) <-chan Pos    { ch := make(chan Pos); close(ch); return ch }
func (v KnightVector) EmptiesBetween(_ Pos) <-chan Pos { ch := make(chan Pos); close(ch); return ch }

type ContinuousVector interface {
	Vector
	Abs() int
	Units(fromRank int) <-chan ContinuousVector
	isContinuousVector()
}

func (v ContinuousVector) EmptiesBetween(from Pos) <-chan Pos {
	ch := make(chan Pos)
	go func() {
		pos := from
		nofrom := false
		for u := range v.Units(from[0]) {
			if nofrom {
				ch <- pos
			} else {
				nofrom = true
			}
			pos = pos.AddVec(u)
		}
	}()
}

type AxisVector struct{ T int8 }

func (v AxisVector) Direc() bool         { return v.T >= 0 }
func (v AxisVector) DirecSign() int8     { return sign(v.T) }
func (v AxisVector) Abs() int8           { return abs(v.T) }
func (v AxisVector) isContinuousVector() {}

type FileVector struct{ AxisVector }

func (v FileVector) File() int8 { return v.T }
func (v FileVector) Rank() int8 { return 0 }
func (v FileVector) Units(_ int8) <-chan FileVector {
	ch := make(chan FileVector)
	go func() {
		for i := v.Abs(); i > 0; i-- {
			ch <- FileVector(v.DirecSign())
		}
		close(ch)
	}()
	return ch
}
func (v FileVector) Moats(from Pos) <-chan Color {
	ch := make(chan Color)
	go func() {
		if from[0] == 0 {
			left := from[1] % 8
			var tm int8
			if v.Direc() {
				tm = 8 - left
			} else {
				tm = left
			}
			start := Color(from[1]/8 + 1)
			moating := v.Abs() - tm
			if moating > 0 {
				if v.Direc() {
					ch <- start
				} else {
					ch <- start.Next().Next()
				}
				if moating > 8 {
					ch <- start.Next()
					if moating > 16 {
						if v.Direc() {
							ch <- start.Next().Next()
						} else {
							ch <- start
						}
					}
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (v FileVector) addTo(pos Pos) { return pos.AddVector([2]int8{0, v.File()}) }

//Board : board array
type Board [6][24]Square

//GPos gets the position's square pointer
func (b *Board) GPos(p Pos) *Square {
	if err := p.Correct(); err != nil {
		panic(err)
	}
	return &((*b)[p[0]][p[1]])
}

//Color : color type
type Color uint8

//Next returns the next color: White, Gray, Black,  White, etc.
func (c Color) Next() Color {
	return c%3 + 1
}

//String returns string "White"/"Gray"/"Black"
func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Gray:
		return "Gray"
	case Black:
		return "Black"
	default:
		panic(uint8(c))
	}
}

//Color constants
const (
	ZeroColor Color = iota
	White
	Gray
	Black
)

//COLORS : array of colors, ordered	const
var COLORS = [3]Color{White, Gray, Black}

//FIRSTRANKNEWGAME : first rank, from 0mod8 to 7mod8
var FIRSTRANKNEWGAME = [8]FigType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

var BRIDGEDMOATS = MoatsState{true, true, true}

var AMFT map[Pos][]Pos = make(map[Pos][]Pos, 144)

func amftinit() {
	var b Board
	var from, to Pos
	for from[0] = 0; from[0] < 6; from[0]++ {
		for from[1] = 0; from[1] < 24; from[1]++ {
			AMFT[from] = make([]Pos, 0, 61) // 61 is the biggest encountered number of to's since 497ca04494eb470c5d1d5778453f7ae026cb00b9
			for to[0] = 0; to[0] < 6; to[0]++ {
				for to[1] = 0; to[1] < 24; to[1]++ {
					if b.queen(from, to, BRIDGEDMOATS) || b.knight(from, to, BRIDGEDMOATS) {
						AMFT[from] = append(AMFT[from], to)
					}
				}
			}
		}
	}
}

//BOARDFORNEWGAME — a newgame board
var BOARDFORNEWGAME Board //newgame board

func boardinit() { //initialize BOARDFORNEWGAME module pseudoconstant
	for ci, c := range COLORS {
		for fi, f := range FIRSTRANKNEWGAME {
			a := ci*8 + fi
			BOARDFORNEWGAME[0][a].FigType = f
			BOARDFORNEWGAME[0][a].Fig.Color = c
			BOARDFORNEWGAME[0][a].NotEmpty = true
			BOARDFORNEWGAME[1][a].Fig.Color = c
			BOARDFORNEWGAME[1][a].Fig.FigType = Pawn
			BOARDFORNEWGAME[1][a].NotEmpty = true
			for l := 2; l < 6; l++ {
				BOARDFORNEWGAME[l][a].Fig.Color = 0
				BOARDFORNEWGAME[l][a].Fig.FigType = 0
				BOARDFORNEWGAME[l][a].NotEmpty = false
			}
		}
	}
}

//NewBoard is a replacement for NEWBOARD
func NewBoard() Board {
	return BOARDFORNEWGAME
}
