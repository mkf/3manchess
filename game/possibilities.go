package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

func (b *Board) canfigstraighthoriz(rank, from, to int8) (bool, bool) { //returns plus, minus
	return b.canfigstraighthorizdirec(rank, from, to, 1), b.canfigstraighthorizdirec(rank, from, to, -1)
}

func (b *Board) canfigstraighthorizdirec(r, f, t, d int8) bool { //rank, from file, to file, direction
	for i := (f + d + 24) % 24; i != t; i = (i + d + 24) % 24 {
		if (*b)[r][i].NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) straight(from Pos, to Pos, m MoatsState) bool { //(bool, bool) { //(whether it can, whether it can capture/check)
	var cantech, canmoat, canfig bool
	//capcheck := true
	if from == to { //Same square
		return false
	}
	if from[0] == to[0] { //same rank
		cantech = true
		if from[0] == 0 { //first rank
			var mshort, mlong bool //, capcheckshort bool
			var direcshort int8
			var fromtominus int8
			if from[1]>>3 == to[1]>>3 { //same color area
				mlong, direcshort, mshort, canmoat = m[0] && m[1] && m[2], sign(to[1]-from[1]), true, true
			} else { //moving to another color's area
				fromto := [2]int8{from[1] >> 3, to[1] >> 3}
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
				fromtominus = fromto[1] - fromto[0]
				if abs(fromtominus) == 2 {
					fromtominus = -fromtominus
				}
				direcshort = sign(fromtominus)
			}
			canfigplus, canfigminus := b.canfigstraighthoriz(0, from[1], to[1])
			canfig = canfigplus || canfigminus

			//as we are on the first rank && moving to another color's area, we gotta check the moats
			canmoat = canmoat || (direcshort == 1 && ((canfigplus && mshort) || (canfigminus && mlong))) ||
				(direcshort == -1 && ((canfigminus && mshort) || (canfigplus && mlong)))

			if canmoat {
				if b.GPos(to).NotEmpty {
					return false
				}
			}
			/*if canmoat {
				cheb := *b
				ourmoasq := new(Square)
				ourmoasq.FromUint8(0)
				*cheb.GPos(to),*cheb.GPos(from) = *cheb.GPos(from),*ourmoasq

			}
			*/

		} else { //if same rank, but not first rank
			canmoat = true
			canfigplus, canfigminus := b.canfigstraighthoriz(from[0], from[1], to[1])
			canfig = canfigplus || canfigminus
		}
	} else if from[1] == to[1] { //if the same file, ie. no passing through center
		cantech, canmoat = true, true
		canfig = b.canfigstraightvertnormal(from[1], from[0], to[0])
	}
	if cantech && canmoat && canfig {
		return true
	}
	if ((from[1] + 12) % 24) == to[1] { //if the adjacent file, passing through center
		return b.canfigstraightvertthrucenter(from[1], from[0], to[0])
	}
	return false
}

func (b *Board) canfigstraightvertthrucenter(s, f, t int8) bool { //startfile (from[0]), from, to
	e := (s + 12) % 24
	//searching for collisions from both sides of the center
	for i := f + 1; i < 6; i++ {
		if (*b)[i][s].NotEmpty {
			return false
		}
	}
	for i := t + 1; i < 6; i++ {
		if (*b)[i][e].NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) canfigstraightvertnormal(file, f, t int8) bool {
	s := sign(t - f)
	for i := f + s; i != t; i += s {
		if (*b)[i][file].NotEmpty {
			return false
		}
	}
	return true
}

var PLUSMINUSPAIRS = [4][2]int8{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

func (p Pos) AddVector(v [2]int8) Pos {
	return Pos{p[0] + v[0], (p[1] + v[1] + 24) % 24}
}

func (p Pos) MinusVector(v [2]int8) Pos {
	return p.AddVector([2]int8{-v[0], -v[1]})
}

/*
var datafordiagonal = [6][6][2]int8{ //fromrank, torank, short&long file vector lenght
	{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}},
	{{1, 1}, {0, 2}, {1, 3}, {2, 4}, {3, 5}, {4, 6}},
	{{2, 2}, {1, 3}, {0, 4}, {1, 5}, {2, 6}, {3, 7}},
	{{3, 3}, {2, 4}, {1, 5}, {0, 6}, {1, 7}, {2, 8}},
	{{4, 4}, {3, 5}, {2, 6}, {1, 7}, {0, 8}, {1, 9}},
	{{5, 5}, {4, 6}, {3, 7}, {2, 8}, {1, 9}, {0, 10}},
}
*/

func tablediagonal(fromrank, torank int8, longnotshort bool) int8 {
	if longnotshort {
		return fromrank + torank
	}
	return abs(fromrank - torank)
}

func techdiagonal(from, to Pos) (short, long bool, znak int8) {
	if from == to {
		return false, false, 1
	}
	shorttd := abs(from[0] - to[0])
	longtd := (from[0] + to[0])
	switch to[1] {
	case (from[1] + shorttd) % 24:
		short = true
		znak = 1
	case (from[1] - shorttd + 24) % 24:
		short = true
		znak = -1
	case (from[1] + longtd) % 24:
		long = true
		znak = 1
	case (from[1] - longtd + 24) % 24:
		long = true
		znak = -1
	}

	if short && (from[1]+znak*longtd+24)%24 == to[1] {
		long = true
	}
	return
}

func (b *Board) canfigshortdiagonal(from, to Pos, znak int8) bool {
	rankznak := sign(to[0] - from[0])
	for a := from.AddVector([2]int8{rankznak, znak}); a != to; a = a.AddVector([2]int8{rankznak, znak}) {
		if b.GPos(a).NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) canfiglongdiagonal(from, to Pos, znak int8) bool {
	czycent := (from[0] == 5)
	var a Pos
	if !czycent {
		a = from.AddVector([2]int8{1, -znak})
	} else {
		a = from
	}
	for ; a[0] < 5; a = a.AddVector([2]int8{1, -znak}) {
		if b.GPos(a).NotEmpty {
			return false
		}
	}
	if a[0] != 5 {
		panic(a)
	}
	if !czycent && b.GPos(a).NotEmpty {
		return false
	}
	a = a.AddVector([2]int8{0, znak * tablediagonal(5, 5, true)})
	for ; a != to; a = a.AddVector([2]int8{-1, -znak}) {
		if b.GPos(a).NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) canfigdiagonal(from, to Pos, znak int8, longnotshort bool) bool {
	if longnotshort {
		return b.canfiglongdiagonal(from, to, znak)
	}
	return b.canfigshortdiagonal(from, to, znak)
}

func bujemny(b bool) int8 {
	if b {
		return -1
	}
	return 1
}

func bdodatni(b bool) int8 {
	if b {
		return 1
	}
	return -1
}

func moatnumdiagonal(from, to Pos, z, longnotshort bool) int8 {
	switch int8(0) {
	case from[0]:
		switch from[1] % 8 {
		case 7:
			if z == longnotshort {
				return (from[1]>>3 + 1) % 3
			}
		case 0:
			if !z == longnotshort {
				return from[1] >> 3
			}
		}
	case to[0]:
		switch to[1] % 8 {
		case 7:
			if !z == longnotshort {
				return (to[1]>>3 + 1) % 3
			}
		case 0:
			if z == longnotshort {
				return to[1] >> 3
			}
		}
	}
	return -1
}

func moatnumsdiagonal(from, to Pos, l, s, z bool) int8 {
	var mns, mnl int8 = -1, -1
	if s {
		mns = moatnumdiagonal(from, to, z, false)
	}
	if l {
		mnl = moatnumdiagonal(from, to, z, true)
	}
	switch mns {
	case -1, mnl:
		return mnl
	}
	if mnl == -1 {
		return mns
	}
	return -1
}

func (b *Board) diagonal(from, to Pos, m MoatsState) bool {
	short, long, znak := techdiagonal(from, to)
	var canfigshort, canfiglong bool
	if !(short || long) {
		return false
	}
	if short {
		canfigshort = b.canfigshortdiagonal(from, to, znak)
	}
	if long {
		canfiglong = b.canfiglongdiagonal(from, to, znak)
	}
	if !(canfigshort || canfiglong) {
		return false
	}
	moatnum := moatnumsdiagonal(from, to, canfiglong, canfigshort, znak > 0)
	return moatnum == -1 || m[moatnum] && b.GPos(to).Empty()
}

func (b *Board) checkbadpc(from Pos, p PawnCenter) {
	n := b.GPos(from).Color()
	if n == Color(from[1]>>3+1) && p {
		panic("pS" + n.String())
	}
}

func (p PawnCenter) ujemny() int8 {
	if p {
		return -1
	}
	return 1
}

func (b *Board) pawnStraight(from Pos, to Pos, p PawnCenter) bool { //(bool,PawnCenter,EnPassant) {
	if from == to {
		return false
	}
	b.checkbadpc(from, p)
	if from[1] == to[1] {
		switch to[0] - from[0] {
		case +2:
			return !bool(p) && from[0] == 1 &&
				(*b)[2][from[1]].Empty() && b.GPos(to).Empty()
			//ep:=e.Appeared(Pos{2,from[1]})
		case p.ujemny():
			return b.GPos(to).Empty()
		default:
			return false //,p,e
		}
	}
	return ((from[1]+12)%24) == to[1] && from[0] == 5 && to[0] == 5 && !bool(p) && b.GPos(to).Empty()
}

func (b *Board) kingMove(from Pos, to Pos, m MoatsState) bool {
	return from != to && b.queen(from, to, m) && absu(from[0]-to[0]) <= 1 &&
		((absu(from[1]-to[1]) == 23 || absu(from[1]-to[1]) <= 1) ||
			// king isn't moving to adjacent or current file (such as from file 1 to 24 or vice versa)
			(from[0] == 5 && to[0] == 5 && // king is moving through the center
				((from[1]+12)%24 == to[1] || // king movin fwd thru center
					((from[1]+10)%24 == to[1] || (from[1]-10+24)%24 == to[1])))) // king movin diag thru center
}

func pawncreek(from Pos, tof int8) bool {
	switch [2]int8{from[1] % 8, tof % 8} { //przemieszczenie między rankami
	case [2]int8{0, 7}, [2]int8{7, 0}:
		return from[0] < 3 //otrze się o creek
	}
	return true //normalnie nie musi się męczyć z creekami
	//założenie: creeki są aktywne nawet jak już nie ma moatów
}

func (e EnPassant) ktorys(to Pos) Pos {
	if e[0] == to {
		return e[0]
	}
	return e[1]
}

func (b *Board) pawnCapture(from Pos, to Pos, e EnPassant, p PawnCenter) bool {
	if from == to {
		return false
	}
	ep := e.ktorys(to)
	nasz := b.GPos(from).Fig.Color
	switch to[0] {
	case from[0]:
		if from[0] == 5 && !p {
			switch to[1] {
			case (from[1] + 24 - 10) % 24, (from[1] + 10) % 24:
				return b.GPos(to).NotEmpty && b.GPos(to).Fig.Color != nasz
			}
		}
	case from[0] + p.ujemny():
		if bool(p) || pawncreek(from, to[1]) {
			switch to[1] {
			case (from[1] + 1) % 24, (from[1] + 24 - 1) % 24:
				switch to {
				case ep:
					ep3, ep2 := (*b)[3][ep[1]], (*b)[2][ep[1]]
					return ep3.FigType == Pawn && ep3.NotEmpty && ep3.Fig.Color != nasz && ep2.Empty()
				}
				doc := b.GPos(to)
				return doc.NotEmpty && doc.Color() != nasz
			}
		}
	}
	return false
}

func xoreq(fr, tr, w int8) bool {
	switch fr {
	case 0:
		return tr == w
	case w:
		return tr == 0
	}
	return false
}

var xrqnmv = map[int8]map[int8]int8{
	6: {0: 1},
	7: {1: 1, 0: 2},
	0: {6: 1, 7: 2},
	1: {7: 1},
}

func techKnight(from, to Pos) bool { //cantech
	//gdziekolor := ColorUint8(uint8(from[1]>>3))
	//analiza wszystkich przypadkow ruchu przez moaty, gdzie wszystkie mozliwosci można wpisać ręcznie
	switch to[1] {
	case (from[1] + 2) % 24, (from[1] - 2 + 24) % 24:
		return abs(from[0]-to[0]) == 1
	case (from[1] + 1) % 24, (from[1] - 1 + 24) % 24:
		return abs(from[0]-to[0]) == 2
	case (from[1] + 1 + 12) % 24, (from[1] - 1 + 12) % 24:
		return from[0] == 5 && to[0] == 4 || from[0] == 4 && to[0] == 5
	case (from[1] + 2 + 12) % 24, (from[1] - 2 + 12) % 24:
		return from[0] == 5 && to[0] == 5
	}
	return false
}

//canmoatKnight returns cantech(canmoat)&cancap
func canmoatKnight(from, to Pos, m MoatsState) (bool, bool) {
	if !techKnight(from, to) {
		return false, false
	}
	if (from[0] < 3 || to[0] < 3) && xoreq(from[0], to[0], xrqnmv[from[1]%8][to[1]%8]) {
		return m[((from[1]+2)/8)%3], false
	}
	return true, true
}

func (b *Board) knightMove(from, to Pos, m MoatsState) bool {
	t, c := canmoatKnight(from, to, m)
	if !t {
		return false
	}
	dosq := b.GPos(to)
	return dosq.Empty() || c && dosq.Color() != b.GPos(from).Color()
}

func (b *Board) multithreatchecking(pa PlayersAlive, ep EnPassant, wheres ...Pos) (check Check) {
	for _, where := range wheres {
		check = b.ThreatChecking(where, pa, ep)
		if check.If {
			return
		}
	}
	return
}

func (b *Board) rank0files3thrchk(pa PlayersAlive, f, d int8) Check {
	return b.multithreatchecking(pa, DEFENPASSANT, Pos{0, f}, Pos{0, f + d}, Pos{0, f + d + d})
}

func (b *Board) castthreat(whopart int8, ksnotqs bool, pa PlayersAlive) Check {
	return b.rank0files3thrchk(pa, whopart*8+4, bdodatni(ksnotqs))
}

func (b *Board) castemptqside(whopart int8) bool {
	p := whopart * 8
	p4 := p + 4
	for i := p + 1; i < p4; i++ {
		if b[0][i].NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) castemptkside(whopart int8) bool {
	s := whopart*8 + 5
	if b[0][s].NotEmpty {
		return false
	}
	return b[0][s+1].Empty()
}

func (b *Board) castling(from Pos, to Pos, cs Castling, pa PlayersAlive) bool {
	if from[0] != 0 || to[0] != 0 || from[1]%8 != 4 {
		return false
	}
	switch to[1] {
	case from[1] - 2: //cuz queen is on the minus
		part := from[1] / 8
		return !(!cs[part][castnbyteQ] || b.castemptqside(part) || b.castthreat(part, false, pa).If)
	case from[1] + 2: //cuz king is on the plus
		part := from[1] / 8
		return !(!cs[part][castnbyteK] || b.castemptkside(part) || b.castthreat(part, true, pa).If)
	}
	return false
}

func (b *Board) rook(from Pos, to Pos, m MoatsState) bool { //whether a rook could move like that
	return b.straight(from, to, m)
}
func (b *Board) knight(from Pos, to Pos, m MoatsState) bool { //whether a knight could move like that
	return b.knightMove(from, to, m)
}
func (b *Board) bishop(from Pos, to Pos, m MoatsState) bool { //whether a boshop could move like that
	return b.diagonal(from, to, m)
}
func (b *Board) king(from Pos, to Pos, m MoatsState, cs Castling, pa PlayersAlive) bool { //whether a king could move like that
	return b.kingMove(from, to, m) || b.castling(from, to, cs, pa)
}
func (b *Board) queen(from Pos, to Pos, m MoatsState) bool { //whether a queen could move like that (concurrency, yay!)
	return b.straight(from, to, m) || b.diagonal(from, to, m)
	/*
		whether := make(chan bool)
		go func() { whether <- b.straight(from, to, m) }()
		go func() { whether <- b.diagonal(from, to, m) }()
		if <-whether {
			return true
		}
		return <-whether
	*/
}
func (b *Board) pawn(from Pos, to Pos, e EnPassant) bool { //whether a pawn could move like that
	var p PawnCenter
	p = (*b)[from[0]][from[1]].PawnCenter
	return b.pawnStraight(from, to, p) || b.pawnCapture(from, to, e, p)
}

//AnyPiece : tell whether the piece being in 'from' could move like that
func (b *Board) AnyPiece(from Pos, to Pos, m MoatsState, cs Castling, e EnPassant, pa PlayersAlive) bool {
	if err := from.Correct(); err != nil {
		panic(err)
	}
	if err := to.Correct(); err != nil {
		panic(err)
	}
	switch (*b)[from[0]][from[1]].What() {
	case Pawn:
		return b.pawn(from, to, e)
	case Rook:
		return b.rook(from, to, m)
	case Knight:
		return b.knight(from, to, m)
	case Bishop:
		return b.bishop(from, to, m)
	case King:
		return b.king(from, to, m, cs, pa)
	case Queen:
		return b.queen(from, to, m)
	default:
		if (*b)[from[0]][from[1]].NotEmpty {
			panic("What it is if it was said to exist???")
		} else {
			return false
		}
	}
}
