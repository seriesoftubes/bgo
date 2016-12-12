package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

type (
	BoardPointState struct {
		IsOwnedByMe bool
		NumChex     uint8
	}

	State struct {
		// points on the board, indexed with 0 being the furthest away from the current player's home
		BoardPoints               [constants.NUM_BOARD_POINTS]BoardPointState
		NumOnMyBar, NumOnEnemyBar uint8
		MyRoll                    game.Roll
	}
)

func uint8Ceiling(x, max uint8) uint8 {
	if x > max {
		return max
	}
	return x
}

// Detects the current state of the game, truncating the checker counts up to a max.
// Returns the State and whether the State's BoardPoints were reversed to account for the player's perspective.
func DetectState(p *plyr.Player, g *game.Game, maxChexToConsider uint8) (State, bool) {
	out := State{MyRoll: g.CurrentRoll.Sorted()}

	isPCC := p == plyr.PCC
	if isPCC {
		out.NumOnMyBar = uint8Ceiling(g.Board.BarCC, maxChexToConsider)
		out.NumOnEnemyBar = uint8Ceiling(g.Board.BarC, maxChexToConsider)
	} else {
		out.NumOnMyBar = uint8Ceiling(g.Board.BarC, maxChexToConsider)
		out.NumOnEnemyBar = uint8Ceiling(g.Board.BarCC, maxChexToConsider)
	}

	out.BoardPoints = [constants.NUM_BOARD_POINTS]BoardPointState{}
	lastPointIndex := int(constants.NUM_BOARD_POINTS - 1)
	for ptIdx, pt := range g.Board.Points {
		chex := uint8Ceiling(pt.NumCheckers, maxChexToConsider)
		// fill them in order of distance from enemy home. so PCC starts as normal
		translatedPtIdx := lastPointIndex - ptIdx
		if isPCC {
			translatedPtIdx = ptIdx
		}
		out.BoardPoints[translatedPtIdx] = BoardPointState{pt.Owner == p, chex}
	}

	return out, !isPCC
}
