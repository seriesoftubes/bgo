package game

import (
	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/random"
)

// Combinations of rolls -- not all permutations.
var RollCombos [21]Roll = [21]Roll{
	// one-of-a-kind rolls
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 4},
	{5, 5},
	{6, 6},
	// Below here, there are 2 of each kind of roll.
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 5},
	{1, 6},
	{2, 3},
	{2, 4},
	{2, 5},
	{2, 6},
	{3, 4},
	{3, 5},
	{3, 6},
	{4, 5},
	{4, 6},
	{5, 6},
}

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
