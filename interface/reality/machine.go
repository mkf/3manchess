package reality

//import "github.com/ArchieT/3manchess/game"

type Pins struct{}

type MPos struct {
	Pos           [2]int8
	Heigth, Catch int8
}

type Machine struct {
	FilePins, RankPins, HeightPins, PullPins Pins
	relpos                                   MPos
	stop                                     chan bool
}

func NewMachine(file, rank, height, pull Pins) (*Machine, error) {
	mc := new(Machine)
	mc.stop = make(chan bool)
	return mc, nil
}

func (mc *Machine) StopChan() chan<- bool {
	return mc.stop
}

func (mc *Machine) RankMove(from, to int8) error {
	return nil
}

func (mc *Machine) FileMove(from, to int8) error {
	return nil
}

func (mc *Machine) SquareMove(from, to [2]int8) error {
	errchan := make(chan error)
	go func() {
		errchan <- mc.RankMove(from[0], to[0])
	}()
	go func() {
		errchan <- mc.FileMove(from[1], to[1])
	}()
	if err := <-errchan; err != nil {
		return err
	}
	return <-errchan
}

func (mc *Machine) HeightMove(from, to int8) error { //how high is the "palm"
	return nil
}

func (mc *Machine) PullMove(from, to int8) error { //how closed are the "fingers"
	return nil
}

func (mc *Machine) GivePos() MPos { //MPos,error {
	return mc.relpos
}
