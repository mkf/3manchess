package game

import "testing"
import "fmt"

func showamft(p Pos) {
	var wyjscie [6][24]bool
	for _, val := range AMFT[p] {
		wyjscie[val[0]][val[1]] = true
	}
	fmt.Println(p)
	for i := int8(5); i >= 0; i-- {
		for j := int8(0); j < 24; j++ {
			if wyjscie[i][j] {
				fmt.Print("▓")
			} else if j == p[1] && i == p[0] {
				fmt.Print("█")
			} else {
				fmt.Print("░")
			}
		}
		fmt.Println()
	}
}

func TestAMFT_zero(t *testing.T) { showamft(Pos{0, 0}) }
func TestAMFT_5a12(t *testing.T) { showamft(Pos{5, 12}) }
