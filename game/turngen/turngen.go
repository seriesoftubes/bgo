// Package turngen generates valid Turns based on the state of a board.
package turngen

import (
	"fmt"

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
func ValidTurns(b *game.Board, r *game.Roll, p *plyr.Player) map[turn.TurnArray]turn.Turn {
	serializedTurns := map[turn.TurnArray]turn.Turn{} // set of serialized Turn strings
	var bestTotalDist uint8                           // placeholder for the max total distance across all potential turns.
	maybeAddToResultSet := func(t turn.Turn) bool {
		sert := t.Arrayify()
		if _, ok := serializedTurns[sert]; ok || !t.IsValid() {
			return false // We already processed it, or it's invalid anyway.
		}

		if totalDist := t.TotalDist(); totalDist > bestTotalDist {
			bestTotalDist = totalDist
		}
		serializedTurns[sert] = t
		return true
	}

	var addPerm func(bb *game.Board, remainingDists []uint8, t turn.Turn)
	addPerm = func(bb *game.Board, remainingDists []uint8, t turn.Turn) {
		if remainingDists == nil || len(remainingDists) == 0 {
			return
		}

		for distIdx, dist := range remainingDists {
			for _, mv := range bb.LegalMoves(p, dist) {
				bcop := bb.Copy()
				bcop.ExecuteMoveUnsafe(mv) // We already know the move is legal, so it's safe to do it.

				legitTurn := t.Copy()
				legitTurn.Update(*mv)
				if !maybeAddToResultSet(legitTurn) {
					continue
				}

				if nextRemaining, err := popSliceUint8(remainingDists, distIdx); err != nil {
					panic("problem popping a value from a uint8 slice: " + err.Error())
				} else {
					addPerm(bcop, nextRemaining, legitTurn)
				}
			}
		}
	}
	addPerm(b, r.MoveDistances(), turn.Turn{})

	for st, t := range serializedTurns {
		if t.TotalDist() != bestTotalDist {
			delete(serializedTurns, st)
		}
	}
	return serializedTurns
}
