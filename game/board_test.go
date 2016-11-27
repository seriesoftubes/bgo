package game

import (
	"reflect"
	"testing"

	"github.com/seriesoftubes/bgo/constants"
)

var defaultPlayer *Player = PCC

// TestLegalMoves tests getting legal moves for the initial, clean board.
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

func mLetters(moves []*Move) []string {
	var out []string
	for _, m := range moves {
		out = append(out, m.Letter)
	}
	return out
}
