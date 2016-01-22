package gui

import "gopkg.in/qml.v1"
import "log"
import "github.com/ArchieT/3manchess/movedet/board"
import "github.com/ArchieT/3manchess/movedet"
import "github.com/ArchieT/3manchess/game"
import "math/cmplx"
import "math"

const (
	Center             = 350 + 350i
	InnerRadius        = 70
	SubsequentRadiiAdd = 35
	DefaultRotation    = math.Pi * (11 / 6)
	OneFile            = math.Pi / 12
)

//adowbiowl — Angle Depening On Whether Black Is On Whites Left
func adowbiowl(p float64, biowl bool) {
	if !biowl {
		return p % (2 * math.Pi)
	}
	return (p + math.Pi) % (2 * math.Pi)
}

type appearing struct {
	game.Pos
	board.Piece
}

type GUI struct {
	disappears          chan<- game.Pos
	appears             chan<- appearing
	replacements        chan<- appearing
	BlackIsOnWhitesLeft bool
	fromtos             <-chan game.FromTo
	Rotated             float64 //zerofile blackmost boundary angle
}

func replacing(r <-chan appearing, a chan<- appearing, d chan<- board.Pos) {
	var y appearing
	for {
		y = <-r
		d <- y.Pos
		a <- y
	}
}

func clicking(s <-chan complex64, d chan<- board.Pos, rot *float64, biowl *bool) {
	var c complex64
	var r, p float64
	var m uint16
	var pr, pf int8
	for {
		c = <-s
		c -= Center
		r, p = cmplx.Polar(c)
		p -= *rot
		r -= InnerRadius
		if r < 0 {
			continue
		}
		p = adowbiowl(p, biowl)
		m = uint16(r) / 35
		if m < 24 {
			pr = int8(m)
		} else {
			continue
		}
		pf = p / OneFile
		d <- game.Pos{pr, pf}
	}
}

func fromtoing(s <-chan game.Pos, d chan<- game.FromTo) {
	var f game.Pos
	for {
		f = <-s
		d <- game.FromTo{f, <-s}
	}
}

func NewGUI() (*GUI, error) {
	gui = new(GUI)
	clicks = make(chan complex64)
	clickpos = make(chan game.Pos)
	disappears = make(chan game.Pos)
	appears = make(chan appearing)
	replacements = make(chan appearing)
	fromtos = make(chan game.FromTo)
	gui.disappears = disappears
	gui.appears = appears
	gui.replacements = replacements
	gui.Rotated = DefaultRotation
	gui.fromtos = fromtos
	go replacing(replacements, appears, disappears)
	go clicking(clicks, clickpos, &(gui.Rotated), &(gui.BlackIsOnWhitesLeft))
	go fromtoing(clickpos, fromtos)
	return gui, nil
}
