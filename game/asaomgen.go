package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//ASAOMGen : generates all states after opponents' moves
func ASAOMGen(gamestate *State, ourcolor Color) <-chan *State {
	all_final_states := make(chan *State)
	go func() {
		if gamestate.MovesNext == ourcolor {
			all_final_states <- gamestate
		} else {
			for move_op1 := range VFTPGen(gamestate) {
				move_to_apply := Move{move_op1.FromTo[0], move_op1.FromTo[1], gamestate, move_op1.PawnPromotion}
				state2, _ := move_to_apply.After()
				if state2.MovesNext == ourcolor {
					all_final_states <- state2
				} else {
					for move_op2 := range VFTPGen(state2) {
						move_to_apply := Move{move_op2.FromTo[0], move_op2.FromTo[1], state2, move_op2.PawnPromotion}
						endstate, _ := move_to_apply.After()
						all_final_states <- endstate
					}
				}
			}
		}
		close(all_final_states)
	}()
	return all_final_states
}
