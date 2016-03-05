package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "time"

type Server interface {
	Initialize(username string, password string, database string) error
	SaveGP(*GameplayData) (key int64, err error)
	LoadGP(key int64, gp *GameplayData) error
	SaveSD(sd *game.StateData, movekeyaddafter int64) (key int64, err error)
	LoadSD(key int64, sd *game.StateData) error
	//	SavePD(*player.PlayerData) (key int64, err error)
	//	LoadPD(key int64, pd *player.PlayerData) error
	SaveMD(*MoveData) (key int64, err error)
	LoadMD(key int64, md *MoveData) error
	ListGP() ([]GameFollow, error)
	////GetMovesFromState selects all the moves where the gamestate is before
	//GetMovesFromState(key int64) (keys []int64, err error)
}

type GameFollow struct {
	Id            int64 `json:"id"`
	*GameplayData `json:"game"`
}

type MoveData struct {
	FromTo        [4]int8
	BeforeGame    int64
	PawnPromotion int8
	Who           int64
}

func (md MoveData) Move(sr Server) game.Move {
	s := new(game.State)
	LoadState(sr, md.Before, s)
	return game.Move{
		From:          Pos{md.FromRank, md.FromFile},
		To:            Pos{md.ToRank, md.ToFile},
		Before:        s,
		PawnPromotion: game.FigType(md.PawnPromotion),
	}
}

type StateFollow struct {
	Key int64
	*game.StateData
}

func SaveState(sr Server, st *game.State) (key int64, err error) {
	return sr.SaveSD(st.Data())
}

func LoadState(sr Server, key int64, s *game.State) error {
	var sd game.StateData
	err := sr.LoadSD(k, &sd)
	s.FromData(&sd)
	return err
}

/*
type PlayerFollow struct {
	Key int64
	*player.PlayerData
}

func SavePlayer(sr Server, pl player.Player) (key int64, err error) {
	s := pl.Data()
	return sr.SavePD(&s)
}

func LoadPlayer(sr Server, key int64, p player.Player) error {
	var s player.PlayerData
	err := sr.LoadPD(k, &s)
	p.FromData(s)
	return err
}
*/

type GameplayData struct {
	State, White, Gray, Black int64
	Date                      time.Time
}

type GameplayFollow struct {
	Key int64
	*GameplayData
}

func FromGameplay(sr Server, gp player.Gameplay) (*GameplayData, error) {
	var d GameplayData
	var err error
	d.Date = time.Now()
	d.State, err = SaveState(sr, gp.State)
	if err != nil {
		return d, err
	}
	//d.White, err = SavePlayer(sr, gp.Players[game.White], c)
	if err != nil {
		return d, err
	}
	//d.Gray, err = SavePlayer(sr, gp.Players[game.Gray], c)
	if err != nil {
		return d, err
	}
	//d.Black, err = SavePlayer(sr, gp.Players[game.Black], c)
	return d, err
}

func SaveGameplay(sr Server, gp player.Gameplay) (key int64, err error) {
	d, err := FromGameplay(sr, gp)
	if err != nil {
		return nil, err
	}
	return sr.SaveGP(d)
}
