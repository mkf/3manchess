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

		} else { //if same rank, but not first rank
			canmoat = true
			canfigplus, canfigminus := b.canfigstraighthoriz(from[0], from[1], to[1])
			canfig = canfigplus || canfigminus
		}
	} else if from[1] == to[1] { //if the same file, ie. no passing through center
		cantech, canmoat = true, true
		canfig = b.canfigstraightvertnormal(from[1], from[0], to[0])
	} else if ((from[1] - 12) % 24) == to[1] { //if the adjacent file, passing through center
		cantech, canmoat = true, true
		canfig = b.canfigstraightvertthrucenter(to[1], from[0], to[0])
	}
	return cantech && canmoat && canfig
}

func (b *Board) canfigstraightvertthrucenter(s, f, t int8) bool { //startfile (from[0]), from, to
	e := (uint8(s) - 12) % 24
	//searching for collisions from both sides of the center
	for i := f; i < 6; i++ {
		if (*b)[i][s].NotEmpty {
			return false
		}
	}
	for i := t; i < 6; i++ {
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
	return Pos{p[0] + v[0], (((p[1] + v[0]) % 24) + 24) % 24}
}

func (p Pos) MinusVector(v [2]int8) Pos {
	return p.AddVector([2]int8{-v[0], -v[1]})
}

func (b *Board) diagonal(from Pos, to Pos, m MoatsState) bool {
	for _, modifyPos := range PLUSMINUSPAIRS {
		for pos := from.AddVector(modifyPos); pos[0] >= 0; pos = pos.AddVector(modifyPos) {
			if max(from[0], to[0]) == 1 && abs(from[1]%8-to[1]%8) == 7 {
				// we crossed a moat (moving diagonally)
				var i int8 // moat index
				if from[1]%8 == 7 {
					i = from[1] / 8
				} else {
					i = to[1] / 8
				}
				if !m[i] {
					break
				}
			}
			if pos[0] > 5 { // we are crossing the center
				modifyPos[0] = -1
				pos = Pos{
					5,
					(pos[1] + (-1+10)*modifyPos[1] + 24) % 24,
				}
			}
			if pos == to {
				return true
			}
			if b.GPos(pos).NotEmpty {
				break
			}
		}
	}
	return false
}

func (b *Board) pawnStraight(from Pos, to Pos, p PawnCenter) bool { //(bool,PawnCenter,EnPassant) {
	if from == to {
		//panic("Same square!")
		return false
	}
	nasz := b.GPos(from)
	gdziekolor := Color(from[1]>>3 + 1)
	if nasz.Color() == gdziekolor && p {
		panic("pS" + nasz.Color().String())
	}
	var sgn int8
	if p {
		sgn = -1
	} else {
		sgn = 1
	}
	if from[1] == to[1] {
		if sign(to[0]-from[0]) != sgn {
			return false //,p,e
		}
		if !bool(p) && from[0] == 1 && to[0] == 3 {
			return (*b)[2][from[1]].Empty() && b.GPos(to).Empty()
			//ep:=e.Appeared(Pos{2,from[1]})
		} else if to[0] == from[0]+sgn {
			return b.GPos(to).Empty()
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

func (b *Board) pawnCapture(from Pos, to Pos, e EnPassant, p PawnCenter) bool {
	nasz := b.GPos(from)
	cancreek := true
	if from == to {
		return false
	}
	if !p {
		creektemp := false                //normalnie nie musi się męczyć z creekami
		fromto := [2]int8{from[0], to[0]} //przemieszczenie między rankami
		switch fromto {
		case [2]int8{0, 1}, [2]int8{1, 0}, [2]int8{1, 2}, [2]int8{2, 1}, [2]int8{2, 3}, [2]int8{3, 2}:
			creektemp = true //otrze się o creek
		}
		cancreek = !(creektemp && ((to[1]%8 == 0 && from[1]%8 == 7) || (from[1]%8 == 0 && to[1]%8 == 7)))
		//założenie: creeki są aktywne nawet jak już nie ma moatów
	}
	var sgn int8
	if p {
		sgn = int8(-1)
	} else {
		sgn = int8(1)
	}
	var ep Pos
	if e[0] == to {
		ep = e[0]
	} else {
		ep = e[1]
	}
	return ((from[0] == 5 && to[0] == 5 && !bool(p)) && //jest na 5 ranku i nie przeszedl przez srodek jeszcze
		(to[1] == ((from[1]+24-10)%24) || to[1] == ((from[1]+10)%24)) && //poprawnie przelecial na skos przez srodek
		b.GPos(to).NotEmpty && b.GPos(to).Fig.Color != nasz.Fig.Color) || //ten co go bijemy jest innego koloru ALBO
		(to[0] == from[0]+sgn && cancreek && //zwykle bicie, o jeden w kierunku sgn na ranku
			((to[1] == (from[1]+1)%24) || (to[1] == (from[1]+24-1)%24)) && //o jeden w tę lub tamtą stronę (wsio mod24) na file'u
			(((e[0] == to || e[1] == to) && //pozycja tego co go bijemy jest w enpassant
				(*b)[3][ep[1]].Fig.FigType == Pawn && //ten co go bijemy jest pionkiem
				(*b)[3][ep[1]].NotEmpty && (*b)[3][ep[1]].Fig.Color != nasz.Fig.Color && //i jest innego koloru
				(*b)[2][ep[1]].Empty()) || //a pole za nim jest puste (jak to po ruchu pre-enpassant) ALBO
				(b.GPos(to).NotEmpty && b.GPos(to).Fig.Color != nasz.Fig.Color))) //ten co go bijemy jest innego koloru
}

func xoreq(n1, n2, w int8) bool {
	switch n1 {
	case 0:
		return n2 == w
	case w:
		return n2 == 0
	}
	return false
}

var xrqnmv = map[int8]map[int8]int8{
	6: {0: 1},
	7: {1: 1, 0: 2},
	0: {6: 1, 7: 2},
	1: {7: 1},
}

func (b *Board) knightMove(from Pos, to Pos, m MoatsState) bool {
	//gdziekolor := ColorUint8(uint8(from[1]>>3))
	//analiza wszystkich przypadkow ruchu przez moaty, gdzie wszystkie mozliwosci można wpisać ręcznie
	var can bool
	switch to[1] {
	case (from[1] + 2) % 24, (from[1] - 2 + 24) % 24:
		can = from[0] == 5 && to[0] == 5 || abs(from[0]-to[0]) == 1
	case (from[1] + 1) % 24, (from[1] - 1 + 24) % 24:
		can = from[0] == 5 && to[0] == 4 || abs(from[0]-to[0]) == 2
	}
	if !can {
		return false
	}
	can = from[0] >= 3 && to[0] >= 3
	if !can {
		switch from[1] {
		case 22, 23, 0, 1:
			can = m[0]
		case 6, 7, 8, 9:
			can = m[1]
		case 14, 15, 16, 17:
			can = m[2]
		}
		if !can && xoreq(from[0], to[0], xrqnmv[from[1]%8][to[1]%8]) {
			return false
		}
	}
	dosq := b.GPos(to)
	return dosq.Empty() || dosq.Color() != b.GPos(from).Color()
}

func (b *Board) castling(from Pos, to Pos, cs Castling, pa PlayersAlive) bool {
	var colorproper bool
	var col Color
	switch from {
	case Pos{0, 4}: //white king starting
		col = White
		colorproper = true
	case Pos{0, 12}: //gray king starting
		col = Gray
		colorproper = true
	case Pos{0, 20}: //black king starting
		col = Black
		colorproper = true
	}
	if !colorproper || to[0] != 0 {
		return false
	}
	queenside := false
	kingside := false
	switch to[1] {
	case from[1] - 2: //cuz queen is on the minus
		queenside = cs.Give(col, 'Q')
	case from[1] + 2: //cuz king is on the plus
		kingside = cs.Give(col, 'K')
	}
	if kingside && (*b)[0][from[1]+1].Empty() && (*b)[0][from[1]+2].Empty() { // kingside and kingside empty
		uncheckedPos := [3]Pos{{from[0], from[0]}, {to[0] + 1, to[1] + 1}, {to[0], to[1]}}
		for _, checkPos := range uncheckedPos {
			check := b.ThreatChecking(checkPos, pa, DEFENPASSANT)
			if check.If == true {
				kingside = false
				break
			}
		}
	}
	if !kingside && queenside && (*b)[0][to[1]+1].Empty() && (*b)[0][to[1]+2].Empty() && (*b)[0][to[1]+3].Empty() { // not kingside, queenside and queenside empty
		uncheckedPos := [3]Pos{{from[0], from[0]}, {to[0] - 1, to[1] - 1}, {to[0], to[1]}}
		for _, checkPos := range uncheckedPos {
			check := b.ThreatChecking(checkPos, pa, DEFENPASSANT)
			if check.If == true {
				kingside = false
				break
			}
		}
	}
	return kingside || queenside
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
