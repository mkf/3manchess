package ai

import "github.com/ArchieT/3manchess/player"

type Player interface {
	player.Player
	Config() Config
}

type Config interface {
	String() string
	Byte() []byte
}
