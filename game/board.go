package game

const (
	numCheckersPerPlayer uint8 = 15
	NUM_BOARD_POINTS     uint8 = 24
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
	/*
		  some rules:
		    1) if you have >=1 piece on the bar, you must move them out first before anything else
		    2) unless there will be all their remaining checkers in the endzone,
		       the ending point's index must be between 0 and 23. endzone game can go beyond
		    3) the ending point must either be owned by me, or have <= 1 checker on it
		        TODO: if there is 1 checker of the enemy's on a space that you move onto, put that piece on the bar before moving yours there.
	      4) the player must have no legal moves left, or have used 2 moves, by the end of their turn.
	          TODO: in input parser, parse the whole turn and reject if it doesn't do this

		    maybe just create a set of all possible moves and check for membership?
		    (and if the set is empty, go to next turn):
		      TO MAKE THE SET:
		        for each piece on the board:
		          for each unique amount on the dice:
		            see if the piece can legally do it.
		              DEFINED BY: the above 5 rules

		    if the player enters an illegal move, display a warning and try again
	*/
	// TODO: consts for y and z, in a shared constants package
	isForBar := m.Letter == "y" || m.Letter == "z"
	numOnTheBar := b.chexOnTheBar(m.Requestor)
	if numOnTheBar > 0 && !isForBar {
		return false // If you have anything on the bar, you must move those things first
	}

	numChexOnCurrentPoint := numOnTheBar
	if !isForBar {
		numChexOnCurrentPoint = b.Points[m.pointIdx()].NumCheckers
	}
	if numChexOnCurrentPoint == 0 {
		return false // Cannot move a checker from an empty point
	}

	nxtIdx := m.nextPointIdx()
	nxtPtExists := nxtIdx >= 0 && nxtIdx < NUM_BOARD_POINTS
	if !nxtPtExists && !b.doesPlayerHaveAllRemainingCheckersInHomeBoard(m.Requestor) {
		return false // Can't move past the finish line unless all your remaining checkers are in your home board
	}

	if nxtPtExists {
		if nxtPt := b.Points[nxtIdx]; nxtPt.Owner != m.Requestor && nxtPt.NumCheckers > 1 {
			return false // Can't move to a point that's controlled (has >1 chex) by the enemy.
		}
	}

	return true
}

func (b *Board) LegalMoves(p *Player, r Roll) []*Move {
	var out []*Move

	/*
		for each boardpoint (and for the player's bar):
		  for each unique_number in the roll:
		    form a Move{Requestor: p, Letter: b.Letter(boardPoint), FowardDistance: unique_number}
		    if b.isLegalMove(move): add
	*/

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
		b.Points[m.pointIdx()].NumCheckers--
	}

	nxtPt := b.Points[m.nextPointIdx()]
	if nxtPt.Owner != nil && nxtPt.Owner != m.Requestor {
		if nxtPt.NumCheckers != 1 {
			panic("making a legal move to overthrow >1 enemy checker")
		}

		nxtPt.NumCheckers--
		b.incrementBar(nxtPt.Owner)
	}
	nxtPt.NumCheckers++

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

func (b *Board) PipCounts() (int, int) {
	var pipC, pipCC int

	for i, p := range b.Points {
		basePips, chex := i+1, int(p.NumCheckers)
		if p.Owner == PC {
			// the clockwise player's closest checker is at points[0].
			pipC += chex * basePips
		} else if p.Owner == PCC {
			// the counter-clockwise player's furthest checker is at points[0].
			pipCC += chex * (int(NUM_BOARD_POINTS) - basePips + 1)
		}
	}

	return pipC, pipCC
}
