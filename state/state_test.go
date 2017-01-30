package state

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

func TestStartingBoard(t *testing.T) {
	b := &game.Board{}
	b.SetUp()
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

	got := DetectState(plyr.PCC, b) // The "X" player.
	want := State{
		0.0, // isRace

		// Starting with player PCC (the hero)

		0.0,       // has1 on bar
		0.0,       // has >1 on bar + amt in excess
		0.0,       // blot % in enemy home
		5.0 / 6.0, // landable % in enemy home
		0.0,       // has1 beared off
		0.0,       // has >1 beared off + amt in excess
		0.0,       // diff % beared off
		// From here on down, for the current iteration's player (first iteration is hero, 2nd is enemy), the player-owned chex count for the point.
		// Point "a"  2
		1.0, // has 1 checker
		1.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "b"  8
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "c"  14
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "d"  20
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "e"  26
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "f"  32
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "g"  38
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "h"  44
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "i"  50
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "j"  56
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "k"  62
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "l"  68
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		1.0, // has 4 chex
		1.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "m"  74
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "n"  80
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "o"  86
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "p"  92
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "q"  98
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "r"  104
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "s"  110
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		1.0, // has 4 chex
		1.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "t"  116
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "u"  122
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "v"  128
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "w"  134
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "x"  140
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6

		// Now onto the enemy player, PC aka "O".  146
		0.0,       // has1 on bar
		0.0,       // has >1 on bar + amt in excess
		0.0,       // blot % in enemy home
		5.0 / 6.0, // landable % in enemy home
		0.0,       // has1 beared off
		0.0,       // has >1 beared off + amt in excess
		0.0,       // diff % beared off
		// From here on down, for the current iteration's player (first iteration is hero, 2nd is enemy), the player-owned chex count for the point.
		// Point "a"  148
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "b"  154
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "c"  160
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "d"  166
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "e"  172
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "f"  178
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		1.0, // has 4 chex
		1.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "g"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "h"
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "i"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "j"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "k"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "l"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "m"
		1.0, // has 1 checker
		1.0, // has 2 checkers
		1.0, // has 3 chex
		1.0, // has 4 chex
		1.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "n"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "o"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "p"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "q"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "r"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "s"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "t"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "u"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "v"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "w"
		0.0, // has 1 checker
		0.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
		// Point "x"
		1.0, // has 1 checker
		1.0, // has 2 checkers
		0.0, // has 3 chex
		0.0, // has 4 chex
		0.0, // has 5 chex
		0.0, // has 6 chex + num beyond 6
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v \nexpected %v", got, want)
		for gi, gv := range got {
			if wv := want[gi]; wv != gv {
				fmt.Println(gi, "Got:", gv, "Want:", wv)
			}
		}
	}
}
