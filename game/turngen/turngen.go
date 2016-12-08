// Package turngen generates valid Turns based on the state of a board.
package turngen

import (
	"fmt"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
)

// TODO: Use go generate for this.
func copySliceUint8(slice []uint8) []uint8 { return append([]uint8(nil), slice...) }
func popSliceUint8(slice []uint8, atIndex int) ([]uint8, error) {
	if slice == nil || len(slice) == 0 {
		return nil, fmt.Errorf("cannot pop from an empty slice")
	} else if atIndex >= len(slice) || atIndex < 0 {
		return nil, fmt.Errorf("index %d out of bounds, must be between [0, %d] inclusive", atIndex, len(slice)-1)
	}
	slice = copySliceUint8(slice)
	return append(slice[:atIndex], slice[atIndex+1:]...), nil
}

// Generates the set of all valid turns for a player, given a roll and a board.
func ValidTurns(b *game.Board, r *game.Roll, p *plyr.Player) map[string]turn.Turn {
	serializedTurns := map[string]turn.Turn{} // set of serialized Turn strings
	var bestTotalDist uint8                   // placeholder for the max total distance across all potential turns.
	maybeAddToResultSet := func(t turn.Turn) bool {
		sert := t.String()
		if _, ok := serializedTurns[sert]; ok || !t.IsValid() {
			return false // We already processed it, or it's invalid anyway.
		}

		if totalDist := t.TotalDist(); totalDist > bestTotalDist {
			bestTotalDist = totalDist
		}
		serializedTurns[sert] = t
		return true
	}

	barLetter := constants.LETTER_BAR_CC
	hasChexOnTheBar := b.BarCC > 0
	if p == plyr.PC {
		barLetter = constants.LETTER_BAR_C
		hasChexOnTheBar = b.BarC > 0
	}

	var addPerm func(bb *game.Board, remainingDists []uint8, t turn.Turn)
	var maybeAddMove func(bcop *game.Board, m *turn.Move, distIdx int, t turn.Turn, remainingDists []uint8)
	addPerm = func(bb *game.Board, remainingDists []uint8, t turn.Turn) {
		if remainingDists == nil || len(remainingDists) == 0 {
			return
		}

		for distIdx, dist := range remainingDists {
			if hasChexOnTheBar {
				maybeAddMove(bb.Copy(), &turn.Move{Requestor: p, Letter: barLetter, FowardDistance: dist}, distIdx, t, remainingDists)
			}

			for ptIdx, pt := range bb.Points {
				if pt.Owner == p {
					maybeAddMove(bb.Copy(), &turn.Move{Requestor: p, Letter: constants.Num2Alpha[uint8(ptIdx)], FowardDistance: dist}, distIdx, t, remainingDists)
				}
			}
		}
	}
	maybeAddMove = func(bcop *game.Board, m *turn.Move, distIdx int, t turn.Turn, remainingDists []uint8) {
		if ok, _ := bcop.ExecuteMoveIfLegal(m); !ok {
			return
		}

		legitTurn := t.Copy()
		legitTurn.Update(*m)
		if added := maybeAddToResultSet(legitTurn); !added {
			return
		}

		if nextRemaining, err := popSliceUint8(remainingDists, distIdx); err != nil {
			panic("problem popping a value from a uint8 slice: " + err.Error())
		} else {
			addPerm(bcop.Copy(), nextRemaining, legitTurn)
		}
	}

	addPerm(b.Copy(), r.MoveDistances(), turn.Turn{})

	for st, t := range serializedTurns {
		if t.TotalDist() != bestTotalDist {
			delete(serializedTurns, st)
		}
	}
	return serializedTurns
}
