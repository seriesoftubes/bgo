package game

import (
	"github.com/seriesoftubes/bgo/random"
)

const (
	minDiceAmt = 1
	maxDiceAmt = 6
)

type Roll [2]uint8

func (r *Roll) moveDistances() []uint8 {
	if first, second := r[0], r[1]; first == second {
		return []uint8{first, first, first, first}
	} else {
		return []uint8{first, second}
	}
}

func (r *Roll) Sorted() Roll {
	if first, second := r[0], r[1]; first <= second {
		return *r
	} else {
		return Roll{second, first}
	}
}

func newRoll() *Roll {
	return &Roll{random.Uint8Between(minDiceAmt, maxDiceAmt), random.Uint8Between(minDiceAmt, maxDiceAmt)}
}
