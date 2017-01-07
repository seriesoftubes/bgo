package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	lastPointIndex = int(constants.FINAL_BOARD_POINT_INDEX)
	numBoardPoints = lastPointIndex + 1

	tinyIncrement = float32(0.01)

	numBoardPointVarsForNonCheckerCounts int = 0
	numBoardPointVarsForCheckerCounts    int = 6 // 1c, 2c, 3c, 4c, 5c, 6+c
	numVarsPerBoardPoint                 int = numBoardPointVarsForNonCheckerCounts + numBoardPointVarsForCheckerCounts
	numNonBoardPointVarsPerPlayer        int = 2
	numNonPlayerSpecificVars             int = 1
	stateLength                          int = numNonPlayerSpecificVars + constants.NUM_PLAYERS*(numNonBoardPointVarsPerPlayer+numBoardPoints*numVarsPerBoardPoint)
)

type State [stateLength]float32

// DetectState detects the current state of the game.
func DetectState(p plyr.Player, b *game.Board) State {
	isPCC := p == plyr.PCC
	slice := make([]float32, 0, stateLength)

	// non player-specific vars
	slice = append(slice, isRace(b))

	// player-specific vars in here
	for _, player := range []plyr.Player{p, p.Enemy()} {
		// this section adds player-level vars
		barChex, offChex := float32(b.BarC), float32(b.OffC)
		if isPCC {
			barChex, offChex = float32(b.BarCC), float32(b.OffCC)
		}
		onChex := float32(constants.NUM_CHECKERS_PER_PLAYER) - offChex
		slice = append(slice,
			barChex/float32(2.0),
			offChex/(onChex+tinyIncrement),
			// "behind" means , look at the furthest forward enemy point. how many hero chex are behind it
			// has 1 behind enemy lines. 2. 3. 4+
		)

		// this section adds boardPoint-specific vars for each player.
		if isPCC {
			for i := 0; i <= lastPointIndex; i++ {
				slice = append(slice, descPoint(b.Points[i], player)...)
			}
		} else {
			for i := lastPointIndex; i >= 0; i-- {
				slice = append(slice, descPoint(b.Points[i], player)...)
			}
		}
	}

	var out State
	for i, v := range slice {
		out[i] = v
	}
	return out
}

func isRace(b *game.Board) float32 {
	// loop thru points. if you see. PCC -> PC -> PCC, or PC -> PCC -> PC, it's not a race.
	var numPlayerTransitions uint8
	var currentPlayer plyr.Player
	for _, pt := range b.Points {
		if p := pt.Owner; p != 0 {
			if currentPlayer != 0 && currentPlayer != p {
				numPlayerTransitions++
				if numPlayerTransitions > 2 {
					return 0.0
				}
			}
			currentPlayer = p
		}
	}

	return 1.0
}

func descPoint(pt *game.BoardPoint, supposedOwner plyr.Player) []float32 {
	subslice := make([]float32, numVarsPerBoardPoint)

	if pt.Owner != supposedOwner { // we only ever want to describe a point owned by the currently analyzed player.
		return subslice
	}

	for ct := uint8(1); ct <= pt.NumCheckers; ct++ {
		cappedCt := int(ct)
		if cappedCt > numBoardPointVarsForCheckerCounts {
			cappedCt = numBoardPointVarsForCheckerCounts
		}
		subslice[cappedCt-1]++
	}

	// chance of moving forward at all.
	//   look at 1, 2, 3, 4, 5, 6 things ahead. calc chance of rolling (1,2,3...etc)
	// chance of hitting at least one enemy blot

	return subslice
}
