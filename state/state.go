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

	// the first 24 elements are the normal board points' states.
	// #25 is [numonmybar, numonenemy].
	// #26 is the dice amounts.
	StateArray [constants.NUM_BOARD_POINTS + 2][2]uint8
)

func (bps BoardPointState) AsArray() [2]uint8 {
	if bps.IsOwnedByMe {
		return [2]uint8{1, bps.NumChex}
	}
	return [2]uint8{0, bps.NumChex}
}

func (s State) AsArray() StateArray {
	out := StateArray{}
	for i, bps := range s.BoardPoints {
		out[i] = bps.AsArray()
	}
	out[constants.NUM_BOARD_POINTS] = [2]uint8{s.NumOnMyBar, s.NumOnEnemyBar}
	out[constants.NUM_BOARD_POINTS+1] = [2]uint8(s.MyRoll)
	return out
}

func (s State) InitFromArray(arr StateArray) {
	for i, v := range arr[:constants.NUM_BOARD_POINTS] {
		s.BoardPoints[i] = BoardPointState{v[0] == 1, v[1]}
	}

	barState := arr[constants.NUM_BOARD_POINTS]
	s.NumOnMyBar, s.NumOnEnemyBar = barState[0], barState[1]

	s.MyRoll = game.Roll(arr[constants.NUM_BOARD_POINTS+1])
}

func uint8Ceiling(x, max uint8) uint8 {
	if x > max {
		return max
	}
	return x
}

// Detects the current state of the game, truncating the checker counts up to a max.
// Returns the State and whether the State's BoardPoints were reversed to account for the player's perspective.
func DetectState(p plyr.Player, g *game.Game, maxChexToConsider uint8) (State, bool) {
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
