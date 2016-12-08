package game

import (
	"reflect"
	"testing"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
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
	defaultPlayer := plyr.PCC // "X"
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
		b.SetUp()

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
	defaultPlayer := plyr.PCC // "X"
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
		b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{plyr.PCC, 2}, {}, {}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 1}, {plyr.PC, 1}, {plyr.PC, 2}, {}, {}, {}, {plyr.PCC, 5},
			{plyr.PC, 5}, {}, {}, {}, {plyr.PCC, 3}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {plyr.PC, 2},
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
	// plyr.PCC == "X", plyr.PC = O
	cases := []struct {
		player      *plyr.Player
		diceAmt     uint8
		wantLetters []string
	}{
		{plyr.PCC, 1, []string{"a", "l", "q", "s"}},
		{plyr.PCC, 2, []string{"l", "q", "s"}},
		{plyr.PCC, 3, []string{"l", "q", "s"}},
		{plyr.PCC, 4, []string{"l", "q", "s"}},
		{plyr.PCC, 5, []string{"l", "q", "s"}},
		{plyr.PCC, 6, []string{"a", "l", "q"}},

		{plyr.PC, 1, []string{"c", "d", "e", "f"}},
		{plyr.PC, 2, []string{"b", "d", "e", "f"}},
		{plyr.PC, 3, []string{"c", "e", "f"}},
		{plyr.PC, 4, []string{"d", "f"}},
		{plyr.PC, 5, []string{"e"}},
		{plyr.PC, 6, []string{"f"}},
	}
	for _, c := range cases {
		b := Board{}
		b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{plyr.PCC, 2}, {plyr.PC, 1}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 5}, {plyr.PC, 5}, {}, {}, {}, {}, {}, {plyr.PCC, 5},
			{}, {}, {}, {}, {plyr.PCC, 3}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {},
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
	// plyr.PCC == "X", plyr.PC = O
	cases := []struct {
		player      *plyr.Player
		diceAmt     uint8
		wantLetters []string
	}{
		{plyr.PCC, 1, []string{"t", "u", "v", "w", "x"}},
		{plyr.PCC, 2, []string{"t", "u", "v", "w"}},
		{plyr.PCC, 3, []string{"t", "u", "v"}},
		{plyr.PCC, 4, []string{"t", "u"}},
		{plyr.PCC, 5, []string{"t"}},
		{plyr.PCC, 6, []string{"t"}},

		{plyr.PC, 1, []string{"b", "c", "d", "e", "f"}},
		{plyr.PC, 2, []string{"b", "c", "d", "e", "f"}},
		{plyr.PC, 3, []string{"c", "d", "e", "f"}},
		{plyr.PC, 4, []string{"d", "e", "f"}},
		{plyr.PC, 5, []string{"e", "f"}},
		{plyr.PC, 6, []string{"f"}},
	}
	for _, c := range cases {
		b := Board{}
		b.OffCC = 5
		b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{}, {plyr.PC, 1}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 5}, {plyr.PC, 5}, {}, {}, {}, {}, {}, {},
			{}, {}, {}, {}, {}, {}, {}, {plyr.PCC, 1}, {plyr.PCC, 1}, {plyr.PCC, 2}, {plyr.PCC, 1}, {plyr.PCC, 5},
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
	b.SetUp()

	m := &turn.Move{Requestor: plyr.PCC, Letter: "a", FowardDistance: 6}

	// Original state
	fromIdx := constants.Alpha2Num[m.Letter]
	toIdx, _ := constants.Alpha2Num["g"]
	fromPt, toPt := b.Points[fromIdx], b.Points[toIdx]
	fromPtChex, toPtChex := fromPt.NumCheckers, toPt.NumCheckers

	ok, reason := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}

	if fromPt.NumCheckers != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != plyr.PCC {
		t.Errorf("Destination point should be owned by plyr.PCC")
	}
}

func TestExecuteMoveIfLegalFromBar(t *testing.T) {
	b := Board{}
	b.SetUp()
	// Simulate having 1 chex on the bar for plyr.PCC.
	b.Points[0].NumCheckers--
	b.BarCC = 1

	m := &turn.Move{Requestor: plyr.PCC, Letter: "y", FowardDistance: 1}

	// Original state
	toIdx, _ := constants.Alpha2Num["a"]
	toPt := b.Points[toIdx]
	fromPtChex, toPtChex := b.BarCC, toPt.NumCheckers

	ok, reason := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}

	if b.BarCC != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != plyr.PCC {
		t.Errorf("Destination point owner should be plyr.PCC")
	}
}

func TestExecuteMoveIfLegalFromBarForPlayerC(t *testing.T) {
	b := Board{}
	b.SetUp()
	// Simulate having 1 chex on the bar for plyr.PCC.
	b.Points[constants.Alpha2Num["x"]].NumCheckers--
	b.BarC = 1

	m := &turn.Move{Requestor: plyr.PC, Letter: "z", FowardDistance: 2}

	// Original state
	toIdx, _ := constants.Alpha2Num["w"]
	toPt := b.Points[toIdx]
	fromPtChex, toPtChex := b.BarC, toPt.NumCheckers

	ok, reason := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}

	if b.BarC != fromPtChex-1 {
		t.Errorf("Did not move any checkers away from the original point.")
	} else if toPt.NumCheckers != toPtChex+1 {
		t.Errorf("Did not move any checkers to the destination.")
	} else if toPt.Owner != plyr.PC {
		t.Errorf("Destination point owner should be plyr.PC")
	}
}

func TestExecuteMoveIfLegalBearOff(t *testing.T) {
	b := Board{}
	b.SetUp()
	b.OffCC = 5
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{}, {plyr.PC, 1}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 5}, {plyr.PC, 5}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {plyr.PCC, 1}, {plyr.PCC, 1}, {plyr.PCC, 2}, {plyr.PCC, 1}, {plyr.PCC, 5},
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

	m := &turn.Move{Requestor: plyr.PCC, Letter: "t", FowardDistance: 6}

	// Original state
	fromIdx, _ := constants.Alpha2Num[m.Letter]
	fromPt := b.Points[fromIdx]
	fromPtChex, toPtChex := fromPt.NumCheckers, b.OffCC

	ok, reason := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
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
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{plyr.PCC, 1}, {}, {}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 1}, {plyr.PC, 1}, {plyr.PC, 2}, {}, {}, {}, {plyr.PCC, 5},
		{plyr.PC, 5}, {}, {}, {}, {plyr.PCC, 4}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {plyr.PC, 2},
		//                                                        clockwise player in top-left.
	}

	m := &turn.Move{Requestor: plyr.PCC, Letter: "a", FowardDistance: 5}
	// Expect the state to be:
	// 0 on "a" (and nil Owner),
	// 1 on "f" (and plyr.PCC Owner)
	// 1 in b.BarC

	// Original state
	originalBarC := b.BarC
	if originalBarC != 0 {
		t.Errorf("thought there were 0 on bar C to begin with, got %v", originalBarC)
	}

	fromIdx, _ := constants.Alpha2Num[m.Letter]
	toIdx := constants.Alpha2Num["f"]
	fromPt, toPt := b.Points[fromIdx], b.Points[toIdx]
	fromPtOwner, toPtOwner := fromPt.Owner, toPt.Owner
	fromPtChex, toPtChex := fromPt.NumCheckers, toPt.NumCheckers

	if fromPtChex != 1 {
		t.Errorf("thought there was only 1 checker on the from point, got %v", fromPtChex)
	} else if toPtChex != 1 {
		t.Errorf("thought there was only 1 checker on the dets point, got %v", toPtChex)
	}
	if fromPtOwner != plyr.PCC {
		t.Errorf("thought the owner of the from point was gonna be plyr.PCC")
	} else if toPtOwner != plyr.PC {
		t.Errorf("thought the owner of the destination was gonna be plyr.PC")
	}

	ok, reason := b.ExecuteMoveIfLegal(m)
	if !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
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
	if toPt.Owner != plyr.PCC {
		t.Errorf("expected plyr.PCC to be the new owner of the destination point, got %v", toPt.Owner)
	}
	if b.BarC != originalBarC+1 {
		t.Errorf("expected there to be 1 new checker on BarC. got %v", b.BarC)
	}
}

// Tests that a proper winner is set on a SingleGame game-winning move and not before then.
func TestExecuteMoveIfLegalWinSingleGame(t *testing.T) {
	b := Board{}
	b.OffCC = 14
	b.OffC = 14
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		{}, {plyr.PC, 1}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {plyr.PCC, 1},
	}

	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 X  -  -  -  -  - |m| -  -  -  -  -  -
	                  |m|
	                  |m|
	                  |m|
	                  |m|
	                  |m|


	                  |w|
	                  |w|
	                  |w|
	                  |w|
	                  |w|
	 -  O  -  -  -  - |w| -  -  -  -  -  -
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	boardPC := b.Copy()
	moveForPC := &turn.Move{Requestor: plyr.PC, Letter: "b", FowardDistance: 6}

	if boardPC.winner != nil {
		t.Errorf("expected no winner to be set but got %v", *boardPC.winner)
	} else if boardPC.winKind != WinKindNotWon {
		t.Errorf("expected win state to be %v but got %v", WinKindNotWon, boardPC.winKind)
	}

	if ok, reason := boardPC.ExecuteMoveIfLegal(moveForPC); !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}
	if boardPC.winner != plyr.PC {
		t.Errorf("expected winner to be %v but got %v", *plyr.PC, *boardPC.winner)
	} else if boardPC.winKind != WinKindSingleGame {
		t.Errorf("expected win state to be %v but got %v", WinKindSingleGame, boardPC.winKind)
	} else if boardPC.OffC != constants.NUM_CHECKERS_PER_PLAYER {
		t.Errorf("expected %d chex to be beared off, but got %v", constants.NUM_CHECKERS_PER_PLAYER, boardPC.OffC)
	}

	boardPCC := b.Copy()
	moveForPCC := &turn.Move{Requestor: plyr.PCC, Letter: "x", FowardDistance: 1}

	if boardPCC.winner != nil {
		t.Errorf("expected no winner to be set but got %v", *boardPCC.winner)
	} else if boardPCC.winKind != WinKindNotWon {
		t.Errorf("expected win state to be %v but got %v", WinKindNotWon, boardPCC.winKind)
	}

	if ok, reason := boardPCC.ExecuteMoveIfLegal(moveForPCC); !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}
	if boardPCC.winner != plyr.PCC {
		t.Errorf("expected winner to be %v but got %v", *plyr.PCC, *boardPCC.winner)
	} else if boardPCC.winKind != WinKindSingleGame {
		t.Errorf("expected win state to be %v but got %v", WinKindSingleGame, boardPCC.winKind)
	} else if boardPCC.OffCC != constants.NUM_CHECKERS_PER_PLAYER {
		t.Errorf("expected %d chex to be beared off, but got %v", constants.NUM_CHECKERS_PER_PLAYER, boardPCC.OffCC)
	}
}

// Tests that a proper winner is set on a Gammon game-winning move and not before then.
func TestExecuteMoveIfLegalWinGammon(t *testing.T) {
	b := Board{}
	b.OffC = 14
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		{}, {plyr.PC, 1}, {}, {}, {}, {}, {}, {plyr.PCC, 15}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	/* Board looks like:
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 -  -  -  -  -  - |m| -  -  -  -  -  -
	                  |m|
	                  |m|
	                  |m|
	                  |m|
	                  |m|


	                  |w|   15
	                  |w|    X
	                  |w|    X
	                  |w|    X
	                  |w|    X
	 -  O  -  -  -  - |w| -  X  -  -  -  -
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	boardPC := b.Copy()
	moveForPC := &turn.Move{Requestor: plyr.PC, Letter: "b", FowardDistance: 2}

	if boardPC.winner != nil {
		t.Errorf("expected no winner to be set but got %v", *boardPC.winner)
	} else if boardPC.winKind != WinKindNotWon {
		t.Errorf("expected win state to be %v but got %v", WinKindNotWon, boardPC.winKind)
	}

	if ok, reason := boardPC.ExecuteMoveIfLegal(moveForPC); !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}
	if boardPC.winner != plyr.PC {
		t.Errorf("expected winner to be %v but got %v", *plyr.PC, *boardPC.winner)
	} else if boardPC.winKind != WinKindGammon {
		t.Errorf("expected win state to be %v but got %v", WinKindGammon, boardPC.winKind)
	} else if boardPC.OffC != constants.NUM_CHECKERS_PER_PLAYER {
		t.Errorf("expected %d chex to be beared off, but got %v", constants.NUM_CHECKERS_PER_PLAYER, boardPC.OffC)
	}
}

// Tests that a proper winner is set on a Backgammon game-winning move and not before then.
func TestExecuteMoveIfLegalWinBackgammon(t *testing.T) {
	b := Board{}
	b.OffC = 14
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		{}, {plyr.PC, 1}, {}, {}, {}, {plyr.PCC, 1}, {}, {plyr.PCC, 14}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	/* Board looks like: (notice how X has 1 left in enemy's home)
	 x  w  v  u  t  s     r  q  p  o  n  m
	=======================================
	 -  -  -  -  -  - |m| -  -  -  -  -  -
	                  |m|
	                  |m|
	                  |m|
	                  |m|
	                  |m|


	                  |w|   14
	                  |w|    X
	                  |w|    X
	                  |w|    X
	                  |w|    X
	 -  O  -  -  -  X |w| -  X  -  -  -  -
	=======================================
	 a  b  c  d  e  f     g  h  i  j  k  l
	*/
	boardPC := b.Copy()
	moveForPC := &turn.Move{Requestor: plyr.PC, Letter: "b", FowardDistance: 2}

	if boardPC.winner != nil {
		t.Errorf("expected no winner to be set but got %v", *boardPC.winner)
	} else if boardPC.winKind != WinKindNotWon {
		t.Errorf("expected win state to be %v but got %v", WinKindNotWon, boardPC.winKind)
	}

	if ok, reason := boardPC.ExecuteMoveIfLegal(moveForPC); !ok {
		t.Errorf("Test move was not legal. Change the test! %v", reason)
	}
	if boardPC.winner != plyr.PC {
		t.Errorf("expected winner to be %v but got %v", *plyr.PC, *boardPC.winner)
	} else if boardPC.winKind != WinKindBackgammon {
		t.Errorf("expected win state to be %v but got %v", WinKindBackgammon, boardPC.winKind)
	} else if boardPC.OffC != constants.NUM_CHECKERS_PER_PLAYER {
		t.Errorf("expected %d chex to be beared off, but got %v", constants.NUM_CHECKERS_PER_PLAYER, boardPC.OffC)
	}
}

func strSlicesEqual(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return len(a) == len(b)
	}
	return reflect.DeepEqual(a, b)
}

func mLetters(moves []*turn.Move) []string {
	var out []string
	for _, m := range moves {
		out = append(out, m.Letter)
	}
	return out
}
