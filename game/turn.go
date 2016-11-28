package game

import (
	"fmt"
	"sort"

	"github.com/seriesoftubes/bgo/constants"
)

const moveDelim = ";"

// A Turn is an ordered list of moves to execute in that order.
type (
	Turn          []Move
	sortableTurns []Turn
)

func (t Turn) totalDist() uint8 {
	var out uint8
	for _, m := range t {
		out += m.FowardDistance
	}
	return out
}

func (s sortableTurns) Len() int      { return len(s) }
func (s sortableTurns) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less determines whether one turn is less than another, first based on the number of moves, then based on the total move distance.
func (s sortableTurns) Less(i, j int) bool {
	return s[i].totalDist() < s[j].totalDist()
}

func (t Turn) String() string {
	out := ""
	for i, m := range t {
		suffix := moveDelim
		if i == len(t)-1 {
			suffix = ""
		}
		out += fmt.Sprintf("%s%d%s", m.Letter, m.FowardDistance, suffix)
	}
	return out
}

func (t Turn) isValid() bool {
	for _, m := range t {
		if !m.isValid() {
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

// Generates all the best permutations of a roll's distances.
func TurnPerms(b *Board, r *Roll, p *Player) []Turn {
	var perms sortableTurns
	appendToPermsIfValid := func(t Turn) {
		if t.isValid() {
			perms = append(perms, t)
		}
	}

	barLetter := constants.LETTER_BAR_CC
	if p == PC {
		barLetter = constants.LETTER_BAR_C
	}

	var addPerm func(bb *Board, remainingDists []uint8, t Turn)
	addPerm = func(bb *Board, remainingDists []uint8, t Turn) {
		if remainingDists == nil || len(remainingDists) == 0 {
			return
		}

		for distIdx, dist := range remainingDists {
			cop := bb.Copy()
			m := &Move{Requestor: p, Letter: barLetter, FowardDistance: dist}
			if ok := cop.ExecuteMoveIfLegal(m); ok {
				legitTurn := append(t, *m)
				appendToPermsIfValid(legitTurn)

				nextRemaining, _ := popSliceUint8(remainingDists, distIdx) // Guaranteed to be no error
				addPerm(cop.Copy(), nextRemaining, legitTurn)
			}

			for ptIdx, pt := range bb.Points {
				if pt.Owner != p {
					continue
				}

				cop := bb.Copy()
				m := &Move{Requestor: p, Letter: constants.Num2Alpha[uint8(ptIdx)], FowardDistance: dist}
				if ok := cop.ExecuteMoveIfLegal(m); ok {
					legitTurn := append(t, *m)
					appendToPermsIfValid(legitTurn)

					nextRemaining, _ := popSliceUint8(remainingDists, distIdx) // Guaranteed to be no error
					addPerm(cop.Copy(), nextRemaining, legitTurn)
				}
			}
		}
	}

	addPerm(b.Copy(), r.moveDistances(), Turn(nil))

	if len(perms) == 0 {
		return perms
	}

	var viableTurns []Turn
	sort.Sort(sort.Reverse(perms))
	bestDist := perms[0].totalDist()
	for _, perm := range perms {
		if perm.totalDist() == bestDist {
			viableTurns = append(viableTurns, perm)
		}
	}
	return viableTurns
}
