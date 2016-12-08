package turn

import (
	"fmt"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game/plyr"
)

// A move being requested by the current player
type Move struct {
	Requestor      *plyr.Player
	Letter         string
	FowardDistance uint8 // validate between 1 and 6
}

func (m *Move) String() string {
	return fmt.Sprintf("%s: go %d spaces starting with %s", *m.Requestor, m.FowardDistance, m.Letter)
}

const (
	invalidNilRequestor = "Must have requestor."
	invalidBadChar      = "Must be a whitelisted lowercase alpha character."
	invalidBadDice      = "Distance must be between [1,6]"
)

func (m *Move) IsValid() (bool, string) {
	if m.Requestor == nil {
		return false, invalidNilRequestor
	}

	if _, ok := constants.Alpha2Num[m.Letter]; !ok {
		return false, invalidBadChar
	}

	if m.FowardDistance < constants.MIN_DICE_AMT || m.FowardDistance > constants.MAX_DICE_AMT {
		return false, invalidBadDice
	}

	return true, ""
}

func (m *Move) IsToMoveSomethingOutOfTheBar() bool {
	return m.Letter == constants.LETTER_BAR_CC || m.Letter == constants.LETTER_BAR_C
}

func (m *Move) PointIdx() uint8 {
	if m.IsToMoveSomethingOutOfTheBar() {
		panic("no point index available for the bar letters (those arent stored in board.Points)")
	}

	return constants.Alpha2Num[m.Letter]
}

// TODO: MoveSet struct, or execute each command, 1 at a time, until the player is out of moves.

// Gets the PointIndex of the next point to go to.
// May return < 0 or > 23, if the move is to bear-off a checker.
func (m *Move) NextPointIdx() (int8, bool) {
	var nxtIdx int8
	if m.IsToMoveSomethingOutOfTheBar() {
		if m.Requestor == plyr.PCC {
			nxtIdx = int8(m.FowardDistance) - 1
		} else {
			nxtIdx = int8(constants.NUM_BOARD_POINTS - m.FowardDistance)
		}
	} else {
		nxtIdx = int8(int(m.PointIdx()) + m.distCC())
	}

	return nxtIdx, nxtIdx >= 0 && nxtIdx < int8(constants.NUM_BOARD_POINTS)
}

// distCC gets the counter-clockwise distance (meaning, moving in a positive direction thru the BoardPoint indices).
// Whatever comes out of this method, you *add* to BoardPoint's index.
func (m *Move) distCC() int {
	if m.Requestor == plyr.PCC {
		return int(m.FowardDistance)
	}
	return int(m.FowardDistance) * -1
}
