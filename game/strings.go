package game

//Rune : Unicode representation of a piece
func (f *Fig) Rune() rune {
	switch f.Color {
	case White:
		switch f.FigType {
		case Pawn:
			return 'P'
		case Rook:
			return 'R'
		case Knight:
			return 'N'
		case Bishop:
			return 'B'
		case Queen:
			return 'Q'
		case King:
			return 'K'
		}
	case Gray:
		switch f.FigType {
		case Pawn:
			return '♙'
		case Rook:
			return '♖'
		case Knight:
			return '♘'
		case Bishop:
			return '♗'
		case Queen:
			return '♕'
		case King:
			return '♔'
		}
	case Black:
		switch f.FigType {
		case Pawn:
			return 'p'
		case Rook:
			return 'r'
		case Knight:
			return 'n'
		case Bishop:
			return 'b'
		case Queen:
			return 'q'
		case King:
			return 'k'
		}
	}
	return '?'
}

//FromRune parses a rune into FigType and Color
func FromRune(r rune) (FigType, Color) {
	switch r {
	case 'P':
		return Pawn, White
	case 'R':
		return Rook, White
	case 'N':
		return Knight, White
	case 'B':
		return Bishop, White
	case 'Q':
		return Queen, White
	case 'K':
		return King, White
	case 'p':
		return Pawn, Black
	case 'r':
		return Rook, Black
	case 'n':
		return Knight, Black
	case 'b':
		return Bishop, Black
	case 'q':
		return Queen, Black
	case 'k':
		return King, Black
	case '♙':
		return Pawn, Gray
	case '♖':
		return Rook, Gray
	case '♘':
		return Knight, Gray
	case '♗':
		return Bishop, Gray
	case '♕':
		return Queen, Gray
	case '♔':
		return King, Gray
	default:
		var a FigType
		var b Color
		return a, b
	}
}
