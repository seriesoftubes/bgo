package main

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/render"
)

func readTurnFromStdin(validTurns map[string]game.Turn) game.Turn {
	for {
		var supposedlySerializedTurn string
		fmt.Scanln(&supposedlySerializedTurn)

		t, err := game.DeserializeTurn(supposedlySerializedTurn)
		if err != nil {
			fmt.Println("could not read your instructions, please try again: " + err.Error())
			continue
		}

		if _, ok := validTurns[t.String()]; ok {
			return t
		} else {
			fmt.Println("invalid turn entered, please try again")
		}
	}
}

func twoPlayerLoop(g *game.Game) bool {
	render.PrintGame(g)

	validTurns := game.ValidTurns(g.Board, g.CurrentRoll, g.CurrentPlayer)
	var chosenTurn game.Turn
	if len(validTurns) == 0 {
		fmt.Println("\tno moves available, sorry!")
	} else if len(validTurns) == 1 {
		for _, t := range validTurns {
			chosenTurn = t
		}
		fmt.Println("\tonly 1 move available, forcing", chosenTurn)
	} else {
		fmt.Println(fmt.Sprintf("\tyour move, %q:", *g.CurrentPlayer))
		chosenTurn = readTurnFromStdin(validTurns)
	}

	for move, numTimes := range chosenTurn {
		mp := &move
		for i := uint8(0); i < numTimes; i++ {
			if ok := g.Board.ExecuteMoveIfLegal(mp); !ok {
				panic("somehow, even with `validTurns` supplied, we couldn't execute a move: " + mp.String())
			}
		}
	}

	if g.Board.Winner() != nil {
		return true
	} else {
		g.NextPlayersTurn()
		return false
	}
}

func main() {
	g := game.NewGame()

	var done bool
	for !done {
		done = twoPlayerLoop(g)
	}

	render.PrintGame(g)
	fmt.Println("\tDONE WITH GAME!")
}
