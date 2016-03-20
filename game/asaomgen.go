package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//ASAOMGen : generates all states after opponents' moves
func ASAOMGen(gamestate *State, ourcolor Color) <-chan *State {
	all_final_states := make(chan *State)
	go func() {
		if !gamestate.PlayersAlive.Give(ourcolor) || gamestate.MovesNext == ourcolor {
			all_final_states <- gamestate
		} else { // we and someone else is alive plus it's not our turn to move
			for move_op1 := range VFTPGen(gamestate) {
				move_to_apply := move_op1.Move(gamestate)
				state2, _ := move_to_apply.EvalAfter()
				if !state2.PlayersAlive.Give(ourcolor) || state2.MovesNext == ourcolor {
					all_final_states <- state2
				} else {
					for move_op2 := range VFTPGen(state2) {
						move_to_apply := move_op2.Move(state2)
						endstate, _ := move_to_apply.EvalAfter()
						all_final_states <- endstate
					}
				}
			}
		}
		close(all_final_states)
	}()
	return all_final_states
}
