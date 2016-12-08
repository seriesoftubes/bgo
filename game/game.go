package game

import (
	"time"

	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
)

type Game struct {
	Board             *Board
	CurrentPlayer     *plyr.Player
	CurrentRoll       *Roll
	numHumanPlayers   uint8
	currentValidTurns map[string]turn.Turn
}

func (g *Game) ValidTurns() map[string]turn.Turn      { return g.currentValidTurns }
func (g *Game) SetValidTurns(vt map[string]turn.Turn) { g.currentValidTurns = vt }

func (g *Game) NextPlayersTurn() {
	if g.CurrentPlayer == plyr.PCC {
		g.CurrentPlayer = plyr.PC
	} else {
		g.CurrentPlayer = plyr.PCC
	}

	g.CurrentRoll = newRoll()
}

func (g *Game) HasAnyHumans() bool    { return g.numHumanPlayers > 0 }
func (g *Game) HasAnyComputers() bool { return g.numHumanPlayers < 2 }
func (g *Game) IsCurrentPlayerHuman() bool {
	if g.numHumanPlayers == 2 {
		return true
	} else if g.numHumanPlayers == 1 {
		return g.CurrentPlayer != plyr.PC // The `PC` player is always the computer
	} else {
		return false
	}
}

func NewGame(numHumanPlayers uint8) *Game {
	b := &Board{}
	b.SetUp()

	player := plyr.PCC
	if time.Now().UnixNano()%2 == 0 {
		player = plyr.PC
	}

	return &Game{Board: b, CurrentPlayer: player, CurrentRoll: newRoll(), numHumanPlayers: numHumanPlayers}
}
