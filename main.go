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
		fmt.Println("\tcan't do anything this turn, sorry!")
	} else if len(validTurns) == 1 {
		for _, t := range validTurns {
			chosenTurn = t
		}
		fmt.Println("\tthis turn only has 1 option, forcing!")
	} else {
		fmt.Println(fmt.Sprintf("\tYour move, %q:", *g.CurrentPlayer))
		chosenTurn = readTurnFromStdin(validTurns)
	}
	fmt.Println("\tChose move:", chosenTurn)

	g.Board.MustExecuteTurn(chosenTurn)

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
