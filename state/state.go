package state

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turngen"
)

const (
	numRollPermutations = float32(36.0)

	lastPointIndex = int(constants.FINAL_BOARD_POINT_INDEX)
	numBoardPoints = lastPointIndex + 1

	numBoardPointVarsForNonCheckerCounts int = 0
	numBoardPointVarsForCheckerCounts    int = 6 // 1c, 2c, 3c, 4c, 5c, 6+c
	numVarsPerBoardPoint                 int = numBoardPointVarsForNonCheckerCounts + numBoardPointVarsForCheckerCounts
	numNonBoardPointVarsPerPlayer        int = 9
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
		barChex, offChex, enemyBar, enemyOff := float32(b.BarC), float32(b.OffC), float32(b.BarCC), float32(b.OffCC)
		if player == plyr.PCC {
			barChex, offChex, enemyBar, enemyOff = float32(b.BarCC), float32(b.OffCC), float32(b.BarC), float32(b.OffC)
		}
		slice = append(slice, descBar(player, b, barChex)...)
		slice = append(slice, descOff(offChex, enemyOff)...)
		slice = append(slice, descPlayerBoard(b, player, enemyBar)...)

		// this section adds boardPoint-specific vars for each player.
		if isPCC {
			for i := 0; i <= lastPointIndex; i++ {
				slice = append(slice, descPoint(b, i, player)...)
			}
		} else {
			for i := lastPointIndex; i >= 0; i-- {
				slice = append(slice, descPoint(b, i, player)...)
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

func descPlayerBoard(b *game.Board, player plyr.Player, enemyBar float32) []float32 {
	var numWinningRolls, numBlotHittingRolls, totalEnemyBlotsHit float32
	for _, roll := range game.RollCombos {
		numRollOccurrences := float32(2.0)
		if isDoubles := roll[0] == roll[1]; isDoubles {
			numRollOccurrences = float32(1.0)
		}

		for _, vt := range turngen.ValidTurns(b, roll, player) {
			bc := b.Copy()
			bc.MustExecuteTurn(vt, false)
			numEnemyBlotsHit, heroOffCount := float32(bc.BarCC)-enemyBar, bc.OffC
			if player == plyr.PCC {
				numEnemyBlotsHit, heroOffCount = float32(bc.BarC)-enemyBar, bc.OffCC
			}
			totalEnemyBlotsHit += (numEnemyBlotsHit * numRollOccurrences)
			if numEnemyBlotsHit > 0 {
				numBlotHittingRolls += numRollOccurrences
			}
			if heroOffCount == constants.NUM_CHECKERS_PER_PLAYER {
				numWinningRolls += numRollOccurrences
			}
		}
	}
	chanceWin := numWinningRolls / numRollPermutations
	chanceHitBlot := numBlotHittingRolls / numRollPermutations
	expectedBlotsHit := totalEnemyBlotsHit / numRollPermutations

	return []float32{chanceWin, chanceHitBlot, expectedBlotsHit}
}

func descPoint(b *game.Board, pointIdx int, supposedOwner plyr.Player) []float32 {
	subslice := make([]float32, numVarsPerBoardPoint)
	pt := b.Points[pointIdx]

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
	if offChex > 0 {
		return []float32{1.0, offChex - 1.0}
	}
	return []float32{0.0, 0.0}
}
