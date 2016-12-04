package game

import (
	"time"
)

type Game struct {
	Board         *Board
	CurrentPlayer *Player
	CurrentRoll   *Roll
}

func (g *Game) NextPlayersTurn() {
	if g.CurrentPlayer == PCC {
		g.CurrentPlayer = PC
	} else {
		g.CurrentPlayer = PCC
	}

	g.CurrentRoll = newRoll()
}

func NewGame() *Game {
	b := &Board{}
	b.setUp()

	player := PCC
	if time.Now().UnixNano()%2 == 0 {
		player = PC
	}

	return &Game{b, player, newRoll()}
}
