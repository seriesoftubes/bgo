// Package turngen generates valid Turns based on the state of a board.
package turngen

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
)

// Generates the set of all valid turns for a player, given a roll and a board.
func ValidTurns(b *game.Board, r game.Roll, p plyr.Player) map[turn.TurnArray]turn.Turn {
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
	if isDoubles := r[0] == r[1]; isDoubles {
		addPerm = func(bb *game.Board, remainingDists []uint8, t turn.Turn) {
			for _, mv := range bb.LegalMoves(p, remainingDists[0]) {
				bcop := bb.Copy()
				bcop.ExecuteMoveUnsafe(mv) // We already know the move is legal, so it's safe to do it.

				legitTurn := t.Copy()
				legitTurn.Update(mv)
				if !maybeAddToResultSet(legitTurn) {
					continue
				}

				if nextRemaining, err := popSliceUint8(remainingDists, 0); err != nil {
					panic("problem popping a value from a uint8 slice: " + err.Error())
				} else if nextRemaining != nil {
					addPerm(bcop, nextRemaining, legitTurn)
				}
			}
		}
	} else {
		addPerm = func(bb *game.Board, remainingDists []uint8, t turn.Turn) {
			for distIdx, dist := range remainingDists {
				for _, mv := range bb.LegalMoves(p, dist) {
					bcop := bb.Copy()
					bcop.ExecuteMoveUnsafe(mv) // We already know the move is legal, so it's safe to do it.

					legitTurn := t.Copy()
					legitTurn.Update(mv)
					if !maybeAddToResultSet(legitTurn) {
						continue
					}

					if nextRemaining, err := popSliceWithOneOrTwoElements(remainingDists, distIdx); err != nil {
						panic("problem popping a value from a uint8 slice: " + err.Error())
					} else if nextRemaining != nil {
						addPerm(bcop, nextRemaining, legitTurn)
					}
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

func copySliceUint8(slice []uint8) []uint8 { return append([]uint8(nil), slice...) }
func popSliceUint8(slice []uint8, atIndex int) ([]uint8, error) {
	if sz := len(slice); sz == 1 {
		return nil, nil
	} else if sz == 0 {
		return nil, fmt.Errorf("cannot pop from an empty slice")
	} else if atIndex >= sz || atIndex < 0 {
		return nil, fmt.Errorf("index %d out of bounds, must be between [0, %d] inclusive", atIndex, sz-1)
	}

	slice = copySliceUint8(slice)
	return append(slice[:atIndex], slice[atIndex+1:]...), nil
}

func popSliceWithOneOrTwoElements(sliceWithOneOrTwoElements []uint8, atIndex int) ([]uint8, error) {
	if sz := len(sliceWithOneOrTwoElements); sz == 1 {
		return nil, nil
	} else if sz == 2 {
		return []uint8{sliceWithOneOrTwoElements[1-atIndex]}, nil
	}

	return nil, fmt.Errorf("index %d out of bounds, must be between [0, %d] inclusive", atIndex, len(sliceWithOneOrTwoElements)-1)
}
