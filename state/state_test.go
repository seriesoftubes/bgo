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
		0.0, 0.0, 0.0, // HeroBar, EnemyBar, Hero-EnemyBar.
		0.0, 0.0, 0.0, // HeroOff, EnemyOff, Hero-EnemyOff.
		// Below here: from back to front, the points from "X"'s POV.

		// Point "a"
		1.0,  // ownerStatus
		0.0,  // numBeyond2
		1.0,  // isSecured
		0.0,  // oppositeDiff
		15.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		5.0,  // distToClosestEnemySecuredPoint
		1.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		15.0, // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		5.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "b"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		15.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		4.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "c"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		15.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		3.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "d"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		15.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		2.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "e"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		15.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		1.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "f"
		-1.0,  // ownerStatus
		3.0,   // numBeyond2
		1.0,   // isSecured
		0.0,   // oppositeDiff
		10.0,  // numEnemyChexInFront
		0.0,   // distToClosestEnemyBlotPoint
		2.0,   // distToClosestEnemySecuredPoint
		-1.0,  // ownerStatus * isSecured
		-3.0,  // ownerStatus * numBeyond2
		-10.0, // ownerStatus * numEnemyChexInFront
		0.0,   // ownerStatus * distToClosestEnemyBlotPoint
		-2.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "g"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		10.0, // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		1.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "h"
		-1.0, // ownerStatus
		1.0,  // numBeyond2
		1.0,  // isSecured
		0.0,  // oppositeDiff
		7.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		5.0,  // distToClosestEnemySecuredPoint
		-1.0, // ownerStatus * isSecured
		-1.0, // ownerStatus * numBeyond2
		-7.0, // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		-5.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "i"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		7.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		4.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "j"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		7.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		3.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "k"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		7.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		2.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "l"
		1.0, // ownerStatus
		3.0, // numBeyond2
		1.0, // isSecured
		0.0, // oppositeDiff
		7.0, // numEnemyChexInFront
		0.0, // distToClosestEnemyBlotPoint
		1.0, // distToClosestEnemySecuredPoint
		1.0, // ownerStatus * isSecured
		3.0, // ownerStatus * numBeyond2
		7.0, // ownerStatus * numEnemyChexInFront
		0.0, // ownerStatus * distToClosestEnemyBlotPoint
		1.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "m"
		-1.0,  // ownerStatus
		3.0,   // numBeyond2
		1.0,   // isSecured
		0.0,   // oppositeDiff
		2.0,   // numEnemyChexInFront
		0.0,   // distToClosestEnemyBlotPoint
		11.0,  // distToClosestEnemySecuredPoint
		-1.0,  // ownerStatus * isSecured
		-3.0,  // ownerStatus * numBeyond2
		-2.0,  // ownerStatus * numEnemyChexInFront
		0.0,   // ownerStatus * distToClosestEnemyBlotPoint
		-11.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "n"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		10.0, // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "o"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		9.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "p"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		8.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "q"
		1.0, // ownerStatus
		1.0, // numBeyond2
		1.0, // isSecured
		0.0, // oppositeDiff
		2.0, // numEnemyChexInFront
		0.0, // distToClosestEnemyBlotPoint
		7.0, // distToClosestEnemySecuredPoint
		1.0, // ownerStatus * isSecured
		1.0, // ownerStatus * numBeyond2
		2.0, // ownerStatus * numEnemyChexInFront
		0.0, // ownerStatus * distToClosestEnemyBlotPoint
		7.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "r"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		6.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "s"
		1.0, // ownerStatus
		3.0, // numBeyond2
		1.0, // isSecured
		0.0, // oppositeDiff
		2.0, // numEnemyChexInFront
		0.0, // distToClosestEnemyBlotPoint
		5.0, // distToClosestEnemySecuredPoint
		1.0, // ownerStatus * isSecured
		3.0, // ownerStatus * numBeyond2
		2.0, // ownerStatus * numEnemyChexInFront
		0.0, // ownerStatus * distToClosestEnemyBlotPoint
		5.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "t"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		4.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "u"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		3.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "v"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		2.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "w"
		0.0,  // ownerStatus
		-2.0, // numBeyond2
		0.0,  // isSecured
		0.0,  // oppositeDiff
		2.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		1.0,  // distToClosestEnemySecuredPoint
		0.0,  // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "x"
		-1.0, // ownerStatus
		0.0,  // numBeyond2
		1.0,  // isSecured
		0.0,  // oppositeDiff
		0.0,  // numEnemyChexInFront
		0.0,  // distToClosestEnemyBlotPoint
		0.0,  // distToClosestEnemySecuredPoint
		-1.0, // ownerStatus * isSecured
		0.0,  // ownerStatus * numBeyond2
		0.0,  // ownerStatus * numEnemyChexInFront
		0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		0.0,  // ownerStatus * distToClosestEnemySecuredPoint
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
