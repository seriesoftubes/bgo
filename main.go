package main

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/render"
)

func main() {
	g := game.NewGame()
	render.PrintGame(g)
	// TODO: try executing a legal move here and print game.
	// also try making a copy of the game by dereferencing the pointer?
}
