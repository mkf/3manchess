package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//Move :  struct describing a single move with the situation before it
type Move struct {
	From, To Pos
	//	What        Fig
	//	AlreadyHere Fig
	Before        *State
	PawnPromotion FigType
}

//FromToProm is a struct containing FromTo and Prom
type FromToProm struct {
	FromTo        `json:"fromto"`
	PawnPromotion FigType `json:"pawnpromotion"`
}

//Move makes a Move from FromToProm and State pointer
func (ftp FromToProm) Move(bef *State) Move {
	m := ftp.FromTo.Move(bef)
	m.PawnPromotion = ftp.PawnPromotion
	return m
}

//IncorrectPos error
type IncorrectPos struct {
	Pos
}

func (ip IncorrectPos) Error() string {
	return ip.Pos.String()
}

//Correct checks if the Pos is 0≤r≤5 and 0≤f≤23, and returns IncorrectPos{Pos} if it is not
func (p Pos) Correct() error {
	if (p[0] < 0) || (p[0] > 5) || (p[1] < 0) || (p[1] > 23) {
		return IncorrectPos{p}
	}
	return nil
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
	return Move{ft.From(), ft.To(), before, 0}
}

//Correct checks if the FromTo is Pos.Correct
func (ft FromTo) Correct() error {
	if err := ft[0].Correct(); err != nil {
		return err
	}
	return ft[1].Correct()
}

//func (m *Move) String() string {
//}

//Where gives the Square of Before.Board[From]
func (m *Move) Where() Square {
	a := m.Before.Board.GPos(m.From)
	return *a
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

//WhereIsKing : where is king of specified color on the board?
func (b *Board) WhereIsKing(who Color) *Pos {
	for _, opos := range ALLPOS {
		if sq := b.GPos(opos); sq.NotEmpty && sq.Fig.Color == who && sq.Fig.FigType == King {
			return &opos
		}
	}
	return nil
}

//IsKingPresent : is king of specified color on the board?
func (b *Board) IsKingPresent(who Color) bool {
	return b.WhereIsKing(who) != nil
}

//CheckChecking :  is `who` in check?
func (b *Board) CheckChecking(who Color, pa PlayersAlive) Check { //true if in check
	if !pa.Give(who) {
		panic("CheckChecking a dead player!: " + who.String())
	}
	wking := b.WhereIsKing(who)
	if wking == nil {
		panic("King not found!!!: " + who.String())
	}
	return b.ThreatChecking(*wking, pa, DEFENPASSANT)
}

//ThreatChecking checks if the piece on where Pos is 'in check'
func (b *Board) ThreatChecking(where Pos, pa PlayersAlive, ep EnPassant) Check {
	who := b.GPos(where).Color()
	var heyitscheck Check
	for _, opos := range ALLPOS {
		if tjf := b.GPos(opos); tjf.NotEmpty && tjf.Color() != who && pa.Give(tjf.Color()) &&
			b.AnyPiece(opos, where, DEFMOATSSTATE, FALSECASTLING, ep, pa) {
			return Check{If: true, From: opos}
		}
	}
	return heyitscheck
}

//FriendsNAllies returns positions of our pieces and their pieces
func (b *Board) FriendsNAllies(who Color, pa PlayersAlive) ([]Pos, <-chan Pos) {
	my := make([]Pos, 0, 16)
	oni := make(chan Pos, 32)
	if pa.Give(who) {
		for _, opos := range ALLPOS {
			tjf := b.GPos(opos)
			if tjf.Color() == who {
				my = append(my, opos)
			} else if tjf.NotEmpty && pa.Give(tjf.Color()) {
				oni <- opos
			}
		}
	}
	close(oni)
	return my, oni
}

//WeAreThreateningTypes returns a list (not a set, dupicates included) of FigTypes we are 'checking'
func (b *Board) WeAreThreateningTypes(who Color, pa PlayersAlive, ep EnPassant) <-chan FigType {
	ret := make(chan FigType, 32)
	my, oni := b.FriendsNAllies(who, pa)
	for ich := range oni {
		for _, nasz := range my {
			if b.AnyPiece(nasz, ich, DEFMOATSSTATE, FALSECASTLING, ep, pa) {
				ret <- (*b)[ich[0]][ich[1]].Fig.FigType
				break
			}
		}
	}
	close(ret)
	return ret
}

//WeAreThreatened returns a list (not a set, dups included) of our FigTypes they are 'checking'
func (b *Board) WeAreThreatened(who Color, pa PlayersAlive, ep EnPassant) <-chan FigType {
	ret := make(chan FigType, 16)
	my, onichan := b.FriendsNAllies(who, pa)
	oni := make([]Pos, 0, len(onichan))
	for ich := range onichan {
		oni = append(oni, ich)
	}
	for _, nasz := range my {
		for _, ich := range oni {
			if b.AnyPiece(ich, nasz, DEFMOATSSTATE, FALSECASTLING, ep, pa) {
				ret <- (*b)[nasz[0]][nasz[1]].Fig.FigType
				break
			}
		}
	}
	close(ret)
	return ret
}

//TODO: Checkmate, stalemate detection. Doing something with the halfmove timer.

//Possible is such a move? Returns an error, same error as After() would give you, ¡¡¡except for CheckChecking!!!
func (m *Move) Possible() error {
	if err := m.From.Correct(); err != nil {
		return err
	}
	if err := m.To.Correct(); err != nil {
		return err
	}
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
	if merr := m.Possible(); merr != nil {
		return nil, merr
	}
	next := *m.Before
	nextboard := *m.Before.Board
	next.Board = &nextboard
	next.MovesNext = next.MovesNext.Next()

	if m.IsItKingSideCastling() {
		empty := next.Board[0][m.From[1]+2]                     //rather senseless, a lazy definition of an empty square
		next.Board[0][m.From[1]+2] = next.Board[0][m.From[1]]   //moving the king to his side
		next.Board[0][m.From[1]+1] = next.Board[0][m.From[1]+3] //moving the rook
		next.Board[0][m.From[1]] = empty                        //emptying king's square
		next.Board[0][m.From[1]+3] = empty                      //emptying rook's square
		next.Castling = next.Castling.OffKing(m.What().Color)
		next.HalfmoveClock++
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.IsItQueenSideCastling() {
		empty := next.Board[0][m.From[1]-2]                     //rather senseless, a lazy definition of an empty square
		next.Board[0][m.From[1]-2] = next.Board[0][m.From[1]]   //moving the king
		next.Board[0][m.From[1]-1] = next.Board[0][m.From[1]+4] //moving the rook
		next.Board[0][m.From[1]] = empty                        //emptying the king's square
		next.Board[0][m.From[1]+4] = empty                      //emptying the rook's square
		next.Castling = next.Castling.OffKing(m.What().Color)
		next.HalfmoveClock++
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.IsItPawnRunningEnPassant() {
		next.Board[3][m.From[1]] = next.Board[1][m.From[1]] //moving the pawn
		next.Board[1][m.From[1]] = next.Board[2][m.From[1]] //emptying the pawn's square
		next.HalfmoveClock = HalfmoveClock(0)               //zeroing half-move clock
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Appeared(Pos{2, m.From[1]}) //a new possibility to capture enpassant
	} else if m.IsItPawnCapturingEnPassant() {
		empty := next.Board[2][m.To[1]]
		next.Board[3][m.To[1]] = empty                    //removing the captured pawn
		next.Board[2][m.To[1]] = next.Board[3][m.From[1]] //moving the capturing pawn
		next.Board[3][m.From[1]] = empty                  //emptying the square of capturing pawn
		next.HalfmoveClock = HalfmoveClock(0)             //zeroing the half-move clock
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
	} else if m.What().FigType == Rook {
		var empty Square                                                //this time, we had to declare empty Square literally ;)
		czyempty := next.Board[m.To[0]][m.To[1]].Empty()                //check if the target square is empty
		next.Board[m.To[0]][m.To[1]] = next.Board[m.From[0]][m.From[1]] //move the piece
		next.Board[m.From[0]][m.From[1]] = empty                        //empty the piece's square
		if m.From[0] == 0 {                                             //if you start from the first rank
			if m.From[1]%8 == 0 { //if you are queenside by the moat
				next.Castling = next.Castling.OffRook(m.Before.MovesNext, 'Q')
			} else if m.From[1]%8 == 7 { //if you are kingside by the moat
				next.Castling = next.Castling.OffRook(m.Before.MovesNext, 'K')
			}
		}
		if czyempty { //if the target square is empty
			next.HalfmoveClock++
		} else {
			next.HalfmoveClock = HalfmoveClock(0) //capturing sth
		}
		next.FullmoveNumber++
		next.EnPassant = next.EnPassant.Nothing()
		moatbridging := true
		if !next.MoatsState[m.From[1]/8] || !next.MoatsState[m.From[1]/8+1] {
			for i := (m.From[1] / 8) * 8; i < ((m.From[1]/8)*8)+8; i++ { //check if all of the color's rank0 is empty
				if next.Board[0][i].NotEmpty { //if one of the squares is not empty
					moatbridging = false //then it is false
					break
				}
			}
		}
		if moatbridging { //if all of the color's rank0 is empty
			next.MoatsState[m.From[1]/8] = true   //bridge queenside
			next.MoatsState[m.From[1]/8+1] = true //bridge kingside
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
		if m.To[0] == 0 && m.From[0] == 1 {
			next.Board[m.To[0]][m.To[1]] = Square{NotEmpty: true, Fig: Fig{FigType: m.PawnPromotion, Color: m.What().Color, PawnCenter: false}}
			switch m.PawnPromotion {
			case ZeroFigType:
				return &next, IllegalMoveError{m, "ZeroPromotion", "Promotion to Zero"}
			case King:
				return &next, IllegalMoveError{m, "KingPromotion", "Promotion to King"}
			} //let's say that you can promote a pawn to a pawn
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

	if heyitscheck := next.AmIInCheck(m.What().Color); heyitscheck.If {
		return &next, IllegalMoveError{m, "Check", "We would be in check! (checking " + heyitscheck.From.String()} //Bug(ArchieT): returns it even if we would not
	}

	return &next, nil
}

//EvalAfter : return the evaluated gamestate afterwards, also error
func (m *Move) EvalAfter() (state *State, err error) {
	if state, err = m.After(); err == nil {
		state.FixMovesNext()
		state.EvalDeath()
		state.FixMovesNext()
	}
	return
}

// FixMovesNext : when someone is alive, sets MovesNext to the color, which will move next
func (s *State) FixMovesNext() {
	if !s.PlayersAlive.Give(s.MovesNext) {
		n := s.MovesNext
		s.MovesNext = s.MovesNext.Next()
		for !s.PlayersAlive.Give(s.MovesNext) && n != s.MovesNext { // while player is dead and we hadn't returned to player 'n'
			s.MovesNext = s.MovesNext.Next()
		}
	}
}
