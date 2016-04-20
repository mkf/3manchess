package game

import "testing"

func TestAMFT_filetranslate(t *testing.T) {
	for rank := int8(0); rank < 6; rank++ {
		dlazera := AMFT[Pos{rank, 0}]
		for file := int8(0); file < 24; file++ {
			dlategotu := AMFT[Pos{rank, file}]
			if len(dlategotu) != len(dlazera) {
				t.Error(len(dlazera), len(dlategotu), rank, file, "\n", dlazera, "\n", dlategotu)
			} else {
				for el := range dlazera {
					orel := dlazera[el].AddVector([2]int8{0, file})
					if orel != dlategotu[el] {
						t.Error(orel, dlategotu[el], "\n", dlazera, "\n", dlategotu, "\n", dlazera[el], rank, file)
					}
				}
			}
		}
	}
}
