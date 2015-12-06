package ai

import "../game"

//type capable of storing 3 players types
type PlayersTypes [3]byte

//representation of players' types
const (
	HUMAN = 'h'
	COMPUTER = 'c'
	DISLIKE_WHITE = 'w'
	DISLIKE_GREY = 'g'
	DISLIKE_BLACK = 'b'
)

//InitPlayers : init players
func InitPlayers(s string) *PlayersTypes {
	PLAYERS_TYPES := new(PlayersTypes)
	for i := 0 ; i < len(s) ; i++ {
		switch s[i] {
			case HUMAN, COMPUTER, DISLIKE_WHITE, DISLIKE_GREY, DISLIKE_BLACK:
				PLAYERS_TYPES[i] = s[i]
			default:
				panic("Invalid player given.")
		}
	}
	return PLAYERS_TYPES
}

//PlayerType : return PlayerType
func (PLAYERS_TYPES *PlayersTypes) PlayerType(c game.Color) byte {
	return PLAYERS_TYPES[c.UInt8()]
}
