package render

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
)

func PrintGame(g *game.Game) {
	fmt.Println(fmt.Sprintf("\n\tPlayer: %s  Rolled: %v", g.CurrentPlayer.Symbol(), *g.CurrentRoll))
	printBoard(g.Board)
}
