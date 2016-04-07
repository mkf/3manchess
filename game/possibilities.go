package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

func (b *Board) canfigstraighthoriz(rank, from, to int8) (bool, bool) { //returns plus, minus
	return b.canfigstraighthorizdirec(rank, from, to, 1), b.canfigstraighthorizdirec(rank, from, to, -1)
}

func (b *Board) canfigstraighthorizdirec(r, f, t, d int8) bool { //rank, from file, to file, direction
	for i := f + d; d*((i-f)%24) < d*((t-f)%24); i = (i + d) % 24 {
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
	for i := f + s; s*i < t; i += s {
		if (*b)[i][file].NotEmpty {
			return false
		}
	}
	return true
}

func (b *Board) diagonal(from Pos, to Pos, m MoatsState) bool { //(bool, bool) {
	nasz := (*b)[from[0]][from[1]] //nasz Square
	if from == to {
		return false
	}

	przel := abs(to[1]-from[1]) % 24
	vectrank := to[0] - from[0]
	rankdirec := sign(vectrank)
	short := abs(vectrank) == przel     //without center
	long := abs(from[0]+to[0]) == przel //with center
	cantech := short || long

	var filedirec int8
	switch to[1] {
	case (from[1] + przel) % 24:
		filedirec = +1
	case (from[1] - przel) % 24:
		filedirec = -1
	default:
		if short && long { //ten warunek był zanegowany od 51d0219 do c984f32 włącznie, odnegowany w 5ff00460b17908610a8365f9e236519759869046
			panic(from.String() + " " + to.String())
		}
	}

	canfigshort, canfiglong, canmoatshort, canmoatlong, capcheckshort, capchecklong := true, true, true, true, true, true

	if from[0] == 0 || to[0] == 0 { //jeżeli jesteśmy na rank 0 i na rank 0 zmierzamy
		var sprawdzamy, mdir int8

		//if from[0] == 0 { // jeżeli wyjeżdżamy do środka
		//mdir := filedirec     // short jedzie w kierunku mdir, long jedzie w -mdir
		//sprawdzamy := from[1] //
		//} else { // czyli inaczej:  else if to[0]==0  czyli  jeżeli jedziemy na brzeg
		//mdir := -filedirec
		//sprawdzamy := to[1]
		//}

		var dirtemp int8     // jak pojedziesz w tę stronę to wjedziesz w moat'a, jak nie wjedziesz to chyba zero
		capchecktemp := true // czy wjedziesz w moat'a jak pojedziesz w stronę `dirtemp`, redundantne
		canmoattemp := true  // czy ten moat jest bridged czy nie

		//switch sprawdzamy { case 0, 23: canmoattemp = m[0];  case 8, 7: m[1];  case 16, 15: m[2] }
		canmoattemp = m[((sprawdzamy+1)/24)%3]
		switch sprawdzamy % 8 {
		case 0: //cross jest jak pojedziemy na minus
			dirtemp = -1
			capchecktemp = false
		case 7: //cross jest jak pojedziemy na plus
			dirtemp = 1
			capchecktemp = false
		}
		if dirtemp == mdir {
			capcheckshort = capchecktemp
			canmoatshort = canmoattemp
		} else if dirtemp == -mdir {
			capchecklong = capchecktemp
			canmoatlong = canmoattemp
		}
	}
	var bijemyostatniego bool
	if short && canmoatshort {
		var i int8
		for i = 1; i < przel; i++ {
			if (*b)[from[0]+(i*rankdirec)][(((from[1]+(i*filedirec))%24)+24)%24].NotEmpty {
				canfigshort = false
				break
			}
		}
		ostatni := (*b)[from[0]+(przel*rankdirec)][(((from[1]+(przel*filedirec))%24)+24)%24]
		if ostatni.NotEmpty {
			bijemyostatniego = ostatni.Color() != nasz.Color()
			if !bijemyostatniego {
				canfigshort = false
			}
		}
	}
	if long && canmoatlong {
		var i int8
		if canfiglong {
			for i = 1; i <= (5 - from[0]); i++ {
				if (*b)[from[0]+i][(((from[1]+(i*filedirec))%24)+24)%24].NotEmpty {
					canfiglong = false
					break
				}
			}
		}
		if canfiglong {
			for i = 0; i+5-from[0] < przel; i++ {
				if (*b)[5-i][(((from[1]+((5-from[0]+i)*filedirec))%24)+24)%24].NotEmpty {
					canfiglong = false
					break
				}
			}
		}
		ostatni := (*b)[(((10-from[0]-przel)%6)+6)%6][(((from[1]+(przel*filedirec))%24)+24)%24]
		if r := recover(); r != nil {
			panic(from[0] + przel)
		}
		if ostatni.NotEmpty {
			if ostatni.Color() != nasz.Color() {
				bijemyostatniego = true
			} else {
				canfiglong = false
			}
		}
	}
	canshort := cantech && canfigshort && canmoatshort && (capcheckshort || !bijemyostatniego)
	canlong := cantech && canfiglong && canmoatlong && (capchecklong || !bijemyostatniego)
	return canshort || canlong
}

func (b *Board) pawnStraight(from Pos, to Pos, p PawnCenter) bool { //(bool,PawnCenter,EnPassant) {
	var cantech, canfig bool
	//pc := p
	//ep := e
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
		realsgn := sign(to[0] - from[0])
		if realsgn != sgn {
			return false //,p,e
		}
		if !bool(p) && from[0] == 1 && to[0] == 3 {
			cantech = true
			canfig = (*b)[2][from[1]].Empty() && b.GPos(to).Empty()
			//ep:=e.Appeared(Pos{2,from[1]})
		} else if to[0] == from[0]+sgn {
			cantech = true
			canfig = b.GPos(to).Empty()
		}
	} else if ((from[1]+12)%24) == to[1] && from[0] == 5 && to[0] == 5 && !bool(p) {
		cantech = true
		canfig = b.GPos(to).Empty()
		//pc = true
	}
	return cantech && canfig //, pc, ep
}

func (b *Board) kingStraight(from Pos, to Pos, m MoatsState) bool {
	if from == to {
		return false
	}
	nasz := b.GPos(from)
	tjf := b.GPos(to)
	switch to {
	case Pos{from[0] + 1, from[1]},
		Pos{from[0] - 1, from[1]},
		Pos{from[0], (from[1] + 1) % 24},
		Pos{from[0], (from[1] - 1) % 24}:
		return !(tjf.NotEmpty && tjf.Color() == nasz.Color())
	}
	return false
}

func (b *Board) pawnCapture(from Pos, to Pos, e EnPassant, p PawnCenter) bool {
	nasz := b.GPos(from)
	gdziekolor := Color(from[1]>>3 + 1)
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
	if nasz.Color() == gdziekolor && !p {
		return false //panic("pC" + nasz.Color().String())
	}
	var sgn int8
	if p {
		sgn = int8(-1)
	} else {
		sgn = int8(1)
	}
	return ((from[0] == 5 && to[0] == 5 && !bool(p)) && //jest na 5 ranku i nie przeszedl przez srodek jeszcze
		(to[1] == ((from[1]+24-10)%24) || to[1] == ((from[1]+10)%24)) && //poprawnie przelecial na skos przez srodek
		b.GPos(to).Color() != nasz.Color()) || //ten co go bijemy jest innego koloru ALBO
		((e[0] == to || e[1] == to) && //pozycja tego co go bijemy jest w enpassant
			(*b)[3][to[1]].What() == Pawn && //ten co go bijemy jest pionkiem
			(*b)[3][to[1]].Color() != nasz.Color() && //i jest innego koloru
			(*b)[2][to[1]].Empty()) || //a pole za nim jest puste (jak to po ruchu pre-enpassant) ALBO
		(to[0] == from[0]+sgn && cancreek && //zwykle bicie, o jeden w kierunku sgn na ranku
			((to[1] == (from[1]+1)%24) || (to[1] == (from[1]+24-1)%24)) && //i o jeden w tę lub tamtą stronę (wsio mod24) na file'u
			b.GPos(to).Color() != nasz.Color()) //a ten co go bijemy jest innego koloru
}

func (b *Board) knightMove(from Pos, to Pos, m MoatsState) bool {
	nasz := (*b)[from[0]][from[1]]
	//gdziekolor := ColorUint8(uint8(from[1]>>3))
	//analiza wszystkich przypadkow ruchu przez moaty, gdzie wszystkie mozliwosci można wpisać ręcznie
	cantech := false
	switch to[1] {
	case (from[1] + 2) % 24, (from[1] - 2) % 24:
		if from[0] == 5 && to[0] == 5 {
			cantech = true
		} else if abs(from[0]-to[0]) == 1 {
			cantech = true
		}
	case (from[1] + 1) % 24, (from[1] - 1) % 24:
		if from[0] == 5 && to[0] == 4 { // doubtful, awaiting email reply
			cantech = true
		} else if abs(from[0]-to[0]) == 2 {
			cantech = true
		}
	}
	canmoat := true
	//cancheck := true
	if cantech && from[0] < 3 && to[0] < 3 {
		var ourmoat bool
		switch from[1] {
		case 22, 23, 0, 1:
			ourmoat = m[0]
		case 6, 7, 8, 9:
			ourmoat = m[1]
		case 14, 15, 16, 17:
			ourmoat = m[2]
		}
		switch from[1] % 8 {
		case 6:
			if to[1]%8 == 0 {
				if from[0]+to[0] == 1 {
					canmoat = ourmoat
				}
			}
		case 7:
			if to[1]%8 == 1 {
				if from[0]+to[0] == 1 {
					canmoat = ourmoat
				}
			} else if to[1]%8 == 0 {
				switch from[0] {
				case 0:
					if to[0] == 2 { //cancheck = false
						canmoat = ourmoat
					}
				case 2:
					if to[0] == 0 { //cancheck = false
						canmoat = ourmoat
					}
				}
			}
		case 0:
			if to[1]%8 == 6 {
				if from[0]+to[0] == 1 {
					canmoat = ourmoat
				}
			} else if to[1]%8 == 7 {
				switch from[0] {
				case 0:
					if to[0] == 2 { //cancheck = false
						canmoat = ourmoat
					}
				case 2:
					if to[0] == 0 { //cancheck = false
						canmoat = ourmoat
					}
				}
			}
		case 1:
			if to[1]%8 == 7 {
				switch from[0] {
				case 0:
					if to[0] == 1 { //cancheck = false
						canmoat = ourmoat
					}
				case 1:
					if to[0] == 0 { //cancheck = false
						canmoat = ourmoat
					}
				}
			}
		}
	}
	canfig := true
	if cantech && canmoat {
		if (*b)[to[0]][to[1]].NotEmpty {
			if (*b)[to[0]][to[1]].Color() == nasz.Color() {
				canfig = false
			}
		}
	}
	return cantech && canmoat && canfig
}

func (b *Board) castling(from Pos, to Pos, cs Castling) bool {
	var colorproper bool
	var col Color
	switch from {
	case Pos{4, 0}: //white king starting
		col = White
		colorproper = true
	case Pos{12, 0}: //gray king starting
		col = Gray
		colorproper = true
	case Pos{20, 0}: //black king starting
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
	return (kingside && (*b)[0][from[1]+1].Empty() && (*b)[0][from[1]+2].Empty()) || //kingside and kingside empty
		(queenside && (*b)[0][to[1]+1].Empty() && (*b)[0][to[1]+2].Empty() && (*b)[0][to[1]+3].Empty())
	//		quenside and queenside empty
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
func (b *Board) king(from Pos, to Pos, m MoatsState, cs Castling) bool { //whether a king could move like that
	return b.kingStraight(from, to, m) || b.castling(from, to, cs)
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
func (b *Board) AnyPiece(from Pos, to Pos, m MoatsState, cs Castling, e EnPassant) bool {
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
		return b.king(from, to, m, cs)
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
