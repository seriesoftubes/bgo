package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

type (
	boardPointState struct {
		isOwnedByMe bool
		numChex     uint8
	}

	State struct {
		// points on the board, indexed with 0 being the furthest away from the current player's home
		boardPoints               [constants.NUM_BOARD_POINTS]boardPointState
		numOnMyBar, numOnEnemyBar uint8
		myRoll                    game.Roll
	}
)

func uint8Ceiling(x, max uint8) uint8 {
	if x > max {
		return max
	}
	return x
}

// Detects the current state of the game, truncating the checker counts up to a max.
// Returns the State and whether the State's boardPoints were reversed to account for the player's perspective.
func DetectState(p *plyr.Player, g *game.Game, maxChexToConsider uint8) (State, bool) {
	out := State{myRoll: g.CurrentRoll.Sorted()}

	isPCC := p == plyr.PCC
	if isPCC {
		out.numOnMyBar = uint8Ceiling(g.Board.BarCC, maxChexToConsider)
		out.numOnEnemyBar = uint8Ceiling(g.Board.BarC, maxChexToConsider)
	} else {
		out.numOnMyBar = uint8Ceiling(g.Board.BarC, maxChexToConsider)
		out.numOnEnemyBar = uint8Ceiling(g.Board.BarCC, maxChexToConsider)
	}

	out.boardPoints = [constants.NUM_BOARD_POINTS]boardPointState{}
	lastPointIndex := int(constants.NUM_BOARD_POINTS - 1)
	for ptIdx, pt := range g.Board.Points {
		chex := uint8Ceiling(pt.NumCheckers, maxChexToConsider)
		// fill them in order of distance from enemy home. so PCC starts as normal
		translatedPtIdx := lastPointIndex - ptIdx
		if isPCC {
			translatedPtIdx = ptIdx
		}
		out.boardPoints[translatedPtIdx] = boardPointState{pt.Owner == p, chex}
	}

	return out, !isPCC
}
