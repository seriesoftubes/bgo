package render

import (
  "fmt"

  "github.com/seriesoftubes/bgo/game"
)

func PrintGame(g *game.Game) {
  fmt.Println("\n\tCurrent player:\t", g.CurrentPlayer.Symbol())
  fmt.Println("\n\tDice:\t\t", *g.CurrentRoll)
  printBoard(g.Board)
}
