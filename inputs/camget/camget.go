package camget

import "github.com/lazywei/go-opencv/opencv"

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
		return cam, "It is nil!"
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
		return c.Capture.RetrieveFrame(1)
	}
	//return c.Capture.QueryFrame(), nil
}

func (v *View) Calibrate() {
	win := opencv.NewWindow("Calibration")
	defer win.Destroy()

	overfunc := func(s string, index int) (string, int, func(pos int, param ...interface{})) {
		var ours string
		ours = s + index.String()
		switch s {
		case "Top":
			myf := func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].TopTangent = pos
			}
			pix := 1943
		case "Left":
			myf := func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].LeftTangent = pos
			}
			pix := 2591
		case "Right":
			myf := func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].RightTangent = pos
			}
			pix := 2591
		case "Bottom":
			myf := func(pos int, param ...interface{}) {
				v.Calibration.Circles[index].BottomTangent = pos
			}
			pix := 1943
		default:
			panic("None of them???")
		}
		return ours, pix, myf
	}

	win.CreateTrackbar("Top0", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[0].TopTangent = pos
	})
	win.CreateTrackbar("Top1", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[1].TopTangent = pos
	})
	win.CreateTrackbar("Top2", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[2].TopTangent = pos
	})
	win.CreateTrackbar("Top3", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[3].TopTangent = pos
	})
	win.CreateTrackbar("Top4", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[4].TopTangent = pos
	})
	win.CreateTrackbar("Top5", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[5].TopTangent = pos
	})
	win.CreateTrackbar("Top6", 0, 1943, func(pos int, param ...interface{}) {
		v.Calibration.Circles[6].TopTangent = pos
	})
	win.CreateTrackbar("Left0", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[0].LeftTangent = pos
	})
	win.CreateTrackbar("Left1", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[1].LeftTangent = pos
	})
	win.CreateTrackbar("Left2", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[2].LeftTangent = pos
	})
	win.CreateTrackbar("Left3", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[3].LeftTangent = pos
	})
	win.CreateTrackbar("Left4", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[4].LeftTangent = pos
	})
	win.CreateTrackbar("Left5", 0, 2591, func(pos int, param ...interface{}) {
		v.Calibration.Circles[5].LeftTangent = pos
	})
	//win.CreateTrackbar("Left6",0,2591,func(pos
}
