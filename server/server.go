package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"
import "time"

type Server interface {
	Initialize(username string, password string, database string) error
	SaveGP(*GameplayData) (key int64, err error)
	LoadGP(key int64, gp *GameplayData) error
	SaveSD(*game.StateData) (key int64, err error)
	LoadSD(key int64, sd *game.StateData) error
	SavePD(*player.PlayerData) (key int64, err error)
	LoadPD(key int64, pd *player.PlayerData) error
	GetDerived(key int64) (keys []int64, err error)
}

func SaveState(sr Server, st *game.State) (key string, err error) {
	return sr.SaveSD(st.Data())
}

func LoadState(sr Server, key string, s *game.State) error {
	var sd game.StateData
	err := sr.LoadSD(k, &sd)
	s.FromData(&sd)
	return err
}

func SavePlayer(sr Server, pl player.Player) (key string, err error) {
	s := pl.Data()
	return sr.SavePD(&s)
}

func LoadPlayer(sr Server, key string, p player.Player) error {
	var s player.PlayerData
	err := sr.LoadPD(k, &s)
	p.FromData(s)
	return err
}

type GameplayData struct {
	State, White, Gray, Black int64
	Date                      time.Time
}

func FromGameplay(sr Server, gp player.Gameplay) (*GameplayData, error) {
	var d GameplayData
	var err error
	d.Date = time.Now()
	d.State, err = SaveState(sr, gp.State)
	if err != nil {
		return d, err
	}
	d.White, err = SavePlayer(sr, gp.Players[game.White], c)
	if err != nil {
		return d, err
	}
	d.Gray, err = SavePlayer(sr, gp.Players[game.Gray], c)
	if err != nil {
		return d, err
	}
	d.Black, err = SavePlayer(sr, gp.Players[game.Black], c)
	return d, err
}

func SaveGameplay(sr Server, gp player.Gameplay) (key string, err error) {
	d, err := FromGameplay(sr, gp)
	if err != nil {
		return nil, err
	}
	return sr.SaveGP(d)
}
