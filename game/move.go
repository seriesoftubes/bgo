package game

import (
	"github.com/seriesoftubes/bgo/constants"
)

var alpha2Num = map[string]uint8{
	"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7, "i": 8, "j": 9,
	"k": 10, "l": 11, "m": 12, "n": 13, "o": 14, "p": 15, "q": 16, "r": 17, "s": 18, "t": 19,
	"u": 20, "v": 21, "w": 22, "x": 23, constants.LETTER_BAR_CC: 24, constants.LETTER_BAR_C: 25,
}

// A move being requested by the current player
type Move struct {
	Requestor      *Player
	Letter         string
	FowardDistance uint8 // validate between 1 and 6
}

func (m *Move) isToMoveSomethingOutOfTheBar() bool {
	return m.Letter == constants.LETTER_BAR_CC || m.Letter == constants.LETTER_BAR_C
}

func (m *Move) pointIdx() uint8 {
	if m.isToMoveSomethingOutOfTheBar() {
		panic("no point index available for the bar letters (those arent stored in board.Points)")
	}

	return alpha2Num[m.Letter]
}

func (m *Move) nextPointIdx() uint8 {
	if m.isToMoveSomethingOutOfTheBar() {
		if m.Requestor == PCC {
			return m.FowardDistance - 1
		}
		return NUM_BOARD_POINTS - m.FowardDistance
	}

	return uint8(int(m.pointIdx()) + m.distCC())
}

// distCC gets the counter-clockwise distance (meaning, moving in a positive direction thru the BoardPoint indices).
// Whatever comes out of this method, you *add* to BoardPoint's index.
func (m *Move) distCC() int {
	if m.Requestor == PCC {
		return int(m.FowardDistance)
	}
	return int(m.FowardDistance) * -1
}
