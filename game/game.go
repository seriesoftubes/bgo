package game

import (
	"time"
)

type Game struct {
	Board             *Board
	CurrentPlayer     *Player
	CurrentRoll       *Roll
	numHumanPlayers   uint8
	currentValidTurns map[string]Turn
}

func (g *Game) ValidTurns() map[string]Turn      { return g.currentValidTurns }
func (g *Game) SetValidTurns(vt map[string]Turn) { g.currentValidTurns = vt }

func (g *Game) NextPlayersTurn() {
	if g.CurrentPlayer == PCC {
		g.CurrentPlayer = PC
	} else {
		g.CurrentPlayer = PCC
	}

	g.CurrentRoll = newRoll()
}

func (g *Game) HasAnyHumans() bool    { return g.numHumanPlayers > 0 }
func (g *Game) HasAnyComputers() bool { return g.numHumanPlayers < 2 }
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

	return &Game{Board: b, CurrentPlayer: player, CurrentRoll: newRoll(), numHumanPlayers: numHumanPlayers}
}
