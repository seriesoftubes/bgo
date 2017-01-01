package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	numVarsPerBoardPoint int = 12
	numNonBoardPointVars int = 3

	lastPointIndex = int(constants.FINAL_BOARD_POINT_INDEX)
)

type State [int(constants.NUM_BOARD_POINTS)*numVarsPerBoardPoint + numNonBoardPointVars]float32

// DetectState detects the current state of the game.
func DetectState(p plyr.Player, b *game.Board) State {
	var out State

	isPCC := p == plyr.PCC
	if isPCC {
		// HeroBar, EnemyBar, Hero-EnemyBar.
		out[0], out[1], out[2] = float32(b.BarCC), float32(b.BarC), float32(b.BarCC-b.BarC)
	} else {
		out[0], out[1], out[2] = float32(b.BarC), float32(b.BarCC), float32(b.BarC-b.BarCC)
	}

	for ptIdx, pt := range b.Points {
		chex := pt.NumCheckers
		// fill them in order of distance from enemy home. so PCC starts as normal
		translatedPtIdx := lastPointIndex - ptIdx
		if isPCC {
			translatedPtIdx = ptIdx
		}
		// the first index in the out array that is relevant to the current boardpoint.
		outIdx := numNonBoardPointVars + (translatedPtIdx * numVarsPerBoardPoint)

		var ownerStatus float32
		if pt.Owner == p {
			ownerStatus = 1.0
		} else if pt.Owner != 0 {
			ownerStatus = -1.0
		}
		out[outIdx+0] = ownerStatus

		numBeyond2 := float32(pt.NumCheckers) - 2.0
		out[outIdx+1] = numBeyond2

		var isSecured float32
		if numBeyond2 >= 0 {
			isSecured = 1.0
		}
		out[outIdx+2] = isSecured

		oppositeDiff := float32(chex - b.Points[constants.FINAL_BOARD_POINT_INDEX-uint8(ptIdx)].NumCheckers)
		out[outIdx+3] = oppositeDiff

		lookaheadDist := int(1)
		var hasBlotDist, hasSecuredDist bool
		var numEnemyChexInFront, distToClosestEnemyBlotPoint, distToClosestEnemySecuredPoint float32
		if isPCC {
			for forwardPtIdx := ptIdx + lookaheadDist; forwardPtIdx < int(constants.NUM_BOARD_POINTS); forwardPtIdx++ {
				if fpt := b.Points[forwardPtIdx]; fpt.Owner == plyr.PC {
					numEnemyChexInFront += float32(fpt.NumCheckers)
					if !hasBlotDist && fpt.NumCheckers == 1 {
						distToClosestEnemyBlotPoint = float32(lookaheadDist)
						hasBlotDist = true
					}
					if !hasSecuredDist && fpt.NumCheckers > 1 {
						distToClosestEnemySecuredPoint = float32(lookaheadDist)
						hasSecuredDist = true
					}
				}
				lookaheadDist++
			}
		} else {
			for forwardPtIdx := ptIdx - lookaheadDist; forwardPtIdx >= 0; forwardPtIdx-- {
				if fpt := b.Points[forwardPtIdx]; fpt.Owner == plyr.PCC {
					numEnemyChexInFront += float32(fpt.NumCheckers)
					if !hasBlotDist && fpt.NumCheckers == 1 {
						distToClosestEnemyBlotPoint = float32(lookaheadDist)
						hasBlotDist = true
					}
					if !hasSecuredDist && fpt.NumCheckers > 1 {
						distToClosestEnemySecuredPoint = float32(lookaheadDist)
						hasSecuredDist = true
					}
				}
				lookaheadDist++
			}
		}
		out[outIdx+4] = numEnemyChexInFront
		out[outIdx+5] = distToClosestEnemyBlotPoint
		out[outIdx+6] = distToClosestEnemySecuredPoint
		out[outIdx+7] = ownerStatus * isSecured
		out[outIdx+8] = ownerStatus * numBeyond2
		out[outIdx+9] = ownerStatus * numEnemyChexInFront
		out[outIdx+10] = ownerStatus * distToClosestEnemyBlotPoint
		out[outIdx+11] = ownerStatus * distToClosestEnemySecuredPoint
	}
	return out
}
