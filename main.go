package main

import (
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/render"
)

func main() {
	g := game.NewGame()
	render.PrintGame(g)
}
