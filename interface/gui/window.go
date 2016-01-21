package gui

import "github.com/andlabs/ui"

type Window struct {
}

func NewWindow() *Window {
	err := ui.Main(func() {
		name := ui.NewEntry()
		boardspace := ui.NewVerticalBox()
		boardspace.Append(ui.NewLabel("test"), false)
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetChild(boardspace)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
