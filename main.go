package main

import (
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/render"
)

func main() {
	g := game.NewGame()
	render.PrintGame(g)
	// TODO: try executing a legal move here and print game.
	// also try making a copy of the game by dereferencing the pointer?

	/*
		b := &game.Board{}
		b.Points = [game.NUM_BOARD_POINTS]*game.BoardPoint{
			// counter-clockwise player is in bottom-left.
			{}, {game.PC, 1}, {game.PC, 2}, {game.PC, 2}, {game.PC, 5}, {game.PC, 5}, {}, {}, {}, {}, {}, {},
			{}, {}, {}, {}, {}, {}, {game.PCC, 5}, {game.PCC, 1}, {game.PCC, 1}, {game.PCC, 2}, {game.PCC, 1}, {game.PCC, 5},
			//                                                        clockwise player in top-left.
		}
		render.PrintBoard(b)
	*/
}
