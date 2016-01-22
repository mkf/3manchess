package camget

import "github.com/lazywei/go-opencv/opencv"
import "fmt"

//import "github.com/ArchieT/3manchess/game"
import "github.com/ArchieT/3manchess/movedet/board"

//import "log"

type NewCamError string

func (nce NewCamError) Error() string {
	return string(nce)
}

type GiveFrameError string

func (gfe GiveFrameError) Error() string {
	return string(gfe)
}

type Camera struct {
	*opencv.Capture
	Index int
}

type Calibration struct {
	Angle   float64
	Circles [7](*Circle) //[0]:Outer [6]:Inner
}

type View struct {
	*Camera
	*Calibration
}

type Circle struct {
	TopTangent, BottomTangent, LeftTangent, RightTangent int
	MoatBW, MoatWG, MoatGB                               [2]int
}

func NewCam(index int) (Camera, error) {
	src := opencv.NewCameraCapture(index)
	cam := Camera{src, index}
	if src == nil {
		return cam, NewCamError("It is nil!")
	}
	go func() {
		defer src.Release()
		for {
		}
	}()
	return cam, nil
}

func (c Camera) GiveFrame() (*opencv.IplImage, error) {
	if c.Capture.GrabFrame() {
		return c.Capture.RetrieveFrame(1), nil
	} else {
		return c.Capture.RetrieveFrame(1), GiveFrameError("Grab failed")
	}
	//return c.Capture.QueryFrame(), nil
}

func (v *View) GiveBoard() (*board.Board, error) {
	var b board.Board
	return &b, nil
}

func (v *View) Calibrate() {
	win := opencv.NewWindow("Calibration")
	defer win.Destroy()

	overfunc := func(s string, index int) (string, int, func(pos int, param ...interface{})) {
		var ours string
		ours = fmt.Sprintf("%s%d", s, index)
		var myf func(pos int, param ...interface{})
		var pix int
		switch s {
		case "Top":
			myf = func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].TopTangent = pos
			}
			pix = 1943
		case "Left":
			myf = func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].LeftTangent = pos
			}
			pix = 2591
		case "Right":
			myf = func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].RightTangent = pos
			}
			pix = 2591
		case "Bottom":
			myf = func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].BottomTangent = pos
			}
			pix = 1943
		default:
			panic("None of them???")
		}
		return ours, pix, myf
	}

	for i := 0; i < 7; i++ {
		tours, tpix, tmyf := overfunc("Top", i)
		lours, lpix, lmyf := overfunc("Left", i)
		rours, rpix, rmyf := overfunc("Right", i)
		bours, bpix, bmyf := overfunc("Bottom", i)
		win.CreateTrackbar(tours, 0, tpix, tmyf)
		win.CreateTrackbar(lours, 0, lpix, lmyf)
		win.CreateTrackbar(rours, 0, rpix, rmyf)
		win.CreateTrackbar(bours, 0, bpix, bmyf)
	}

	//win.CreateTrackbar("Left5", 0, 2591, func(pos int, param ...interface{}) {
	//	v.Calibration.Circles[5].LeftTangent = pos
	//})
}
