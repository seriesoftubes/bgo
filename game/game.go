package game

import (
	"time"
)

type Game struct {
	Board           *Board
	CurrentPlayer   *Player
	CurrentRoll     *Roll
	numHumanPlayers uint8
}

func (g *Game) NextPlayersTurn() {
	if g.CurrentPlayer == PCC {
		g.CurrentPlayer = PC
	} else {
		g.CurrentPlayer = PCC
	}

	g.CurrentRoll = newRoll()
}

func (g *Game) HasAnyHumans() bool { return g.numHumanPlayers > 0 }
func (g *Game) IsCurrentPlayerHuman() bool {
	if g.numHumanPlayers == 2 {
		return true
	} else if g.numHumanPlayers == 1 {
		return g.CurrentPlayer != PC // The `PC` player is always the computer
	} else {
		return false
	}
}

func NewGame(numHumanPlayers uint8) *Game {
	b := &Board{}
	b.setUp()

	player := PCC
	if time.Now().UnixNano()%2 == 0 {
		player = PC
	}

	return &Game{b, player, newRoll(), numHumanPlayers}
}
