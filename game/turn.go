package game

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/seriesoftubes/bgo/constants"
)

const moveDelim = ";"

// A Turn contains the moves to execute during a player's turn and the number of times to make each move.
type Turn map[Move]uint8

func (t Turn) update(m Move) { t[m]++ }
func (t Turn) copy() Turn {
	out := make(Turn, len(t))
	for m, times := range t {
		out[m] = times
	}
	return out
}

func (t Turn) totalDist() uint8 {
	var out uint8
	for m, numTimes := range t {
		out += m.FowardDistance * numTimes
	}
	return out
}

// String serializes a Turn into a string like "X;a3;a3;b3;d3".
func (t Turn) String() string {
	var out []string

	first := true
	for m, numTimes := range t {
		if first {
			out = append(out, string(*m.Requestor))
			first = false
		}

		for i := uint8(0); i < numTimes; i++ {
			out = append(out, fmt.Sprintf("%s%d", m.Letter, m.FowardDistance))
		}
	}

	sort.Strings(out)
	return strings.Join(out, moveDelim)
}

// DeserializeTurn creates a Turn from a string like "X;a3;a3;b3;d3".
func DeserializeTurn(s string) (Turn, error) {
	out := Turn{}

	moveStrings := strings.Split(s, moveDelim)

	var p *Player
	if Player(moveStrings[0]) == *PCC {
		p = PCC
	} else if Player(moveStrings[0]) == *PC {
		p = PC
	} else {
		return nil, fmt.Errorf("invalid player in serialized turn %q", s)
	}

	for _, moveString := range moveStrings[1:len(moveStrings)] {
		letter := string(moveString[0])
		dist, err := strconv.Atoi(string(moveString[1]))
		distUint8 := uint8(dist)
		if err != nil || distUint8 < minDiceAmt || distUint8 > maxDiceAmt {
			return nil, fmt.Errorf("invalid distance %v in moveString %v: %v", moveString[1], moveString, err)
		}

		out.update(Move{p, letter, distUint8})
	}

	return out, nil
}

func (t Turn) isValid() bool {
	var p *Player // Placeholder for the first player listed in the turn's moves.
	for m := range t {
		if p == nil {
			p = m.Requestor
		}

		if moveOk, _ := m.isValid(); !moveOk || m.Requestor != p {
			return false
		}
	}

	return true
}

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
func ValidTurns(b *Board, r *Roll, p *Player) map[string]Turn {
	serializedTurns := map[string]Turn{} // set of serialized Turn strings
	var bestTotalDist uint8              // placeholder for the max total distance across all potential turns.
	maybeAddToResultSet := func(t Turn) bool {
		sert := t.String()
		if _, ok := serializedTurns[sert]; ok || !t.isValid() {
			return false // We already processed it, or it's invalid anyway.
		}

		if totalDist := t.totalDist(); totalDist > bestTotalDist {
			bestTotalDist = totalDist
		}
		serializedTurns[sert] = t
		return true
	}

	barLetter := constants.LETTER_BAR_CC
	hasChexOnTheBar := b.BarCC > 0
	if p == PC {
		barLetter = constants.LETTER_BAR_C
		hasChexOnTheBar = b.BarC > 0
	}

	var addPerm func(bb *Board, remainingDists []uint8, t Turn)
	var maybeAddMove func(bcop *Board, m *Move, distIdx int, t Turn, remainingDists []uint8)
	addPerm = func(bb *Board, remainingDists []uint8, t Turn) {
		if remainingDists == nil || len(remainingDists) == 0 {
			return
		}

		for distIdx, dist := range remainingDists {
			if hasChexOnTheBar {
				maybeAddMove(bb.Copy(), &Move{Requestor: p, Letter: barLetter, FowardDistance: dist}, distIdx, t, remainingDists)
			}

			for ptIdx, pt := range bb.Points {
				if pt.Owner == p {
					maybeAddMove(bb.Copy(), &Move{Requestor: p, Letter: constants.Num2Alpha[uint8(ptIdx)], FowardDistance: dist}, distIdx, t, remainingDists)
				}
			}
		}
	}
	maybeAddMove = func(bcop *Board, m *Move, distIdx int, t Turn, remainingDists []uint8) {
		if ok, _ := bcop.ExecuteMoveIfLegal(m); !ok {
			return
		}

		legitTurn := t.copy()
		legitTurn.update(*m)
		if added := maybeAddToResultSet(legitTurn); !added {
			return
		}

		if nextRemaining, err := popSliceUint8(remainingDists, distIdx); err != nil {
			panic("problem popping a value from a uint8 slice: " + err.Error())
		} else {
			addPerm(bcop.Copy(), nextRemaining, legitTurn)
		}
	}

	addPerm(b.Copy(), r.moveDistances(), Turn{})

	for st, t := range serializedTurns {
		if t.totalDist() != bestTotalDist {
			delete(serializedTurns, st)
		}
	}
	return serializedTurns
}
