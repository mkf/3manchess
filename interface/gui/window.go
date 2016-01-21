package gui

import "github.com/andlabs/ui"

type Window struct {
}

type BoardHandler struct {
	*Window
}

func (bh *BoardHandler) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
}

func (bh *BoardHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
}

func (bh *BoardHandler) MouseCrossed(a *ui.Area, left bool) {
}

func (bh *BoardHandler) DragBroken(a *ui.Area) {
}

func (bh *BoardHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	return false
}

func NewWindow() *BoardHandler {
	ourwindow := new(Window)
	bh := BoardHandler{ourwindow}
	err := ui.Main(func() {
		boardarea := ui.NewScrollingArea(&bh, 150, 150)
		window := ui.NewWindow("Hello", 300, 200, false)
		window.SetChild(boardarea)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
	return &bh
}
