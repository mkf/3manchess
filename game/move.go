package game

//Move :  struct describing a single move with the situation before it
type Move struct {
	From Pos
	To   Pos
	//	What        Fig
	//	AlreadyHere Fig
	Before *State
}

//FromTo is a type useful for AI and tests
type FromTo [2]Pos

//From gives you the From field
func (ft FromTo) From() Pos {
	return ft[0]
}

//To gives you the To field
func (ft FromTo) To() Pos {
	return ft[1]
}

//Move gives you a Move with the given Before *State
func (ft FromTo) Move(before *State) Move {
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

//PiecePossible is such a move?
func (m *Move) PiecePossible() bool {
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
	var ourpos Pos
	for i = 0; i < 6; i++ {
		for j = 0; j < 24; j++ {
			ourpos = Pos{i, j}
			if !((b.AnyPiece(ourpos, where, DEFMOATSSTATE, FALSECASTLING, DEFENPASSANT)) || ((*b)[i][j].NotEmpty && pa[(*b)[i][j].Color().UInt8()])) {
				MoveTrace.Println("CheckChecking: TRUE!", ourpos, (*b)[i][j])
				return true
			}
		}
	}
	return false
}

//TODO: Checkmate, stalemate detection. Doing something with the halfmove timer.

//Possible is such a move? Returns an error, same error as After() would give you, ¡¡¡except for CheckChecking!!!
func (m *Move) Possible() error {
	if m.Where().Empty() {
		return IllegalMoveError{m, "NothingHereAlready", "How do you move that which does not exist?"}
	}
	if m.What().Color != m.Before.MovesNext {
		return IllegalMoveError{m, "ThatColorDoesNotMoveNow", "That is not " + m.What().Color.String() + `'` + "s move, but " + m.Before.MovesNext.String() + `'` + "s"}
	}
	if m.What().Color == m.AlreadyHere().Color {
		return IllegalMoveError{m, "SameColorHereAlready", "Same color on that square already!"}
	}
	if !(m.PiecePossible()) {
		return IllegalMoveError{m, "Impossible", "Illegal/impossible move"}
	}
	return nil
}

//After : return the gamestate afterwards, also error
func (m *Move) After() (*State, error) { //situation after
	MoveTrace.Println("After: ", m.From, m.To)
	if merr := m.Possible(); merr != nil {
		return nil, merr
	}
	next := *m.Before
	nextboard := *m.Before.Board
	next.Board = &nextboard
	next.MovesNext = next.MovesNext.Next()

	if m.IsItKingSideCastling() {
		empty := next.Board[0][m.From[1]+2]
		next.Board[0][m.From[1]+2] = next.Board[0][m.From[1]]
		next.Board[0][m.From[1]+1] = next.Board[0][m.From[1]+3]
		next.Board[0][m.From[1]] = empty
		next.Board[0][m.From[1]+3] = empty
		next.Castling = next.Castling.OffKing(m.What().Color)
		next.HalfmoveClock++
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.IsItQueenSideCastling() {
		empty := next.Board[0][m.From[1]-2]
		next.Board[0][m.From[1]-2] = next.Board[0][m.From[1]]
		next.Board[0][m.From[1]-1] = next.Board[0][m.From[1]+4]
		next.Board[0][m.From[1]] = empty
		next.Board[0][m.From[1]+4] = empty
		next.Castling = next.Castling.OffKing(m.What().Color)
		next.HalfmoveClock++
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.IsItPawnRunningEnPassant() {
		next.Board[3][m.From[1]] = next.Board[1][m.From[1]]
		next.Board[1][m.From[1]] = next.Board[2][m.From[1]]
		next.HalfmoveClock = HalfmoveClock(0)
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Appeared(Pos{2, m.From[1]})
	} else if m.IsItPawnCapturingEnPassant() {
		next.Board[3][m.To[1]] = next.Board[2][m.To[1]]
		empty := next.Board[2][m.To[1]]
		next.Board[2][m.To[1]] = next.Board[3][m.From[1]]
		next.Board[3][m.From[1]] = empty
		next.HalfmoveClock = HalfmoveClock(0)
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.What().FigType == Rook {
		var empty Square
		czyempty := next.Board[m.To[0]][m.To[1]].Empty()
		next.Board[m.To[0]][m.To[1]] = next.Board[m.From[0]][m.From[1]]
		next.Board[m.From[0]][m.From[1]] = empty
		if m.From[0] == 0 {
			if m.From[1]%8 == 0 {
				next.Castling = next.Castling.OffRook(m.Before.MovesNext, 'Q')
			} else if m.From[1]%8 == 7 {
				next.Castling = next.Castling.OffRook(m.Before.MovesNext, 'K')
			}
		}
		if czyempty {
			next.HalfmoveClock++
		} else {
			next.HalfmoveClock = HalfmoveClock(0)
		}
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if next.Board[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			next.MoatsState[m.From[1]/8] = true
			next.MoatsState[m.From[1]/8+1] = true
		}
	} else if m.What().FigType == King {
		var empty Square
		czyempty := next.Board[m.To[0]][m.To[1]].Empty()
		next.Board[m.To[0]][m.To[1]] = next.Board[m.From[0]][m.From[1]]
		next.Board[m.From[0]][m.From[1]] = empty
		next.Castling = next.Castling.OffKing(m.Before.MovesNext)
		if czyempty {
			next.HalfmoveClock++
		} else {
			next.HalfmoveClock = HalfmoveClock(0)
		}
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if next.Board[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			next.MoatsState[m.From[1]/8] = true
			next.MoatsState[m.From[1]/8+1] = true
		}
	} else if m.What().FigType == Pawn {
		var empty Square
		//czyempty := nextboard[m.To[0]][m.To[1]].Empty()
		next.Board[m.To[0]][m.To[1]] = next.Board[m.From[0]][m.From[1]]
		next.Board[m.From[0]][m.From[1]] = empty
		next.HalfmoveClock = HalfmoveClock(0)
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if next.Board[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			next.MoatsState[m.From[1]/8] = true
			next.MoatsState[m.From[1]/8+1] = true
		}
	} else {
		var empty Square
		czyempty := next.Board[m.To[0]][m.To[1]].Empty()
		next.Board[m.To[0]][m.To[1]] = next.Board[m.From[0]][m.From[1]]
		next.Board[m.From[0]][m.From[1]] = empty
		if czyempty {
			next.HalfmoveClock++
		} else {
			next.HalfmoveClock = HalfmoveClock(0)
		}
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
		moatbridging := true
		for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ {
			if next.Board[0][i].NotEmpty {
				moatbridging = false
			}
		}
		if moatbridging {
			next.MoatsState[m.From[1]/8] = true
			next.MoatsState[m.From[1]/8+1] = true
		}
	}

	if next.AmIInCheck(m.What().Color) {
		return &next, IllegalMoveError{m, "Check", "We would be in check!"} //Bug(ArchieT): returns it even if we would not
	}

	return &next, nil
}
