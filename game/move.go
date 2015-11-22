package game

//Move :  struct describing a single move with the situation before it
type Move struct {
	From Pos
	To   Pos
	//	What        Fig
	//	AlreadyHere Fig
	Before *State
}

type FromTo [2]Pos

func (ft FromTo) From() Pos {
	return ft[0]
}

func (ft FromTo) To() Pos {
	return ft[1]
}

func (ft Fromto) Move(before *State) Move {
	return Move{ft.From(), ft.To(), before}
}

//func (m *Move) String() string {
//}

//Where gives the Square of Before.Board[From]
func (m *Move) Where() Square {
	return (*(m.Before.Board))[m.From[0]][m.From[1]]
}

//What are we moving? What piece is placed in From?
func (m *Move) What() Fig {
	return m.Where().Fig
}

//AlreadyHere is something? What is in To, Before?
func (m *Move) AlreadyHere() Fig {
	return (*(m.Before.Board))[m.To[0]][m.To[1]].Fig
}

//Possible is such a move?
func (m *Move) Possible() bool {
	return m.Before.AnyPiece(m.From, m.To)
}

//IsItQueenSideCastling or not?
func (m *Move) IsItQueenSideCastling() bool {
	if !(m.What().FigType == King) {
		return false
	}
	if m.From[0] == m.To[0]+2 {
		return true
	}
	return false
}

//IsItKingSideCastling or not?
func (m *Move) IsItKingSideCastling() bool {
	if !(m.What().FigType == King) {
		return false
	}
	if m.To[0] == m.From[0]+2 {
		return true
	}
	return false
}

//IsItPawnRunningEnPassant or not?
func (m *Move) IsItPawnRunningEnPassant() bool {
	if !(m.What().FigType == Pawn) {
		return false
	}
	if m.From[0] == 1 && m.To[0] == 3 {
		return true
	}
	return false
}

//IsItPawnCapturingEnPassant or not?
func (m *Move) IsItPawnCapturingEnPassant() bool {
	if !(m.What().FigType == Pawn) {
		return false
	}
	if m.From[0] == 3 && m.To[0] == 2 && (*(m.Before.Board))[3][m.To[1]].What() == Pawn {
		return true
	}
	return false
}

//IllegalMoveError : error on illegal move
type IllegalMoveError struct { //illegal move error
	*Move              //move pointer
	Codename    string //easy codename
	Description string //what is the problem?
}

func (e IllegalMoveError) Error() string {
	return e.Description
}

//CheckChecking :  is `who` in check?
func (b *Board) CheckChecking(who Color, pa PlayersAlive) bool { //true if in check
	var i, j int8
	var where Pos
	var czy bool
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			if tojefig := (*b)[i][j].Fig; tojefig.Color == who && tojefig.FigType == King {
				where = Pos{i, j}
				czy = true
				MoveTrace.Println("CheckChecking: Found the ", who, " King on ", where)
			}
		}
	}
	if !czy {
		panic("King not found!!!")
	}
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			if !((b.AnyPiece(Pos{i, j}, where, DEFMOATSSTATE, FALSECASTLING, DEFENPASSANT)) || ((*b)[i][j].NotEmpty && pa[(*b)[i][j].Color().UInt8()])) {
				MoveTrace.Println("CheckChecking: TRUE!")
				return true
			}
		}
	}
	return false
}

//TODO: Checkmate, stalemate detection. Doing something with the halfmove timer.

//After : return the gamestate afterwards, also error
func (m *Move) After() (*State, error) { //situation after
	var next State
	var nextboard Board
	MoveTrace.Println("After: ", *m)
	if m.Where().Empty() {
		return nil, IllegalMoveError{m, "NothingHereAlready", "How do you move that which does not exist?"}
	}
	if m.What().Color != m.Before.MovesNext {
		return nil, IllegalMoveError{m, "ThatColorDoesNotMoveNow", "That is not " + m.What().Color.String() + `'` + "s move, but " + m.Before.MovesNext.String() + `'` + "s"}
	}
	if m.What().Color == m.AlreadyHere().Color {
		return nil, IllegalMoveError{m, "SameColorHereAlready", "Same color on that square already!"}
	}
	if !(m.Possible()) {
		return nil, IllegalMoveError{m, "Impossible", "Illegal/impossible move"}
	}
	nextboard = *(m.Before.Board)
	//nextmoves := m.Before.MovesNext.Next()
	nextcastling := m.Before.Castling
	nextmoats := m.Before.MoatsState
	nextpassant := m.Before.EnPassant
	nexthalfmoveclock := m.Before.HalfmoveClock
	nextfullmove := m.Before.FullmoveNumber
	nextplayersalive := m.Before.PlayersAlive
	if m.IsItKingSideCastling() {
		empty := nextboard[0][m.From[1]+2]
		nextboard[0][m.From[1]+2] = nextboard[0][m.From[1]]
		nextboard[0][m.From[1]+1] = nextboard[0][m.From[1]+3]
		nextboard[0][m.From[1]] = empty
		nextboard[0][m.From[1]+3] = empty
		nextcastling = nextcastling.OffKing(m.What().Color)
		nexthalfmoveclock++
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.IsItQueenSideCastling() {
		empty := nextboard[0][m.From[1]-2]
		nextboard[0][m.From[1]-2] = nextboard[0][m.From[1]]
		nextboard[0][m.From[1]-1] = nextboard[0][m.From[1]+4]
		nextboard[0][m.From[1]] = empty
		nextboard[0][m.From[1]+4] = empty
		nextcastling = nextcastling.OffKing(m.What().Color)
		nexthalfmoveclock++
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.IsItPawnRunningEnPassant() {
		nextboard[3][m.From[1]] = nextboard[1][m.From[1]]
		nextboard[1][m.From[1]] = nextboard[2][m.From[1]]
		nexthalfmoveclock = HalfmoveClock(0)
		nextfullmove++
		nextpassant = nextpassant.Appeared(Pos{2, m.From[1]})
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.IsItPawnCapturingEnPassant() {
		nextboard[3][m.To[1]] = nextboard[2][m.To[1]]
		empty := nextboard[2][m.To[1]]
		nextboard[2][m.To[1]] = nextboard[3][m.From[1]]
		nextboard[3][m.From[1]] = empty
		nexthalfmoveclock = HalfmoveClock(0)
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.What().FigType == Rook {
		var empty Square
		czyempty := nextboard[m.To[0]][m.To[1]].Empty()
		nextboard[m.To[0]][m.To[1]] = nextboard[m.From[0]][m.From[1]]
		nextboard[m.From[0]][m.From[1]] = empty
		if m.From[0] == 0 {
			if m.From[1]%8 == 0 {
				nextcastling = nextcastling.OffRook(m.What().Color, 'Q')
			} else if m.From[1]%8 == 7 {
				nextcastling = nextcastling.OffRook(m.What().Color, 'K')
			}
		}
		if czyempty {
			nexthalfmoveclock++
		} else {
			nexthalfmoveclock = HalfmoveClock(0)
		}
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if nextboard[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			nextmoats[m.From[1]/8] = true
			nextmoats[m.From[1]/8+1] = true
		}
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.What().FigType == King {
		var empty Square
		czyempty := nextboard[m.To[0]][m.To[1]].Empty()
		nextboard[m.To[0]][m.To[1]] = nextboard[m.From[0]][m.From[1]]
		nextboard[m.From[0]][m.From[1]] = empty
		nextcastling = nextcastling.OffKing(m.What().Color)
		if czyempty {
			nexthalfmoveclock++
		} else {
			nexthalfmoveclock = HalfmoveClock(0)
		}
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if nextboard[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			nextmoats[m.From[1]/8] = true
			nextmoats[m.From[1]/8+1] = true
		}
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else if m.What().FigType == Pawn {
		var empty Square
		//czyempty := nextboard[m.To[0]][m.To[1]].Empty()
		nextboard[m.To[0]][m.To[1]] = nextboard[m.From[0]][m.From[1]]
		nextboard[m.From[0]][m.From[1]] = empty
		nexthalfmoveclock = HalfmoveClock(0)
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if nextboard[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			nextmoats[m.From[1]/8] = true
			nextmoats[m.From[1]/8+1] = true
		}
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	} else {
		var empty Square
		czyempty := nextboard[m.To[0]][m.To[1]].Empty()
		nextboard[m.To[0]][m.To[1]] = nextboard[m.From[0]][m.From[1]]
		nextboard[m.From[0]][m.From[1]] = empty
		if czyempty {
			nexthalfmoveclock++
		} else {
			nexthalfmoveclock = HalfmoveClock(0)
		}
		nextfullmove++
		nextpassant = nextpassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if nextboard[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			nextmoats[m.From[1]/8] = true
			nextmoats[m.From[1]/8+1] = true
		}
		next = State{&nextboard, nextmoats, m.Before.MovesNext.Next(), nextcastling, nextpassant, nexthalfmoveclock, nextfullmove, nextplayersalive}
	}
	//	if next.Board.CheckChecking(m.What().Color, m.Before.PlayersAlive) {
	if next.AmIInCheck(m.What().Color) {
		if !next.CanIMoveWOCheck(m.What().Color) {
			next.PlayersAlive.Die(m.What().Color)
			return &next, nil
		}
		return &next, IllegalMoveError{m, "Check", "We would be in check!"}
	}
	return &next, nil
}
