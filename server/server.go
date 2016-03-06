package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "time"

type Server interface {
	Initialize(username string, password string, database string) error
	SaveGP(*GameplayData) (key int64, err error)
	LoadGP(key int64, gp *GameplayData) error
	SaveSD(sd *game.StateData) (key int64, err error)
	LoadSD(key int64, sd *game.StateData) error
	SaveMD(*MoveData) (key int64, err error)
	LoadMD(key int64, md *MoveData) error
	ListGP() ([]GameFollow, error)
	////GetMovesFromState selects all the moves where the gamestate is before
	//GetMovesFromState(key int64) (keys []int64, err error)
}

type GameplayData struct {
	State              int64
	White, Gray, Black *int64
	Date               time.Time
}

type MoveData struct {
	FromTo        [4]int8
	BeforeGame    int64
	AfterGame     int64
	PawnPromotion *int8
	Who           int64
}

func (md MoveData) Move(sr Server) game.Move {
	s := new(game.State)
	LoadState(sr, md.BeforeGame, s)
	return game.Move{
		From:          game.Pos{md.FromTo[0], md.FromTo[1]},
		To:            game.Pos{md.FromTo[2], md.FromTo[3]},
		Before:        s,
		PawnPromotion: game.FigType(md.PawnPromotion),
	}
}

func SaveState(sr Server, st *game.State, movekeyaddafter int64) (key int64, err error) {
	return sr.SaveSD(st.Data(), movekeyaddafter)
}

func LoadState(sr Server, key int64, s *game.State) error {
	var sd game.StateData
	err := sr.LoadSD(key, &sd)
	s.FromData(&sd)
	return err
}

/*
type StateFollow struct {
	Key int64
	*game.StateData
}
*/

/*
type GameplayFollow struct {
	Key int64 `json:"id"`
	*GameplayData `json:"game"`
}*/
