package game

import "testing"
import "fmt"

func TestAMFT_filetranslate(t *testing.T) {
	for rank := int8(0); rank < 6; rank++ {
		dlazera := AMFT[Pos{rank, 0}]
		for file := int8(0); file < 24; file++ {
			dlategotu := AMFT[Pos{rank, file}]
			if len(dlategotu) != len(dlazera) {
				t.Error(len(dlazera), len(dlategotu), rank, file, "\n", dlazera, "\n", dlategotu)
			} else {
				for _, el := range dlazera {
					orel := el.AddVector([2]int8{0, file})
					jesttensam := false
					//if orel != dlategotu[el] {
					for _, dlael := range dlategotu {
						if orel == dlael {
							jesttensam = true
						}
					}
					if !jesttensam {
						t.Error(orel, "\n", dlazera, "\n", dlategotu, "\n", el, rank, file)
					}
				}
			}
		}
	}
}

func showamft(p Pos, vfile int8) {
	var wyjscie [6][24]bool
	for _, val := range AMFT[p] {
		wyjscie[val[0]][val[1]] = true
	}
	fmt.Println(p)
	for i := int8(5); i >= 0; i-- {
		for j := vfile; j < 24+vfile; j++ {
			if wyjscie[i][(j+24)%24] {
				fmt.Print("▓")
			} else if (j+24)%24 == p[1] && i == p[0] {
				fmt.Print("█")
			} else {
				fmt.Print("░")
			}
		}
		fmt.Println()
	}
}

func TestAMFT_0a0(t *testing.T)  { showamft(Pos{0, 0}, 12) }
func TestAMFT_0a11(t *testing.T) { showamft(Pos{0, 11}, -1) }
func TestAMFT_0a12(t *testing.T) { showamft(Pos{0, 12}, 0) }
func TestAMFT_0a7(t *testing.T)  { showamft(Pos{0, 7}, 7+12) }
func TestAMFT_5a0(t *testing.T)  { showamft(Pos{5, 0}, 12) }
func TestAMFT_5a12(t *testing.T) { showamft(Pos{5, 12}, 0) }
func TestAMFT_5a7(t *testing.T)  { showamft(Pos{5, 7}, 7+12) }
func TestAMFT_3a0(t *testing.T)  { showamft(Pos{3, 0}, 12) }
func TestAMFT_3a12(t *testing.T) { showamft(Pos{3, 12}, 0) }
func TestAMFT_3a7(t *testing.T)  { showamft(Pos{3, 7}, 7+12) }
