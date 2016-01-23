package desktop

import "gopkg.in/qml.v1"

//import "log"
import "github.com/ArchieT/3manchess/interface/gui"

type DesktopEngine struct {
	engine    *qml.Engine
	component *qml.Object
	window    *qml.Window
	ErrChan   <-chan error
	errchan   chan error
}

func (de *DesktopEngine) Initialize(clicks gui.Boardclicker) error {
	de.engine = qml.NewEngine()
	de.engine.Context().SetVar("clickinto", clicks)
	component, err := de.engine.LoadFile("okno.qml")
	de.component = &component
	if err != nil {
		return err
	}
	de.errchan = make(chan error)
	de.ErrChan = de.errchan
	de.window = component.CreateWindow(nil)
	de.window.Show()
	go de.run()
	return nil
}

func (de *DesktopEngine) ErrorChan() <-chan error { return de.ErrChan }

func (de *DesktopEngine) run() {
	de.window.Wait()
}
