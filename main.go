package main

import (
	"fmt"

	"github.com/seriesoftubes/bgo/ctrl"
	"github.com/seriesoftubes/bgo/game"
)

func main() {
	g := game.NewGame(1)
	ctrl := ctrl.New(g)
	winner, winKind := ctrl.PlayOneGame()
	fmt.Println(winner, winKind)
}
