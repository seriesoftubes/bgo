package main

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/render"
)

func main() {
	g := game.NewGame()
	render.PrintGame(g)

	cop := game.CopyGame(*g)
	pointerCop := game.GamePointer(cop) // etc
	fmt.Println(pointerCop == g, cop == *g, cop.Board == g.Board) // not a deep copy
	// TODO: try executing a legal move here and print game.
	// also try making a copy of the game by dereferencing the pointer?
}
