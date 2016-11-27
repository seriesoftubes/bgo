package game

import (
	"github.com/seriesoftubes/bgo/constants"
)

const (
	numCheckersPerPlayer uint8 = 15
	NUM_BOARD_POINTS     uint8 = 24
	barPips              uint8 = NUM_BOARD_POINTS + 1
	alphabet                   = "abcdefghijklmnopqrstuvwxyz"
)

type BoardPoint struct {
	Owner       *Player
	NumCheckers uint8
}

func (p *BoardPoint) Symbol() string {
	return p.Owner.Symbol()
}

type Board struct {
	Points      [NUM_BOARD_POINTS]*BoardPoint
	BarCC, BarC uint8 // # of checkers on each player's bar
	OffCC, OffC uint8 // # of checkers that each player has beared off
}

func (b *Board) doesPlayerHaveAllRemainingCheckersInHomeBoard(p *Player) bool {
	totalChexInHomeOrBearedOff := b.OffC
	if p == PCC {
		totalChexInHomeOrBearedOff = b.OffCC
	}

	homeStart, homeEnd := p.homePointIndices()
	for i := homeStart; i <= homeEnd; i++ {
		if pt := b.Points[i]; pt.Owner == p {
			totalChexInHomeOrBearedOff += pt.NumCheckers
		}
	}

	return totalChexInHomeOrBearedOff == numCheckersPerPlayer
}

func (b *Board) chexOnTheBar(p *Player) uint8 {
	if p == PC {
		return b.BarC
	}
	return b.BarCC
}

func (b *Board) isLegalMove(m *Move) bool {
	// The player must have no legal moves left, or have used 2 moves, by the end of their turn.
	// TODO: in input parser, parse the whole turn and reject if it doesn't do this

	isForBar := m.Letter == constants.LETTER_BAR_CC || m.Letter == constants.LETTER_BAR_C
	numOnTheBar := b.chexOnTheBar(m.Requestor)
	if numOnTheBar > 0 && !isForBar {
		return false // If you have anything on the bar, you must move those things first
	}
	expectedLetter := constants.LETTER_BAR_C
	if m.Requestor == PCC {
		expectedLetter = constants.LETTER_BAR_CC
	}
	if isForBar && m.Letter != expectedLetter {
		return false // Can't move the enemy's chex.
	}

	numChexOnCurrentPoint := numOnTheBar
	if !isForBar {
		fromPt := b.Points[m.pointIdx()]
		if fromPt.Owner != m.Requestor {
			return false // Can only move your own checkers.
		}
		numChexOnCurrentPoint = fromPt.NumCheckers
	}
	if numChexOnCurrentPoint == 0 {
		return false // Cannot move a checker from an empty point
	}

	nxtIdx, nxtPtExists := m.nextPointIdx()
	if !nxtPtExists {
		if !b.doesPlayerHaveAllRemainingCheckersInHomeBoard(m.Requestor) {
			return false // Can't move past the finish line unless all your remaining checkers are in your home board
		}
		if (m.Requestor == PCC && nxtIdx < 0) || (m.Requestor == PC && nxtIdx >= int8(NUM_BOARD_POINTS)) {
			return false // Must move past the correct finish line.
		}
	}

	if nxtPtExists {
		if nxtPt := b.Points[nxtIdx]; nxtPt.Owner != m.Requestor && nxtPt.NumCheckers > 1 {
			return false // Can't move to a point that's controlled (has >1 chex) by the enemy.
		}
	}

	return true
}

func (b *Board) LegalMoves(p *Player, diceAmt uint8) []*Move {
	var out []*Move

	// Moves off the bar.
	if p == PCC && b.BarCC > 0 {
		m := &Move{Requestor: p, Letter: constants.LETTER_BAR_CC, FowardDistance: diceAmt}
		if b.isLegalMove(m) {
			out = append(out, m)
		}
	} else if p == PC && b.BarC > 0 {
		m := &Move{Requestor: p, Letter: constants.LETTER_BAR_C, FowardDistance: diceAmt}
		if b.isLegalMove(m) {
			out = append(out, m)
		}
	}

	for pointIdx := range b.Points {
		m := &Move{Requestor: p, Letter: string(alphabet[pointIdx]), FowardDistance: diceAmt}
		if b.isLegalMove(m) {
			out = append(out, m)
		}
	}

	return out
}

func (b *Board) incrementBar(p *Player) {
	if p == PCC {
		b.BarCC++
	} else {
		b.BarC++
	}
}

func (b *Board) ExecuteMoveIfLegal(m *Move) bool {
	if !b.isLegalMove(m) {
		return false
	}

	if m.isToMoveSomethingOutOfTheBar() {
		if m.Requestor == PCC {
			b.BarCC--
		} else {
			b.BarC--
		}
	} else {
		fromPt := b.Points[m.pointIdx()]
		fromPt.NumCheckers--
		if fromPt.NumCheckers == 0 {
			fromPt.Owner = nil
		}
	}

	nextPointIdx, nxtPtExists := m.nextPointIdx()
	if !nxtPtExists {
		if m.Requestor == PCC {
			b.OffCC++
		} else {
			b.OffC++
		}
		return true
	}

	nxtPt := b.Points[nextPointIdx]
	if nxtPt.Owner != nil && nxtPt.Owner != m.Requestor {
		nxtPt.NumCheckers--
		b.incrementBar(nxtPt.Owner)
	}
	nxtPt.NumCheckers++
	nxtPt.Owner = m.Requestor

	return true
}

// receive moves like "j1;k3" or "j18;m6". show a preview (with a command + exit command)
// record the move entered so we can undo them. actually just show a preview and
// you can accept the preview, like Y. or no. if yes, update the board's points.
// or just apply the move to the non-copy of it (have a executeMove method, that relies on a
// series of QA checks for whether the move is legit).

func (b *Board) setUp() {
	b.Points = [NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{PCC, 2}, {}, {}, {}, {}, {PC, 5}, {}, {PC, 3}, {}, {}, {}, {PCC, 5},
		{PC, 5}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {PC, 2},
		//                                                        clockwise player in top-left.
	}
}

func (b *Board) PipCounts() (uint, uint) {
	var pipC, pipCC uint

	for i, p := range b.Points {
		basePips, chex := uint(i)+1, uint(p.NumCheckers)
		if p.Owner == PC {
			// the clockwise player's closest checker is at points[0].
			pipC += chex * basePips
		} else if p.Owner == PCC {
			// the counter-clockwise player's furthest checker is at points[0].
			pipCC += chex * (uint(NUM_BOARD_POINTS) - basePips + 1)
		}
	}
	pipC += uint(b.BarC) * uint(barPips)
	pipCC += uint(b.BarCC) * uint(barPips)

	return pipC, pipCC
}
