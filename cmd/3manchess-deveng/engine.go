package main

import (
	"fmt"
	"github.com/ArchieT/3manchess/ai"
	"github.com/ArchieT/3manchess/game"
	"os"
)

func main() {
	fmt.Println(`3manchess experimental engine
Please add required players. Options are as follows:
h - human    c - computer
w,g,b - computer especially disliking White, Grey and Black player respectively
For now only 3 humans (hhh) are supported`)
	var input_players string
	fmt.Scanf("%s", &input_players)
	if len(input_players) != 3 {
		fmt.Println("Invalid number of players (%d given).", len(input_players))
		os.Exit(1)
	}
	players := ai.InitPlayers(input_players)
	game_state := game.InitializeNewGame()
	var input_move string
	var exit_code int
	for !game.IsItEndOfGame(game_state) { // TODO: implement
		if players.PlayerType(game_state.MovesNext) == ai.HUMAN { // human's move
			for {
				fmt.Println("%s (%d): ", game_state.MovesNext.String(), game_state.FullmoveNumber)
				fmt.Scanf("%s", &input_move)
				exit_code = game.ApplyMove(game_state, game.String2Move(input_move)) // TODO: implement ApplyMove
				if exit_code == game.OK || exit_code == game.ENDGAME {
					break
				}
				fmt.Println("Invalid move: %s", input_move)
			}
			if exit_code == game.ENDGAME {
				break
			}
			game_state.Board.Print()
		} else { // AI's move
			move := ai.Think(game_state, players) // TODO: implement
			fmt.Print("%s's ", game_state.MovesNext.String())
			exit_code = game.ApplyMove(game_state, move)
			fmt.Println("move is: %s", game.Move2String(move))
			game_state.Board.Print()
		}
		if exit_code == game.ENDGAME_HUMAN {
			fmt.Print("No more human players. Continue anyway? [Y/n] ")
			var answer byte
			fmt.Scanf("%c", &answer)
			if answer == 'N' || answer == 'n' {
				break
			}
		} else if exit_code == game.ENDGAME {
			fmt.Println("Game over. %s.", game.Winner(game_state))
			break
		}
	}
}
