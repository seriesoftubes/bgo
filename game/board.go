package game

import (
	"fmt"
	"sort"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
)

const (
	WinKindNotWon     WinKind = 0
	WinKindSingleGame WinKind = 1
	WinKindGammon     WinKind = 2
	WinKindBackgammon WinKind = 3
)

const (
	barPips = uint16(constants.NUM_BOARD_POINTS + 1)

	illegalBarFirst                   = "If you have anything on the bar, you must move those things first"
	illegalEnemyBarChex               = "Can't move the enemy's chex off the enemy's bar."
	illegalEnemyRegularChex           = "Can only move your own checkers."
	illegalEmptyPoint                 = "Cannot move a checker from an empty point"
	illegalCantBearoffUntilAllAreHome = "Can't move past the finish line unless all your remaining checkers are in your home board"
	illegalWrongFinishLine            = "Must move past the correct finish line."
	illegalBearoffOthersFirst         = "If the amount on the dice > the point's distance away from 0, then you must have already beared off all chex behind the point."
	illegalEnemyControlsIt            = "Can't move to a point that's controlled (has >1 chex) by the enemy."
	illegalBearOntoEnemyHome          = "Must move to a point inside the enemy's home."
)

type (
	WinKind uint8

	BoardPoint struct {
		Owner       plyr.Player
		NumCheckers uint8
	}
	Board struct {
		Points      *[constants.NUM_BOARD_POINTS]*BoardPoint
		BarCC, BarC uint8 // # of checkers on each player's bar
		OffCC, OffC uint8 // # of checkers that each player has beared off
		// These win-related fields must only be set by the board itself.
		winner  plyr.Player
		winKind WinKind
	}
)

type (
	motimesPair struct {
		mo    turn.Move
		times uint8
	}
	sortableMotimesPairs []motimesPair
)

func detectWinKind(b *Board, p plyr.Player) WinKind {
	otherPlayer := plyr.PC
	numOtherPlayerHasBearedOff := b.OffC

	if p == plyr.PCC {
		if b.OffCC != constants.NUM_CHECKERS_PER_PLAYER {
			return WinKindNotWon
		}
	} else {
		if b.OffC != constants.NUM_CHECKERS_PER_PLAYER {
			return WinKindNotWon
		}

		otherPlayer = plyr.PCC
		numOtherPlayerHasBearedOff = b.OffCC
	}

	if numOtherPlayerHasBearedOff > 0 {
		return WinKindSingleGame
	}

	homeStart, homeEnd := p.HomePointIndices()
	for i := homeStart; i <= homeEnd; i++ {
		if pt := b.Points[i]; pt.Owner == otherPlayer {
			return WinKindBackgammon
		}
	}

	return WinKindGammon
}

func (smp sortableMotimesPairs) Len() int      { return len(smp) }
func (smp sortableMotimesPairs) Swap(i, j int) { smp[i], smp[j] = smp[j], smp[i] }
func (smp sortableMotimesPairs) Less(i, j int) bool {
	if left, right := smp[i], smp[j]; left.mo.Requestor == plyr.PCC {
		return left.mo.Letter < right.mo.Letter // PCC needs to exec lo letters first, then hi ones
	} else {
		return left.mo.Letter > right.mo.Letter // PC needs to exec hi letters first, then lo ones
	}
}

func (p *BoardPoint) Symbol() string { return p.Owner.Symbol() }

// Copy returns a pointer to a deepcopy of a Board.
func (b *Board) Copy() *Board {
	cop := &Board{}
	pts := b.Points
	cop.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		{pts[0].Owner, pts[0].NumCheckers},
		{pts[1].Owner, pts[1].NumCheckers},
		{pts[2].Owner, pts[2].NumCheckers},
		{pts[3].Owner, pts[3].NumCheckers},
		{pts[4].Owner, pts[4].NumCheckers},
		{pts[5].Owner, pts[5].NumCheckers},
		{pts[6].Owner, pts[6].NumCheckers},
		{pts[7].Owner, pts[7].NumCheckers},
		{pts[8].Owner, pts[8].NumCheckers},
		{pts[9].Owner, pts[9].NumCheckers},
		{pts[10].Owner, pts[10].NumCheckers},
		{pts[11].Owner, pts[11].NumCheckers},
		{pts[12].Owner, pts[12].NumCheckers},
		{pts[13].Owner, pts[13].NumCheckers},
		{pts[14].Owner, pts[14].NumCheckers},
		{pts[15].Owner, pts[15].NumCheckers},
		{pts[16].Owner, pts[16].NumCheckers},
		{pts[17].Owner, pts[17].NumCheckers},
		{pts[18].Owner, pts[18].NumCheckers},
		{pts[19].Owner, pts[19].NumCheckers},
		{pts[20].Owner, pts[20].NumCheckers},
		{pts[21].Owner, pts[21].NumCheckers},
		{pts[22].Owner, pts[22].NumCheckers},
		{pts[23].Owner, pts[23].NumCheckers},
	}

	cop.BarC, cop.BarCC = b.BarC, b.BarCC
	cop.OffC, cop.OffCC = b.OffC, b.OffCC
	cop.winner, cop.winKind = b.winner, b.winKind

	return cop
}

func (b *Board) Winner() plyr.Player { return b.winner }
func (b *Board) WinKind() WinKind    { return b.winKind }

func (b *Board) LegalMoves(p plyr.Player, diceAmt uint8) []turn.Move {
	var out []turn.Move

	if p == plyr.PCC && b.BarCC > 0 {
		m := turn.Move{Requestor: p, Letter: constants.LETTER_BAR_CC, FowardDistance: diceAmt}
		if ok, _ := b.isLegalMoveForBearingOn(m); ok {
			return append(out, m)
		}
		return out
	} else if p == plyr.PC && b.BarC > 0 {
		m := turn.Move{Requestor: p, Letter: constants.LETTER_BAR_C, FowardDistance: diceAmt}
		if ok, _ := b.isLegalMoveForBearingOn(m); ok {
			return append(out, m)
		}
		return out
	}

	for pointIdx, pt := range b.Points {
		if pt.Owner != p {
			continue
		}

		m := turn.Move{Requestor: p, Letter: constants.Num2Alpha[uint8(pointIdx)], FowardDistance: diceAmt}
		if ok, _ := b.isLegalMoveForNonBearingOn(m); ok {
			out = append(out, m)
		}
	}

	return out
}

// MustExecuteTurn takes a Turn, and executes its individual moves, in an order that won't explode the game.
// This is mainly to support the stdin UX of supplying entire, serialized Turns (the UX should be improved to do 1 Move at a time instead of a whole Turn though).
func (b *Board) MustExecuteTurn(t turn.Turn, debug bool) {
	var mustExec func(m turn.Move, times uint8)
	if debug {
		mustExec = func(m turn.Move, times uint8) {
			for i := uint8(0); i < times; i++ {
				if ok, reason := b.ExecuteMoveIfLegal(m); !ok {
					panic(fmt.Sprintf("we couldn't execute Move %v for the %d'th time, as part of supposedly-valid Turn %v, because %s", m, i, t, reason))
				}
			}
		}
	} else {
		mustExec = func(m turn.Move, times uint8) {
			for i := uint8(0); i < times; i++ {
				b.ExecuteMoveUnsafe(m)
			}
		}
	}

	var sortable sortableMotimesPairs
	for move, numTimes := range t {
		if p := move.Requestor; (p == plyr.PCC && move.Letter == constants.LETTER_BAR_CC) || (p == plyr.PC && move.Letter == constants.LETTER_BAR_C) {
			mustExec(move, numTimes)
			continue
		}
		sortable = append(sortable, motimesPair{move, numTimes})
	}
	sort.Sort(sortable)

	for _, mtp := range sortable {
		mustExec(mtp.mo, mtp.times)
	}
}

func (b *Board) ExecuteMoveUnsafe(m turn.Move) {
	if m.IsToMoveSomethingOutOfTheBar() {
		b.decrementBar(m.Requestor)
	} else {
		fromPt := b.Points[m.PointIdx()]
		fromPt.NumCheckers--
		if fromPt.NumCheckers == 0 {
			fromPt.Owner = 0
		}
	}

	nextPointIdx, nxtPtExists := m.NextPointIdx()
	if !nxtPtExists {
		b.incrementBearoffZone(m.Requestor)
		return
	}

	nxtPt := b.Points[nextPointIdx]
	if nxtPt.Owner != 0 && nxtPt.Owner != m.Requestor {
		b.incrementBar(nxtPt.Owner)
		nxtPt.NumCheckers--
	}
	nxtPt.NumCheckers++
	nxtPt.Owner = m.Requestor
}

func (b *Board) ExecuteMoveIfLegal(m turn.Move) (bool, string) {
	moveOk, moveReason := m.IsValid()
	boardOk, boardReason := b.isLegalMove(m)
	if !moveOk || !boardOk {
		return false, moveReason + boardReason
	}

	b.ExecuteMoveUnsafe(m)
	return true, ""
}

func (b *Board) SetUp() {
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{plyr.PCC, 2}, {}, {}, {}, {}, {plyr.PC, 5}, {}, {plyr.PC, 3}, {}, {}, {}, {plyr.PCC, 5},
		{plyr.PC, 5}, {}, {}, {}, {plyr.PCC, 3}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {plyr.PC, 2},
		//                                                        clockwise player in top-left.
	}
}

func (b *Board) PipCounts() (uint16, uint16) {
	var pipC, pipCC uint16

	for i, p := range b.Points {
		basePips, chex := uint16(i)+1, uint16(p.NumCheckers)
		if p.Owner == plyr.PC {
			// the clockwise player's closest checker is at points[0].
			pipC += chex * basePips
		} else if p.Owner == plyr.PCC {
			// the counter-clockwise player's furthest checker is at points[0].
			pipCC += chex * (uint16(constants.NUM_BOARD_POINTS) - basePips + 1)
		}
	}
	pipC += uint16(b.BarC) * barPips
	pipCC += uint16(b.BarCC) * barPips

	return pipC, pipCC
}

func (b *Board) doesPlayerHaveAllRemainingCheckersInHomeBoard(p plyr.Player) bool {
	totalChexInHomeOrBearedOff := b.OffC
	if p == plyr.PCC {
		totalChexInHomeOrBearedOff = b.OffCC
	}

	homeStart, homeEnd := p.HomePointIndices()
	for i := homeStart; i <= homeEnd; i++ {
		if pt := b.Points[i]; pt.Owner == p {
			totalChexInHomeOrBearedOff += pt.NumCheckers
		}
	}

	return totalChexInHomeOrBearedOff == constants.NUM_CHECKERS_PER_PLAYER
}

func (b *Board) chexOnTheBar(p plyr.Player) uint8 {
	if p == plyr.PC {
		return b.BarC
	}
	return b.BarCC
}

// Specifically determines whether the given move is OK for moving a checker off the bar and back onto the board.
// Before running this method, you must be certain that `m` specifically is for moving a checker back onto the board!
func (b *Board) isLegalMoveForBearingOn(m turn.Move) (bool, string) {
	if (m.Requestor == plyr.PCC && m.Letter != constants.LETTER_BAR_CC) ||
		(m.Requestor == plyr.PC && m.Letter != constants.LETTER_BAR_C) {
		return false, illegalEnemyBarChex
	}

	if b.chexOnTheBar(m.Requestor) < 1 {
		return false, illegalEmptyPoint
	}

	enemy := m.Requestor.Enemy()
	enemyHomeStart, enemyHomeEnd := enemy.HomePointIndices()

	toPtIdx, _ := m.NextPointIdx()
	if toPtIdx < int8(enemyHomeStart) || toPtIdx > int8(enemyHomeEnd) {
		return false, illegalBearOntoEnemyHome
	}
	if pt := b.Points[toPtIdx]; pt.Owner == enemy && pt.NumCheckers > 1 {
		return false, illegalEnemyControlsIt
	}

	return true, ""
}

func (b *Board) isLegalMoveForNonBearingOn(m turn.Move) (bool, string) {
	if b.chexOnTheBar(m.Requestor) > 0 {
		return false, illegalBarFirst
	}

	fromPt := b.Points[m.PointIdx()]
	if fromPt.Owner != m.Requestor {
		return false, illegalEnemyRegularChex
	}
	if fromPt.NumCheckers < 1 {
		return false, illegalEmptyPoint
	}

	nxtIdx, nxtPtExists := m.NextPointIdx()
	if !nxtPtExists {
		if !b.doesPlayerHaveAllRemainingCheckersInHomeBoard(m.Requestor) {
			return false, illegalCantBearoffUntilAllAreHome
		}
		if (m.Requestor == plyr.PCC && nxtIdx < 0) || (m.Requestor == plyr.PC && nxtIdx >= int8(constants.NUM_BOARD_POINTS)) {
			return false, illegalWrongFinishLine
		}
		if ((m.Requestor == plyr.PCC && nxtIdx > int8(constants.NUM_BOARD_POINTS)) || (m.Requestor == plyr.PC && nxtIdx < -1)) && b.doesPlayerHaveAnyRemainingCheckersBehindPoint(m.Requestor, m.PointIdx()) {
			// E.g., if you roll a 6, and you have chex on your 5 and 6 point, you can only bear off the ones on the 6 point (and not the ones on the 5 until all the chex on 6 are gone).
			return false, illegalBearoffOthersFirst
		}
	} else {
		if nxtPt := b.Points[nxtIdx]; nxtPt.Owner != m.Requestor && nxtPt.NumCheckers > 1 {
			return false, illegalEnemyControlsIt
		}
	}

	return true, ""
}

func (b *Board) isLegalMove(m turn.Move) (bool, string) {
	if isForBar := m.Letter == constants.LETTER_BAR_CC || m.Letter == constants.LETTER_BAR_C; isForBar {
		return b.isLegalMoveForBearingOn(m)
	}
	return b.isLegalMoveForNonBearingOn(m)
}

func (b *Board) doesPlayerHaveAnyRemainingCheckersBehindPoint(p plyr.Player, pointIdx uint8) bool {
	homeStart, homeEnd := p.HomePointIndices()

	if p == plyr.PCC {
		for i := pointIdx - 1; i >= homeStart; i-- {
			if pt := b.Points[i]; pt.Owner == p && pt.NumCheckers > 0 {
				return true
			}
		}
	} else {
		for i := pointIdx + 1; i <= homeEnd; i++ {
			if pt := b.Points[i]; pt.Owner == p && pt.NumCheckers > 0 {
				return true
			}
		}
	}
	return false
}

func (b *Board) incrementBar(p plyr.Player) {
	if p == plyr.PCC {
		b.BarCC++
	} else {
		b.BarC++
	}
}

func (b *Board) decrementBar(p plyr.Player) {
	if p == plyr.PCC {
		b.BarCC--
	} else {
		b.BarC--
	}
}

func (b *Board) incrementBearoffZone(p plyr.Player) {
	if p == plyr.PCC {
		b.OffCC++
		if b.OffCC == constants.NUM_CHECKERS_PER_PLAYER {
			b.winner, b.winKind = plyr.PCC, detectWinKind(b, plyr.PCC)
		}
	} else {
		b.OffC++
		if b.OffC == constants.NUM_CHECKERS_PER_PLAYER {
			b.winner, b.winKind = plyr.PC, detectWinKind(b, plyr.PC)
		}
	}
}
