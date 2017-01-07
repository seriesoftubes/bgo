package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	numBoardPointVarsForNonCheckerCounts = 0
	numBoardPointVarsForCheckerCounts    = 6 // 1c, 2c, 3c, 4c, 5c, 6+c
	numVarsPerBoardPoint                 = numBoardPointVarsForNonCheckerCounts + numBoardPointVarsForCheckerCounts
	numNonBoardPointVars                 = 2
	stateLength                          = constants.NUM_PLAYERS * (numNonBoardPointVars + int(constants.NUM_BOARD_POINTS)*numVarsPerBoardPoint)

	lastPointIndex = int(constants.FINAL_BOARD_POINT_INDEX)

	tinyIncrement = float32(0.01)
)

type State [stateLength]float32

// DetectState detects the current state of the game.
func DetectState(p plyr.Player, b *game.Board) State {
	isPCC := p == plyr.PCC
	slice := make([]float32, 0, stateLength)
	for _, player := range []plyr.Player{p, p.Enemy()} {
		barChex, offChex := float32(b.BarC), float32(b.OffC)
		if isPCC {
			barChex, offChex = float32(b.BarCC), float32(b.OffCC)
		}
		onChex := float32(constants.NUM_CHECKERS_PER_PLAYER) - offChex
		slice = append(slice,
			barChex/float32(2.0),
			offChex/(onChex+tinyIncrement),
		)

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

	return subslice
}
