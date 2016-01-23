package gui

import "log"

//import "github.com/ArchieT/3manchess/movedet"
import "github.com/ArchieT/3manchess/game"
import "math/cmplx"
import "math"

const (
	Center             = 350 + 350i
	InnerRadius        = 70
	SubsequentRadiiAdd = 35
	DefaultRotation    = -math.Pi / 6
	OneFile            = math.Pi / 12
)

//adowbiowl — Angle Depening On Whether Black Is On Whites Left
func adowbiowl(p float64, biowl bool) float64 {
	if !biowl {
		return math.Remainder(p, 2*math.Pi)
	}
	return math.Remainder(p+math.Pi, 2*math.Pi)
}

type GUI struct {
	appears             chan<- game.BoardDiff
	BlackIsOnWhitesLeft bool
	fromtos             <-chan game.FromTo
	Rotated             float64 //zerofile blackmost boundary angle
	errchan             chan error
	ErrorChan           <-chan error
	GUIEngine
	bm *boardmap
}

type GUIEngine interface {
	Initialize(Boardclicker) error
	Appear(game.BoardDiff)
	ErrorChan() <-chan error
}

type boardmap [6][24]string

func (gui *GUI) Appear(w game.BoardDiff) {
	gui.bm[w.Pos[0]][w.Pos[1]] = FigURIs[w.Fig.Uint8()]
	gui.GUIEngine.Appear(w)
}

type Boardclicker struct {
	c     chan game.Pos
	rot   *float64
	biowl *bool
}

func (bckr Boardclicker) ClickedIt(x, y int) {
	bckr.c <- clicking(complex(float64(x), float64(y)), *bckr.rot, *bckr.biowl)
}

func (bckr Boardclicker) ClickPos(rank, file int8) error {
	p := game.Pos{rank, file}
	if err := p.Correct(); err == nil {
		bckr.c <- p
	} else {
		return err
	}
	return nil
}

func clicking(c complex128, rot float64, biowl bool) game.Pos {
	var r, p float64
	var m uint16
	var pr, pf int8
	log.Println("RawClick:", c)
	c -= Center
	r, p = cmplx.Polar(c)
	p -= rot
	r -= InnerRadius
	if r < 0 {
		return game.Pos{-1, -1}
	}
	p = adowbiowl(p, biowl)
	m = uint16(r) / 35
	if m < 24 {
		pr = int8(m)
	} else {
		return game.Pos{127, 127}
	}
	pf = int8(p / OneFile)
	return game.Pos{pr, pf}
}

func fromtoing(s <-chan game.Pos, d chan<- game.FromTo) {
	var f game.Pos
	for {
		f = <-s
		d <- game.FromTo{f, <-s}
	}
}

func NewGUI(ge GUIEngine) (*GUI, error) {
	gui := new(GUI)
	var clicks Boardclicker
	clicks.rot = &gui.Rotated
	clicks.biowl = &gui.BlackIsOnWhitesLeft
	clickpos := make(chan game.Pos)
	clicks.c = clickpos
	appears := make(chan game.BoardDiff)
	fromtos := make(chan game.FromTo)
	gui.appears = appears
	gui.Rotated = DefaultRotation
	gui.fromtos = fromtos
	go fromtoing(clickpos, fromtos)
	err := ge.Initialize(clicks)
	if err != nil {
		return gui, err
	}
	gui.errchan = make(chan error)
	gui.ErrorChan = gui.errchan
	return gui, nil
}
