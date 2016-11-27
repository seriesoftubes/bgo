package game

import (
	"time"
)

type Game struct {
	Board         *Board
	CurrentPlayer *Player
	CurrentRoll   *Roll
}

func CopyGame(g Game) Game     { return g } // is this a deep copy though?
func GamePointer(g Game) *Game { return &g }

func NewGame() *Game {
	b := &Board{}
	b.setUp()

	player := PCC
	if time.Now().UnixNano()%2 == 0 {
		player = PC
	}

	return &Game{b, player, newRoll()}
}
