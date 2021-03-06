package turn

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	moveDelim = ";"
)

type (
	// A Turn contains the moves to execute during a player's turn and the number of times to make each move.
	Turn      map[Move]uint8
	TurnArray [constants.MAX_MOVES_PER_TURN]MoveArray
)

type (
	sortableMoves []Move
)

// DeserializeTurn creates a Turn from a string like "X;a3;a3;b3;d3".
func DeserializeTurn(s string) (Turn, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == moveDelim {
		return nil, fmt.Errorf("no turn provided in serialized turn %q", s)
	}

	out := Turn{}

	moveStrings := strings.Split(s, moveDelim)
	if len(moveStrings) < 2 {
		return nil, fmt.Errorf("invalid format of serialized turn %q", s)
	}
	playerString := moveStrings[0]
	if len(playerString) != 1 {
		return nil, fmt.Errorf("invalid format of serialized turn %q", s)
	}

	var p plyr.Player
	if plyr.Player(playerString[0]) == plyr.PCC {
		p = plyr.PCC
	} else if plyr.Player(playerString[0]) == plyr.PC {
		p = plyr.PC
	} else {
		return nil, fmt.Errorf("invalid player in serialized turn %q", s)
	}

	for _, moveString := range moveStrings[1:len(moveStrings)] {
		if len(moveString) != 2 {
			return nil, fmt.Errorf("invalid part %q of serialized turn %q", moveString, s)
		}

		letter := moveString[0]
		dist, err := strconv.Atoi(string(moveString[1]))
		distUint8 := uint8(dist)
		if err != nil || distUint8 < constants.MIN_DICE_AMT || distUint8 > constants.MAX_DICE_AMT {
			return nil, fmt.Errorf("invalid distance %v in moveString %v: %v", moveString[1], moveString, err)
		}

		out.Update(Move{p, letter, distUint8})
	}

	return out, nil
}

func (t Turn) Update(m Move) { t[m]++ }
func (t Turn) Copy() Turn {
	out := make(Turn, len(t))
	for m, times := range t {
		out[m] = times
	}
	return out
}

func (t Turn) TotalDist() uint8 {
	var out uint8
	for m, numTimes := range t {
		out += m.FowardDistance * numTimes
	}
	return out
}

func (t Turn) Arrayify() TurnArray {
	smoves := make(sortableMoves, 0, len(t))
	for m := range t {
		smoves = append(smoves, m)
	}
	sort.Sort(smoves)

	ta := TurnArray{}
	for i, mov := range smoves {
		ta[i] = mov.arrayify(t[mov])
	}
	return ta
}

// String serializes a Turn into a string like "X;a3;a3;b3;d3".
func (t Turn) String() string {
	if len(t) == 0 {
		return ""
	}

	var out []string

	smoves := make(sortableMoves, 0, len(t))
	for m := range t {
		smoves = append(smoves, m)
	}
	sort.Sort(smoves)
	out = append(out, string(smoves[0].Requestor))

	for _, mov := range smoves {
		if numTimes := int(t[mov]); numTimes == 1 {
			out = append(out, fmt.Sprintf("%s%d", string(mov.Letter), mov.FowardDistance))
		} else {
			reps := strings.Repeat(fmt.Sprintf("%s%d;", string(mov.Letter), mov.FowardDistance), numTimes)
			out = append(out, reps[0:len(reps)-1])
		}
	}

	return strings.Join(out, moveDelim)
}

func (t Turn) IsValid() bool {
	var p plyr.Player // Placeholder for the first player listed in the turn's moves.
	for m := range t {
		if p == 0 {
			p = m.Requestor
		}

		if moveOk, _ := m.IsValid(); !moveOk || m.Requestor != p {
			return false
		}
	}

	return true
}

func (sm sortableMoves) Len() int      { return len(sm) }
func (sm sortableMoves) Swap(i, j int) { sm[j], sm[i] = sm[i], sm[j] }
func (sm sortableMoves) Less(i, j int) bool {
	if left, right := sm[i], sm[j]; left.Letter < right.Letter {
		return true
	} else if left.Letter > right.Letter {
		return false
	} else {
		return left.FowardDistance < right.FowardDistance
	}
}
