package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//ASAOMGen : generates all states after opponents' moves
func ASAOMGen(gamestate *State, ourcolor Color) <-chan *State {
	allFinalStates := make(chan *State)
	go func() {
		if !gamestate.PlayersAlive.Give(ourcolor) || gamestate.MovesNext == ourcolor {
			allFinalStates <- gamestate
		} else { // we and someone else is alive plus it's not our turn to move
			for moveOp1 := range VFTPGen(gamestate) {
				moveToApply := moveOp1.Move(gamestate)
				state2, _ := moveToApply.EvalAfter()
				if !state2.PlayersAlive.Give(ourcolor) || state2.MovesNext == ourcolor {
					allFinalStates <- state2
				} else {
					for moveOp2 := range VFTPGen(state2) {
						moveToApply := moveOp2.Move(state2)
						endstate, _ := moveToApply.EvalAfter()
						allFinalStates <- endstate
					}
				}
			}
		}
		close(allFinalStates)
	}()
	return allFinalStates
}
