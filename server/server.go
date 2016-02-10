package server

import "github.com/ArchieT/3manchess/player"
import "github.com/ArchieT/3manchess/game"

type Server interface {
	SaveGP(*player.Gameplay) (key string, err error)
	LoadGP(key string, gp *player.Gameplay) error
	SaveSD(*game.StateData) (key string, err error)
	LoadSD(key string, sd *game.StateData) error
	SavePD(*player.PlayerData) (key string, err error)
	LoadPD(key string, pd *player.PlayerData) error
	GetDerived(key string) (keys []string, err error)
}

func (sr Server) SaveState(st *game.State) (key string, err error) {
	return sr.SaveSD(st.Data())
}

func (sr Server) LoadState(key string, s *game.State) error {
	var sd game.StateData
	err := sr.LoadSD(k, &sd)
	s.FromData(&sd)
	return err
}

func (sr Server) SavePlayer(pl player.Player) (key string, err error) {
	s := pl.Data()
	return sr.SavePD(&s)
}

func (sr Server) LoadPlayer(key string, p player.Player) error {
	var s player.PlayerData
	err := sr.LoadPD(k, &s)
	p.FromData(s)
	return err
}

type GameplayData struct {
	State, White, Gray, Black string
	Date                      time.Time
}

func (gd GameplayData) FromGameplay(gp player.Gameplay) error {
	var d GameplayData
	var err error
	d.Date = time.Now()
	d.State, err = SaveState(gp.State)
	if err != nil {
		return d, err
	}
	d.White, err = SavePlayer(gp.Players[game.White], c)
	if err != nil {
		return d, err
	}
	d.Gray, err = SavePlayer(gp.Players[game.Gray], c)
	if err != nil {
		return d, err
	}
	d.Black, err = SavePlayer(gp.Players[game.Black], c)
	return d, err
}

func (sr Server) SaveGameplay(gp player.Gameplay) (key string, err error) {
	var d GameplayData
	if err := d.FromGameplay(gp); err != nil {
		return nil, err
	}
	return sr.SaveGP(&d)
}
