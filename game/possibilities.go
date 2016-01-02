package game

import "log"

func (b *Board) straight(from Pos, to Pos, m MoatsState) bool { //(bool, bool) { //(whether it can, whether it can capture/check)
	var cantech, canmoat, canfig bool
	//capcheck := true
	if from == to { //Same square
		//panic("Same square!")
		return false
	}
	if from[0] == to[0] { //same rank
		cantech = true
		if from[0] == 0 { //first rank
			var mshort, mlong bool //, capcheckshort bool
			var direcshort int8
			var fromtominus int8
			if from[1]/8 == to[1]/8 { //same color area
				//capcheckshort = true
				canmoat = true
				mshort = true
				if m[0] && m[1] && m[2] {
					mlong = true
				}
				direcshort = sign(to[1] - from[1])
			} else { //moving to another color's area
				//capcheckshort = false
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
				fromtominus = fromto[1] - fromto[0]
				if abs(fromtominus) == 2 {
					fromtominus = -fromtominus
				}
				direcshort = sign(fromtominus)
			}
			canfigminus := true
			//straight in +file direction, mod24 ofcoz
			for i := from[1] + 1; ((i-from[1])%24 < (to[1]-from[1])%24) && canfig; i = (i + 1) % 24 {
				//if something between A and B
				if !((*b)[0][i].Empty()) {
					canfig = false
				}
			}
			//straight in -file direction, mod24 ofcoz
			for i := from[1] - 1; ((i-from[1])%24 > (to[1]-from[1])%24) && canfigminus; i = (i - 1) % 24 {
				//if something is between A and B
				if !((*b)[0][i].Empty()) {
					canfigminus = false
				}
			}
			canfigplus := canfig //legacy mess, but not much
			canfig = canfigplus || canfigminus

			//as we are on the first rank && moving to another color's area, we gotta check the moats
			if direcshort == 1 {
				if canfigplus && mshort {
					canmoat = true
				} else if canfigminus && mlong {
					canmoat = true
				}
			} else if direcshort == -1 {
				if canfigminus && mshort {
					canmoat = true
				} else if canfigplus && mlong {
					canmoat = true
				}
			} else {
				//panic(direcshort)
				return false
			}
		} else { //if same rank, but not first rank
			canmoat = true
			canfigplus := true
			//straight direc +file (mod24 ofcoz)
			for i := from[1] + 1; ((i-from[1])%24 < (to[1]-from[1])%24) && canfigplus; i = (i + 1) % 24 {
				if !((*b)[from[0]][i].Empty()) {
					canfigplus = false
				}
			}
			canfigminus := true
			//straight direc -file (mod24 ofcoz)
			for i := from[1] - 1; ((i-from[1])%24 > (to[1]-from[1])%24) && canfigminus; i = (i - 1) % 24 {
				if !((*b)[from[0]][i].Empty()) {
					canfigminus = false
				}
			}
			canfig = canfigplus || canfigminus
		}
	} else if from[1] == to[1] { //if the same file, ie. no passing through center
		cantech = true
		canmoat = true
		canfig = true
		sgn := sign(to[0] - from[0])
		for i := from[0] + sgn; (sgn*i < to[0]) && canfig; i += sgn {
			if !((*b)[i][from[1]].Empty()) {
				canfig = false
			}
		}
	} else if ((from[1] - 12) % 24) == to[1] { //if the adjacent file, passing through center
		cantech = true
		canmoat = true
		canfig = true
		//searching for collisions from both sides of the center
		for i, j := from[0], to[0]; canfig && (i < 6 && j < 6); i, j = i+1, j+1 {
			if !((*b)[i][from[1]].Empty()) {
				canfig = false
			}
			if !((*b)[j][to[1]].Empty()) {
				canfig = false
			}
		}
	} else { //not the same rank and not the same file nor adjacent
		cantech = false
	}
	final := cantech && canmoat && canfig
	return final //, capcheck && final
}

func (b *Board) diagonal(from Pos, to Pos, m MoatsState) bool { //(bool, bool) {
	nasz := (*b)[from[0]][from[1]] //nasz Square
	//if from[0] != 0 && from == to {
	//panic("Same square and not the first rank!")
	//}
	if from == to {
		//panic("Same square!") //make sure //awaiting email reply
		//if such thing was legal, one could easily escape a zugzwang, having such a possibility around
		//also, that would make move detection *really* hard
		return false
	}

	przel := abs(to[1]-from[1]) % 24
	vectrank := to[0] - from[0]
	rankdirec := sign(vectrank)
	short := abs(vectrank) == przel     //without center
	long := abs(from[0]+to[0]) == przel //with center
	cantech := short || long

	var filedirec int8
	if (from[1]+przel)%24 == to[1] {
		filedirec = +1
	} else if (from[1]-przel)%24 == to[1] {
		filedirec = -1
	} else if !(short && long) {
		panic(from.String() + " " + to.String())
	}

	//var canfig bool
	canfigshort := true
	canfiglong := true
	canmoatshort := true
	canmoatlong := true
	capcheckshort := true
	capchecklong := true

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

		switch sprawdzamy {
		case 0: //cross jest jak pojedziemy na minus
			dirtemp = -1
			capchecktemp = false
			canmoattemp = m[0]
		case 23: //cross jest jak pojedziemy na plus
			dirtemp = 1
			capchecktemp = false
			canmoattemp = m[0]
		case 8:
			dirtemp = -1
			capchecktemp = false
			canmoattemp = m[1]
		case 7:
			dirtemp = 1
			capchecktemp = false
			canmoattemp = m[1]
		case 16:
			dirtemp = -1
			capchecktemp = false
			canmoattemp = m[2]
		case 15:
			dirtemp = 1
			capchecktemp = false
			canmoattemp = m[2]
		}
		if dirtemp == mdir {
			capcheckshort = capchecktemp
			canmoatshort = canmoattemp
		} else if dirtemp == -mdir {
			capchecklong = capchecktemp
			canmoatlong = canmoattemp
		}
	}
	bijemyostatniego := false
	if short && canmoatshort {
		var i int8
		for i = 1; i < przel; i++ {
			if (*b)[from[0]+(i*rankdirec)][(from[1]+(i*filedirec))%24].NotEmpty {
				canfigshort = false
				break
			}
		}
		ostatni := (*b)[from[0]+(przel*rankdirec)][(from[1]+(przel*filedirec))%24]
		if ostatni.NotEmpty {
			if ostatni.Color() != nasz.Color() {
				bijemyostatniego = true
			} else {
				canfigshort = false
			}
		}
	}
	if long && canmoatlong {
		var i int8
		for i = 1; i <= (5 - from[0]); i++ {
			if (*b)[from[0]+i][(from[1]+(i*filedirec))%24].NotEmpty {
				canfiglong = false
				break
			}
			if recerr := recover(); recerr != nil {
				log.Println(*b)
				log.Println(from, to)
				log.Println("i", i)
				log.Println("fd", filedirec)
				panic(recerr)
			}
		}
		for i = 0; i+5-from[0] < przel; i++ {
			if (*b)[5-i][(from[1]+((5-from[0]+i)*filedirec))%24].NotEmpty {
				canfiglong = false
				break
			}
		}
		ostatni := (*b)[(10-from[0]-przel)%6][(from[1]+(przel*filedirec))%24]
		if r := recover(); r != nil {
			panic(from[0] + przel)
		}
		if !(ostatni.Empty()) {
			if ostatni.Color() != nasz.Color() {
				bijemyostatniego = true
			} else {
				canfiglong = false
			}
		}
	}
	//canfig = canfigshort || canfiglong
	//if canfigshort && canfiglong {
	//	canmoat = canmoatshort || canmoatlong
	//} // dalej: co jeśli jedno z nich? rozpatrywać przypadki tylko short i tylko long
	canshort := cantech && canfigshort && canmoatshort && (capcheckshort || !bijemyostatniego)
	canlong := cantech && canfiglong && canmoatlong && (capchecklong || !bijemyostatniego)
	/*	if canshort && canlong { capcheck = capcheckshort || capchecklong
		} else if canshort {     capcheck = capcheckshort
		} else if canlong {      capcheck = capchecklong                 }  */
	return canshort || canlong //	, capcheck
}

func (b *Board) pawnStraight(from Pos, to Pos, p PawnCenter) bool { //(bool,PawnCenter,EnPassant) {
	var cantech, canfig bool
	//pc := p
	//ep := e
	if from == to {
		//panic("Same square!")
		return false
	}
	nasz := (*b)[from[0]][from[1]]
	gdziekolor := ColorUint8(uint8(from[1] / 8))
	if nasz.Color() == gdziekolor && p {
		panic("pS" + nasz.Color().String())
	}
	var sgn int8
	if p {
		sgn = int8(-1)
	} else {
		sgn = int8(1)
	}
	if from[1] == to[1] {
		realsgn := sign(to[0] - from[0])
		if realsgn != sgn {
			return false //,p,e
		}
		if !p && from[0] == 1 && to[0] == 3 {
			cantech = true
			canfig = (*b)[2][from[1]].Empty() && (*b)[3][from[1]].Empty()
			//ep:=e.Appeared(Pos{2,from[1]})
		} else if to[0] == from[0]+sgn {
			cantech = true
			canfig = (*b)[to[0]][from[1]].Empty()
		}
	} else if ((from[1]-12)%24) == to[1] && from[0] == 5 && to[0] == 5 && !bool(p) {
		cantech = true
		canfig = (*b)[5][to[1]].Empty()
		//pc = true
	}
	return cantech && canfig //, pc, ep
}

func (b *Board) kingStraight(from Pos, to Pos, m MoatsState) bool {
	if from == to {
		return false
	}
	nasz := (*b)[from[0]][from[1]]
	switch to {
	case Pos{from[0] + 1, from[1]}, Pos{from[0] - 1, from[1]}, Pos{from[0], (from[1] + 1) % 24}, Pos{from[0], (from[1] - 1) % 24}:
		if (*b)[to[0]][to[1]].NotEmpty {
			if (*b)[to[0]][to[1]].Color() == nasz.Color() {
				return false
			}
			return true
		}
		return true
	default:
		return false
	}
}

func (b *Board) pawnCapture(from Pos, to Pos, e EnPassant, p PawnCenter) bool {
	nasz := (*b)[from[0]][from[1]]
	gdziekolor := ColorUint8(uint8(from[1] / 8))
	//cancreek := true
	if from == to {
		return false
	}
	if !p {
		creektemp1 := false
		fromto := [2]int8{from[0], to[0]}
		switch fromto {
		case [2]int8{0, 1}, [2]int8{1, 0}, [2]int8{1, 2}, [2]int8{2, 1}, [2]int8{2, 3}, [2]int8{3, 2}:
			creektemp1 = true
		}
		if ((to[1]%8 == 0 && from[1]%8 == 7) || (from[1]%8 == 0 && to[1]%8 == 7)) && creektemp1 {
			//cancreek = false
		}
	}
	if nasz.Color() == gdziekolor && !p {
		return false
		//panic("pC" + nasz.Color().String())
	}
	var sgn int8
	if p {
		sgn = int8(-1)
	} else {
		sgn = int8(1)
	}
	if from[0] == 5 && !bool(p) && to[0] == 5 && (to[1] == ((from[1]-10)%24) || to[1] == ((from[1]+10)%24)) && (*b)[to[0]][to[1]].Color() != nasz.Color() {
		return true
	}
	if (e[0] == to || e[1] == to) && (*b)[3][to[1]].What() == Pawn && (*b)[3][to[1]].Color() != nasz.Color() && (*b)[2][to[1]].Empty() {
		return true
	} else if to[0] == from[0]+sgn && ((to[1] == (from[1]+1)%24) || (to[1] == (from[1]-1)%24)) && (*b)[to[0]][to[1]].Color() != nasz.Color() {
		return true
	}
	return false
}

func (b *Board) knightMove(from Pos, to Pos, m MoatsState) bool {
	nasz := (*b)[from[0]][from[1]]
	//gdziekolor := ColorUint8(uint8(from[1] / 8))
	//analiza wszystkich przypadkow ruchu przez moaty, gdzie wszystkie mozliwosci można wpisać ręcznie
	cantech := false
	switch to[1] {
	case (from[1] + 2) % 24, (from[1] - 2) % 24:
		if from[0] == 5 && to[0] == 5 {
			cantech = true
		} else if from[0] == to[0]+1 || from[0] == to[0]-1 {
			cantech = true
		}
	case (from[1] + 1) % 24, (from[1] - 1) % 24:
		if from[0] == 5 && to[0] == 4 { // doubtful, awaiting email reply
			cantech = true
		} else if from[0] == to[0]+2 || from[0] == to[0]-2 {
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
		case 7:
			if to[1]%8 == 1 {
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
	whether := make(chan bool)
	go func() { whether <- b.straight(from, to, m) }()
	go func() { whether <- b.diagonal(from, to, m) }()
	if <-whether {
		return true
	}
	return <-whether
}
func (b *Board) pawn(from Pos, to Pos, e EnPassant) bool { //whether a pawn could move like that
	var p PawnCenter
	p = (*b)[from[0]][from[1]].PawnCenter
	return b.pawnStraight(from, to, p) || b.pawnCapture(from, to, e, p)
}

//AnyPiece : tell whether the piece being in 'from' could move like that
func (b *Board) AnyPiece(from Pos, to Pos, m MoatsState, cs Castling, e EnPassant) bool {
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
