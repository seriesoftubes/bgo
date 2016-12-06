package game

import (
	"github.com/seriesoftubes/bgo/constants"
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

func (r *Roll) reverse() *Roll {
	r[1], r[0] = r[0], r[1]
	return r
}

func randBetween(min, max int) uint8 { return uint8(constants.Rand.Intn(max-min+1) + min) }
func newRoll() *Roll {
	return &Roll{randBetween(minDiceAmt, maxDiceAmt), randBetween(minDiceAmt, maxDiceAmt)}
}
