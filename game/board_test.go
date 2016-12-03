package game

import (
	"reflect"
	"testing"

	"github.com/seriesoftubes/bgo/constants"
)

// TestLegalMovesPlainBoard tests getting legal moves for the initial, clean board.
func TestLegalMovesPlainBoard(t *testing.T) {
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 O  -  -  -  -  X |m| -  X  -  -  -  O
	 O              X |m|    X           O
	                X |m|    X           O
	                X |m|                O
	                X |m|                O
	                  |m|


	                  |w|
	                O |w|                X
	                O |w|                X
	                O |w|    O           X
	 X              O |w|    O           X
	 X  -  -  -  -  O |w| -  O  -  -  -  X
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	defaultPlayer := PCC // "X"
	cases := []struct {
		diceAmt     uint8
		wantLetters []string
	}{
		{1, []string{"a", "q", "s"}},
		{2, []string{"a", "l", "q", "s"}},
		{3, []string{"a", "l", "q", "s"}},
		{4, []string{"a", "l", "q", "s"}},
		{5, []string{"l", "q"}},
		{6, []string{"a", "l", "q"}},
	}
	for _, c := range cases {
		b := Board{}
		b.setUp()

		gotLetters := mLetters(b.LegalMoves(defaultPlayer, c.diceAmt))
		if !reflect.DeepEqual(gotLetters, c.wantLetters) {
			t.Errorf("LegalMoves for diceAmt: %v unexpected. got %v want %v", c.diceAmt, gotLetters, c.wantLetters)
		}

		// Proves that, whenever there is at least 1 on the bar, the player can only move that bar checker.
		b.incrementBar(defaultPlayer)
		incGotLetters := mLetters(b.LegalMoves(defaultPlayer, c.diceAmt))
		var incWantLetters []string
		if c.diceAmt != 6 {
			incWantLetters = append(incWantLetters, constants.LETTER_BAR_CC)
		}
		if !strSlicesEqual(incGotLetters, incWantLetters) {
			t.Errorf("LegalMoves (with bar) for diceAmt: %v unexpected. got %v wanted %v", c.diceAmt, incGotLetters, incWantLetters)
		}
	}
}

// TestLegalMovesTakeOverEnemy tests getting legal moves when one move can involve stepping on an enemy checker.
func TestLegalMovesTakeOverEnemy(t *testing.T) {
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 O  -  -  -  -  X |m| -  X  -  -  -  O
	 O              X |m|    X           O
	                X |m|    X           O
	                X |m|                O
	                X |m|                O
	                  |m|


	                  |w|
	                  |w|                X
	                  |w|                X
	                  |w|                X
	 X        O  O    |w|    O           X
	 X  -  -  O  O  O |w| O  O  -  -  -  X
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	defaultPlayer := PCC // "X"
	cases := []struct {
		diceAmt     uint8
		wantLetters []string
	}{
		{1, []string{"a", "q", "s"}},
		{2, []string{"a", "l", "q", "s"}},
		{3, []string{"l", "q", "s"}},
		{4, []string{"l", "q", "s"}},
		{5, []string{"a", "l", "q"}}, // We can move "a" to "f"
		{6, []string{"a", "l", "q"}}, // We can move "a" to "g"
	}
	for _, c := range cases {
		b := Board{}
		b.Points = [NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{PCC, 2}, {}, {}, {PC, 2}, {PC, 2}, {PC, 1}, {PC, 1}, {PC, 2}, {}, {}, {}, {PCC, 5},
			{PC, 5}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {PC, 2},
			//                                                        clockwise player in top-left.
		}

		gotLetters := mLetters(b.LegalMoves(defaultPlayer, c.diceAmt))
		if !strSlicesEqual(gotLetters, c.wantLetters) {
			t.Errorf("LegalMoves for diceAmt: %v unexpected. got %v want %v", c.diceAmt, gotLetters, c.wantLetters)
		}
	}
}

// TestLegalMovesBearOff tests getting legal moves when one move can involve bearing off.
func TestLegalMovesBearOff(t *testing.T) {
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 -  -  -  -  -  X |m| -  X  -  -  -  -
	                X |m|    X
	                X |m|    X
	                X |m|
	                X |m|
	                  |m|


	                  |w|
	             O  O |w|                X
	             O  O |w|                X
	             O  O |w|                X
	 X     O  O  O  O |w|                X
	 X  O  O  O  O  O |w| -  -  -  -  -  X
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	// PCC == "X", PC = O
	cases := []struct {
		player      *Player
		diceAmt     uint8
		wantLetters []string
	}{
		{PCC, 1, []string{"a", "l", "q", "s"}},
		{PCC, 2, []string{"l", "q", "s"}},
		{PCC, 3, []string{"l", "q", "s"}},
		{PCC, 4, []string{"l", "q", "s"}},
		{PCC, 5, []string{"l", "q", "s"}},
		{PCC, 6, []string{"a", "l", "q"}},

		{PC, 1, []string{"c", "d", "e", "f"}},
		{PC, 2, []string{"b", "d", "e", "f"}},
		{PC, 3, []string{"c", "e", "f"}},
		{PC, 4, []string{"d", "f"}},
		{PC, 5, []string{"e"}},
		{PC, 6, []string{"f"}},
	}
	for _, c := range cases {
		b := Board{}
		b.Points = [NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{PCC, 2}, {PC, 1}, {PC, 2}, {PC, 2}, {PC, 5}, {PC, 5}, {}, {}, {}, {}, {}, {PCC, 5},
			{}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {},
			//                                                        clockwise player in top-left.
		}

		gotLetters := mLetters(b.LegalMoves(c.player, c.diceAmt))
		if !strSlicesEqual(gotLetters, c.wantLetters) {
			t.Errorf("LegalMoves for Player %q and diceAmt %v unexpected. got %v want %v", *c.player, c.diceAmt, gotLetters, c.wantLetters)
		}
	}
}

// TestLegalMovesBearOffBoth tests getting legal moves when one move can involve bearing off for either player.
func TestLegalMovesBearOffBoth(t *testing.T) {
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 X  X  X  X  X    |m| -  -  -  -  -  -
	 X     X          |m|
	 X                |m|
	 X                |m|
	 X                |m|
	                  |m|


	                  |w|
	             O  O |w|
	             O  O |w|
	             O  O |w|
	       O  O  O  O |w|
	 -  O  O  O  O  O |w| -  -  -  -  -  -
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	// PCC == "X", PC = O
	cases := []struct {
		player      *Player
		diceAmt     uint8
		wantLetters []string
	}{
		{PCC, 1, []string{"t", "u", "v", "w", "x"}},
		{PCC, 2, []string{"t", "u", "v", "w"}},
		{PCC, 3, []string{"t", "u", "v"}},
		{PCC, 4, []string{"t", "u"}},
		{PCC, 5, []string{"t"}},
		{PCC, 6, []string{"t"}},

		{PC, 1, []string{"b", "c", "d", "e", "f"}},
		{PC, 2, []string{"b", "c", "d", "e", "f"}},
		{PC, 3, []string{"c", "d", "e", "f"}},
		{PC, 4, []string{"d", "e", "f"}},
		{PC, 5, []string{"e", "f"}},
		{PC, 6, []string{"f"}},
	}
	for _, c := range cases {
		b := Board{}
		b.OffCC = 5
		b.Points = [NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{}, {PC, 1}, {PC, 2}, {PC, 2}, {PC, 5}, {PC, 5}, {}, {}, {}, {}, {}, {},
			{}, {}, {}, {}, {}, {}, {}, {PCC, 1}, {PCC, 1}, {PCC, 2}, {PCC, 1}, {PCC, 5},
			//                                                        clockwise player in top-left.
		}

		gotLetters := mLetters(b.LegalMoves(c.player, c.diceAmt))
		if !strSlicesEqual(gotLetters, c.wantLetters) {
			t.Errorf("LegalMoves for Player %q and diceAmt %v unexpected. got %v want %v", *c.player, c.diceAmt, gotLetters, c.wantLetters)
		}
	}
}

func TestExecuteMoveIfLegal(t *testing.T) {
	b := Board{}
	b.setUp()

	m := &Move{Requestor: PCC, Letter: "a", FowardDistance: 6}

	// Original state
	fromIdx := alpha2Num[m.Letter]
	toIdx, _ := alpha2Num["g"]
	fromPt, toPt := b.Points[fromIdx], b.Points[toIdx]
	fromPtChex, toPtChex := fromPt.NumCheckers, toPt.NumCheckers

	ok := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test!")
	}

	if fromPt.NumCheckers != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != PCC {
		t.Errorf("Destination point should be owned by PCC")
	}
}

func TestExecuteMoveIfLegalFromBar(t *testing.T) {
	b := Board{}
	b.setUp()
	// Simulate having 1 chex on the bar for PCC.
	b.Points[0].NumCheckers--
	b.BarCC = 1

	m := &Move{Requestor: PCC, Letter: "y", FowardDistance: 1}

	// Original state
	toIdx, _ := alpha2Num["a"]
	toPt := b.Points[toIdx]
	fromPtChex, toPtChex := b.BarCC, toPt.NumCheckers

	ok := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test!")
	}

	if b.BarCC != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != PCC {
		t.Errorf("Destination point owner should be PCC")
	}
}

func TestExecuteMoveIfLegalFromBarForPlayerC(t *testing.T) {
	b := Board{}
	b.setUp()
	// Simulate having 1 chex on the bar for PCC.
	b.Points[alpha2Num["x"]].NumCheckers--
	b.BarC = 1

	m := &Move{Requestor: PC, Letter: "z", FowardDistance: 2}

	// Original state
	toIdx, _ := alpha2Num["w"]
	toPt := b.Points[toIdx]
	fromPtChex, toPtChex := b.BarC, toPt.NumCheckers

	ok := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test!")
	}

	if b.BarC != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != PC {
		t.Errorf("Destination point owner should be PC")
	}
}

func TestExecuteMoveIfLegalBearOff(t *testing.T) {
	b := Board{}
	b.setUp()
	b.OffCC = 5
	b.Points = [NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{}, {PC, 1}, {PC, 2}, {PC, 2}, {PC, 5}, {PC, 5}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {PCC, 1}, {PCC, 1}, {PCC, 2}, {PCC, 1}, {PCC, 5},
		//                                                        clockwise player in top-left.
	}
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 X  X  X  X  X    |m| -  -  -  -  -  -
	 X     X          |m|
	 X                |m|
	 X                |m|
	 X                |m|
	                  |m|


	                  |w|
	             O  O |w|
	             O  O |w|
	             O  O |w|
	       O  O  O  O |w|
	 -  O  O  O  O  O |w| -  -  -  -  -  -
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/

	m := &Move{Requestor: PCC, Letter: "t", FowardDistance: 6}

	// Original state
	fromIdx, _ := alpha2Num[m.Letter]
	fromPt := b.Points[fromIdx]
	fromPtChex, toPtChex := fromPt.NumCheckers, b.OffCC

	ok := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test!")
	}

	if fromPt.NumCheckers != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if b.OffCC != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	}
}

// TestExecuteMoveIfLegalTakeoverEnemy tests executing a move that captures an enemy checker.
func TestExecuteMoveIfLegalTakeoverEnemy(t *testing.T) {
	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 O  -  -  -  -  X |m| -  X  -  -  -  O
	 O              X |m|    X           O
	                X |m|    X           O
	                X |m|    X           O
	                X |m|                O
	                  |m|


	                  |w|
	                  |w|                X
	                  |w|                X
	                  |w|                X
	          O  O    |w|    O           X
	 X  -  -  O  O  O |w| O  O  -  -  -  X
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	b := Board{}
	b.Points = [NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{PCC, 1}, {}, {}, {PC, 2}, {PC, 2}, {PC, 1}, {PC, 1}, {PC, 2}, {}, {}, {}, {PCC, 5},
		{PC, 5}, {}, {}, {}, {PCC, 4}, {}, {PCC, 5}, {}, {}, {}, {}, {PC, 2},
		//                                                        clockwise player in top-left.
	}

	m := &Move{Requestor: PCC, Letter: "a", FowardDistance: 5}
	// Expect the state to be:
	// 0 on "a" (and nil Owner),
	// 1 on "f" (and PCC Owner)
	// 1 in b.BarC

	// Original state
	originalBarC := b.BarC
	if originalBarC != 0 {
		t.Errorf("thought there were 0 on bar C to begin with, got %v", originalBarC)
	}

	fromIdx, _ := alpha2Num[m.Letter]
	toIdx := alpha2Num["f"]
	fromPt, toPt := b.Points[fromIdx], b.Points[toIdx]
	fromPtOwner, toPtOwner := fromPt.Owner, toPt.Owner
	fromPtChex, toPtChex := fromPt.NumCheckers, toPt.NumCheckers

	if fromPtChex != 1 {
		t.Errorf("thought there was only 1 checker on the from point, got %v", fromPtChex)
	} else if toPtChex != 1 {
		t.Errorf("thought there was only 1 checker on the dets point, got %v", toPtChex)
	}
	if fromPtOwner != PCC {
		t.Errorf("thought the owner of the from point was gonna be PCC")
	} else if toPtOwner != PC {
		t.Errorf("thought the owner of the destination was gonna be PC")
	}

	ok := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test!")
	}

	if fromPt.Owner != nil {
		t.Errorf("expected there to be no owner of the from point, got %v", *fromPt.Owner)
	}
	if fromPt.NumCheckers != 0 {
		t.Errorf("expected there to be 0 checkers on the from point, got %v", fromPt.NumCheckers)
	}
	if toPt.NumCheckers != 1 {
		t.Errorf("expected there to be 1 checker on the destination point, got %v", toPt.NumCheckers)
	}
	if toPt.Owner != PCC {
		t.Errorf("expected PCC to be the new owner of the destination point, got %v", toPt.Owner)
	}
	if b.BarC != originalBarC+1 {
		t.Errorf("expected there to be 1 new checker on BarC. got %v", b.BarC)
	}
}

func strSlicesEqual(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return len(a) == len(b)
	}
	return reflect.DeepEqual(a, b)
}

func mLetters(moves []*Move) []string {
	var out []string
	for _, m := range moves {
		out = append(out, m.Letter)
	}
	return out
}
