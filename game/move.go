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

	next := *m.Before
	next.MovesNext = next.MovesNext.Next()
	if !m.Before.CanIMoveWOCheck(m.Before.MovesNext) {
		next.PlayersAlive.Die(m.Before.MovesNext)
		return &next, nil
	}

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
				next.Castling = next.Castling.OffRook(m.What().Color, 'Q')
			} else if m.From[1]%8 == 7 {
				next.Castling = next.Castling.OffRook(m.What().Color, 'K')
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
		next.Castling = next.Castling.OffKing(m.What().Color)
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
		return &next, IllegalMoveError{m, "Check", "We would be in check!"}
	}
	return &next, nil
}
