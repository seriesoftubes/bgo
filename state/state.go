package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	lastPointIndex = int(constants.FINAL_BOARD_POINT_INDEX)
	numBoardPoints = lastPointIndex + 1

	numBoardPointVarsForNonCheckerCounts int = 0
	numBoardPointVarsForCheckerCounts    int = 6 // 1c, 2c, 3c, 4c, 5c, 6+c
	numVarsPerBoardPoint                 int = numBoardPointVarsForNonCheckerCounts + numBoardPointVarsForCheckerCounts
	numNonBoardPointVarsPerPlayer        int = 7
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

	for _, player := range []plyr.Player{p, p.Enemy()} {
		// this section adds player-level vars
		barChex, offChex, enemyOff := float32(b.BarC), float32(b.OffC), float32(b.OffCC)
		if player == plyr.PCC {
			barChex, offChex, enemyOff = float32(b.BarCC), float32(b.OffCC), float32(b.OffC)
		}
		slice = append(slice, descBar(player, b, barChex)...)
		slice = append(slice, descOff(offChex, enemyOff)...)

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
	if b.BarCC+b.BarC > 0 {
		return 0.0
	}

	// loop thru points. if you see. PCC -> PC -> PCC, or PC -> PCC -> PC, it's not a race.
	var hasSwitched bool
	var currentPlayer plyr.Player
	for _, pt := range b.Points {
		if p := pt.Owner; p != 0 {
			if currentPlayer != 0 && currentPlayer != p {
				if hasSwitched {
					return 0.0
				}
				hasSwitched = true
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

	// chance of hitting at least one enemy blot (go thru all moves for all 21 possible rolls)

	return subslice
}

func descBar(p plyr.Player, b *game.Board, barChex float32) []float32 {
	enemy := p.Enemy()
	enemyHomeStart, enemyHomeEnd := enemy.HomePointIndices()
	var numEnemyBlots, numLandingPlaces float32
	for ptIdx := enemyHomeStart; ptIdx <= enemyHomeEnd; ptIdx++ {
		if pt := b.Points[ptIdx]; pt.Owner == enemy && pt.NumCheckers == 1 {
			numEnemyBlots++
			numLandingPlaces++
		} else if pt.NumCheckers == 0 || pt.Owner == p {
			numLandingPlaces++
		}
	}
	pctBlot, pctLandable := numEnemyBlots/6.0, numLandingPlaces/6.0

	if barChex > 0 {
		return []float32{1.0, barChex - 1.0, pctBlot, pctLandable}
	}
	return []float32{0.0, 0.0, pctBlot, pctLandable}
}

func descOff(offChex, enemyOff float32) []float32 {
	// TODO: % chance that i will have won after the next move
	//   chex := total chex in home board currently.
	//	 if chex > 4 || not all in home board, chance= 0
	//   elif chex >= 3, chance of getting doubles >= highest point
	//   elif chex > 0, chance of getting >= highest one off + >= lower one.  eg 64, 65, 66 may work.
	//   elif chex == 0, chance=100%
	diffPct := (offChex - enemyOff) / (enemyOff + 1.0/15.0)

	if offChex > 0 {
		return []float32{1.0, offChex - 1.0, diffPct}
	}
	return []float32{0.0, 0.0, diffPct}
}
