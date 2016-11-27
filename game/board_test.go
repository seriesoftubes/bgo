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
		if !reflect.DeepEqual(incGotLetters, incWantLetters) {
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
		if !reflect.DeepEqual(gotLetters, c.wantLetters) {
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

		// TODO: this behavior is wrong. it should only allow bearing off when everything behind it has been beared off already.
		{PC, 1, []string{"c", "d", "e", "f"}},
		{PC, 2, []string{"b", "d", "e", "f"}},
		{PC, 3, []string{"b", "c", "e", "f"}},
		{PC, 4, []string{"b", "c", "d", "f"}},
		{PC, 5, []string{"b", "c", "d", "e"}},
		{PC, 6, []string{"b", "c", "d", "e", "f"}},
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
		if !reflect.DeepEqual(gotLetters, c.wantLetters) {
			t.Errorf("LegalMoves for diceAmt: %v unexpected. got %v want %v", c.diceAmt, gotLetters, c.wantLetters)
		}
	}
}

func mLetters(moves []*Move) []string {
	var out []string
	for _, m := range moves {
		out = append(out, m.Letter)
	}
	return out
}
