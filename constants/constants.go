// Package constants contains variables that never change.
package constants

const (
	NUM_POINTS_IN_HOME_BOARD uint8 = 6
	NUM_BOARD_POINTS         uint8 = 24
	NUM_CHECKERS_PER_PLAYER  uint8 = 15
	MIN_DICE_AMT                   = 1
	MAX_DICE_AMT                   = 6

	LETTER_BAR_CC = "y" // Accesses chex on the bar for the CC player
	LETTER_BAR_C  = "z" // Accesses chex on the bar for the C player
)

var (
	Num2Alpha = map[uint8]string{
		0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h", 8: "i", 9: "j",
		10: "k", 11: "l", 12: "m", 13: "n", 14: "o", 15: "p", 16: "q", 17: "r", 18: "s", 19: "t",
		20: "u", 21: "v", 22: "w", 23: "x", 24: LETTER_BAR_CC, 25: LETTER_BAR_C,
	}
	Alpha2Num = map[string]uint8{
		"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7, "i": 8, "j": 9,
		"k": 10, "l": 11, "m": 12, "n": 13, "o": 14, "p": 15, "q": 16, "r": 17, "s": 18, "t": 19,
		"u": 20, "v": 21, "w": 22, "x": 23, LETTER_BAR_CC: 24, LETTER_BAR_C: 25,
	}
)
