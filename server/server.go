package server

import "github.com/ArchieT/3manchess/player"

type Server interface {
	SaveGPD(*player.Gameplay) (key string)
	OpenGP(key string) *player.Gameplay
	GetDerived(key string) (derived_keys []string)
}
