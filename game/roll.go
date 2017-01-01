package game

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/random"
)

type Roll [2]uint8

func newRoll() Roll {
	return Roll{random.Uint8Between(constants.MIN_DICE_AMT, constants.MAX_DICE_AMT), random.Uint8Between(constants.MIN_DICE_AMT, constants.MAX_DICE_AMT)}
}

func (r *Roll) MoveDistances() []uint8 {
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
