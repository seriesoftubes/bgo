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
		0: 0.0, 1: 0.0, 2: 0.0, // HeroBar, EnemyBar, Hero-EnemyBar.
		// Below here: from back to front, the points from "X"'s POV.

		// Point "a"
		3:  1.0,  // ownerStatus
		4:  0.0,  // numBeyond2
		5:  1.0,  // isSecured
		6:  0.0,  // oppositeDiff
		7:  15.0, // numEnemyChexInFront
		8:  0.0,  // distToClosestEnemyBlotPoint
		9:  5.0,  // distToClosestEnemySecuredPoint
		10: 1.0,  // ownerStatus * isSecured
		11: 0.0,  // ownerStatus * numBeyond2
		12: 15.0, // ownerStatus * numEnemyChexInFront
		13: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		14: 5.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "b"
		15: 0.0,  // ownerStatus
		16: -2.0, // numBeyond2
		17: 0.0,  // isSecured
		18: 0.0,  // oppositeDiff
		19: 15.0, // numEnemyChexInFront
		20: 0.0,  // distToClosestEnemyBlotPoint
		21: 4.0,  // distToClosestEnemySecuredPoint
		22: 0.0,  // ownerStatus * isSecured
		23: 0.0,  // ownerStatus * numBeyond2
		24: 0.0,  // ownerStatus * numEnemyChexInFront
		25: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		26: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "c"
		27: 0.0,  // ownerStatus
		28: -2.0, // numBeyond2
		29: 0.0,  // isSecured
		30: 0.0,  // oppositeDiff
		31: 15.0, // numEnemyChexInFront
		32: 0.0,  // distToClosestEnemyBlotPoint
		33: 3.0,  // distToClosestEnemySecuredPoint
		34: 0.0,  // ownerStatus * isSecured
		35: 0.0,  // ownerStatus * numBeyond2
		36: 0.0,  // ownerStatus * numEnemyChexInFront
		37: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		38: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "d"
		39: 0.0,  // ownerStatus
		40: -2.0, // numBeyond2
		41: 0.0,  // isSecured
		42: 0.0,  // oppositeDiff
		43: 15.0, // numEnemyChexInFront
		44: 0.0,  // distToClosestEnemyBlotPoint
		45: 2.0,  // distToClosestEnemySecuredPoint
		46: 0.0,  // ownerStatus * isSecured
		47: 0.0,  // ownerStatus * numBeyond2
		48: 0.0,  // ownerStatus * numEnemyChexInFront
		49: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		50: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "e"
		51: 0.0,  // ownerStatus
		52: -2.0, // numBeyond2
		53: 0.0,  // isSecured
		54: 0.0,  // oppositeDiff
		55: 15.0, // numEnemyChexInFront
		56: 0.0,  // distToClosestEnemyBlotPoint
		57: 1.0,  // distToClosestEnemySecuredPoint
		58: 0.0,  // ownerStatus * isSecured
		59: 0.0,  // ownerStatus * numBeyond2
		60: 0.0,  // ownerStatus * numEnemyChexInFront
		61: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		62: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "f"
		63: -1.0,  // ownerStatus
		64: 3.0,   // numBeyond2
		65: 1.0,   // isSecured
		66: 0.0,   // oppositeDiff
		67: 10.0,  // numEnemyChexInFront
		68: 0.0,   // distToClosestEnemyBlotPoint
		69: 2.0,   // distToClosestEnemySecuredPoint
		70: -1.0,  // ownerStatus * isSecured
		71: -3.0,  // ownerStatus * numBeyond2
		72: -10.0, // ownerStatus * numEnemyChexInFront
		73: 0.0,   // ownerStatus * distToClosestEnemyBlotPoint
		74: -2.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "g"
		75: 0.0,  // ownerStatus
		76: -2.0, // numBeyond2
		77: 0.0,  // isSecured
		78: 0.0,  // oppositeDiff
		79: 10.0, // numEnemyChexInFront
		80: 0.0,  // distToClosestEnemyBlotPoint
		81: 1.0,  // distToClosestEnemySecuredPoint
		82: 0.0,  // ownerStatus * isSecured
		83: 0.0,  // ownerStatus * numBeyond2
		84: 0.0,  // ownerStatus * numEnemyChexInFront
		85: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		86: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "h"
		87: -1.0, // ownerStatus
		88: 1.0,  // numBeyond2
		89: 1.0,  // isSecured
		90: 0.0,  // oppositeDiff
		91: 7.0,  // numEnemyChexInFront
		92: 0.0,  // distToClosestEnemyBlotPoint
		93: 5.0,  // distToClosestEnemySecuredPoint
		94: -1.0, // ownerStatus * isSecured
		95: -1.0, // ownerStatus * numBeyond2
		96: -7.0, // ownerStatus * numEnemyChexInFront
		97: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		98: -5.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "i"
		99:  0.0,  // ownerStatus
		100: -2.0, // numBeyond2
		101: 0.0,  // isSecured
		102: 0.0,  // oppositeDiff
		103: 7.0,  // numEnemyChexInFront
		104: 0.0,  // distToClosestEnemyBlotPoint
		105: 4.0,  // distToClosestEnemySecuredPoint
		106: 0.0,  // ownerStatus * isSecured
		107: 0.0,  // ownerStatus * numBeyond2
		108: 0.0,  // ownerStatus * numEnemyChexInFront
		109: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		110: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "j"
		111: 0.0,  // ownerStatus
		112: -2.0, // numBeyond2
		113: 0.0,  // isSecured
		114: 0.0,  // oppositeDiff
		115: 7.0,  // numEnemyChexInFront
		116: 0.0,  // distToClosestEnemyBlotPoint
		117: 3.0,  // distToClosestEnemySecuredPoint
		118: 0.0,  // ownerStatus * isSecured
		119: 0.0,  // ownerStatus * numBeyond2
		120: 0.0,  // ownerStatus * numEnemyChexInFront
		121: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		122: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "k"
		123: 0.0,  // ownerStatus
		124: -2.0, // numBeyond2
		125: 0.0,  // isSecured
		126: 0.0,  // oppositeDiff
		127: 7.0,  // numEnemyChexInFront
		128: 0.0,  // distToClosestEnemyBlotPoint
		129: 2.0,  // distToClosestEnemySecuredPoint
		130: 0.0,  // ownerStatus * isSecured
		131: 0.0,  // ownerStatus * numBeyond2
		132: 0.0,  // ownerStatus * numEnemyChexInFront
		133: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		134: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "l"
		135: 1.0, // ownerStatus
		136: 3.0, // numBeyond2
		137: 1.0, // isSecured
		138: 0.0, // oppositeDiff
		139: 7.0, // numEnemyChexInFront
		140: 0.0, // distToClosestEnemyBlotPoint
		141: 1.0, // distToClosestEnemySecuredPoint
		142: 1.0, // ownerStatus * isSecured
		143: 3.0, // ownerStatus * numBeyond2
		144: 7.0, // ownerStatus * numEnemyChexInFront
		145: 0.0, // ownerStatus * distToClosestEnemyBlotPoint
		146: 1.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "m"
		147: -1.0,  // ownerStatus
		148: 3.0,   // numBeyond2
		149: 1.0,   // isSecured
		150: 0.0,   // oppositeDiff
		151: 2.0,   // numEnemyChexInFront
		152: 0.0,   // distToClosestEnemyBlotPoint
		153: 11.0,  // distToClosestEnemySecuredPoint
		154: -1.0,  // ownerStatus * isSecured
		155: -3.0,  // ownerStatus * numBeyond2
		156: -2.0,  // ownerStatus * numEnemyChexInFront
		157: 0.0,   // ownerStatus * distToClosestEnemyBlotPoint
		158: -11.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "n"
		159: 0.0,  // ownerStatus
		160: -2.0, // numBeyond2
		161: 0.0,  // isSecured
		162: 0.0,  // oppositeDiff
		163: 2.0,  // numEnemyChexInFront
		164: 0.0,  // distToClosestEnemyBlotPoint
		165: 10.0, // distToClosestEnemySecuredPoint
		166: 0.0,  // ownerStatus * isSecured
		167: 0.0,  // ownerStatus * numBeyond2
		168: 0.0,  // ownerStatus * numEnemyChexInFront
		169: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		170: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "o"
		171: 0.0,  // ownerStatus
		172: -2.0, // numBeyond2
		173: 0.0,  // isSecured
		174: 0.0,  // oppositeDiff
		175: 2.0,  // numEnemyChexInFront
		176: 0.0,  // distToClosestEnemyBlotPoint
		177: 9.0,  // distToClosestEnemySecuredPoint
		178: 0.0,  // ownerStatus * isSecured
		179: 0.0,  // ownerStatus * numBeyond2
		180: 0.0,  // ownerStatus * numEnemyChexInFront
		181: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		182: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "p"
		183: 0.0,  // ownerStatus
		184: -2.0, // numBeyond2
		185: 0.0,  // isSecured
		186: 0.0,  // oppositeDiff
		187: 2.0,  // numEnemyChexInFront
		188: 0.0,  // distToClosestEnemyBlotPoint
		189: 8.0,  // distToClosestEnemySecuredPoint
		190: 0.0,  // ownerStatus * isSecured
		191: 0.0,  // ownerStatus * numBeyond2
		192: 0.0,  // ownerStatus * numEnemyChexInFront
		193: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		194: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "q"
		195: 1.0, // ownerStatus
		196: 1.0, // numBeyond2
		197: 1.0, // isSecured
		198: 0.0, // oppositeDiff
		199: 2.0, // numEnemyChexInFront
		200: 0.0, // distToClosestEnemyBlotPoint
		201: 7.0, // distToClosestEnemySecuredPoint
		202: 1.0, // ownerStatus * isSecured
		203: 1.0, // ownerStatus * numBeyond2
		204: 2.0, // ownerStatus * numEnemyChexInFront
		205: 0.0, // ownerStatus * distToClosestEnemyBlotPoint
		206: 7.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "r"
		207: 0.0,  // ownerStatus
		208: -2.0, // numBeyond2
		209: 0.0,  // isSecured
		210: 0.0,  // oppositeDiff
		211: 2.0,  // numEnemyChexInFront
		212: 0.0,  // distToClosestEnemyBlotPoint
		213: 6.0,  // distToClosestEnemySecuredPoint
		214: 0.0,  // ownerStatus * isSecured
		215: 0.0,  // ownerStatus * numBeyond2
		216: 0.0,  // ownerStatus * numEnemyChexInFront
		217: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		218: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "s"
		219: 1.0, // ownerStatus
		220: 3.0, // numBeyond2
		221: 1.0, // isSecured
		222: 0.0, // oppositeDiff
		223: 2.0, // numEnemyChexInFront
		224: 0.0, // distToClosestEnemyBlotPoint
		225: 5.0, // distToClosestEnemySecuredPoint
		226: 1.0, // ownerStatus * isSecured
		227: 3.0, // ownerStatus * numBeyond2
		228: 2.0, // ownerStatus * numEnemyChexInFront
		229: 0.0, // ownerStatus * distToClosestEnemyBlotPoint
		230: 5.0, // ownerStatus * distToClosestEnemySecuredPoint

		// Point "t"
		231: 0.0,  // ownerStatus
		232: -2.0, // numBeyond2
		233: 0.0,  // isSecured
		234: 0.0,  // oppositeDiff
		235: 2.0,  // numEnemyChexInFront
		236: 0.0,  // distToClosestEnemyBlotPoint
		237: 4.0,  // distToClosestEnemySecuredPoint
		238: 0.0,  // ownerStatus * isSecured
		239: 0.0,  // ownerStatus * numBeyond2
		240: 0.0,  // ownerStatus * numEnemyChexInFront
		241: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		242: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "u"
		243: 0.0,  // ownerStatus
		244: -2.0, // numBeyond2
		245: 0.0,  // isSecured
		246: 0.0,  // oppositeDiff
		247: 2.0,  // numEnemyChexInFront
		248: 0.0,  // distToClosestEnemyBlotPoint
		249: 3.0,  // distToClosestEnemySecuredPoint
		250: 0.0,  // ownerStatus * isSecured
		251: 0.0,  // ownerStatus * numBeyond2
		252: 0.0,  // ownerStatus * numEnemyChexInFront
		253: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		254: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "v"
		255: 0.0,  // ownerStatus
		256: -2.0, // numBeyond2
		257: 0.0,  // isSecured
		258: 0.0,  // oppositeDiff
		259: 2.0,  // numEnemyChexInFront
		260: 0.0,  // distToClosestEnemyBlotPoint
		261: 2.0,  // distToClosestEnemySecuredPoint
		262: 0.0,  // ownerStatus * isSecured
		263: 0.0,  // ownerStatus * numBeyond2
		264: 0.0,  // ownerStatus * numEnemyChexInFront
		265: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		266: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "w"
		267: 0.0,  // ownerStatus
		268: -2.0, // numBeyond2
		269: 0.0,  // isSecured
		270: 0.0,  // oppositeDiff
		271: 2.0,  // numEnemyChexInFront
		272: 0.0,  // distToClosestEnemyBlotPoint
		273: 1.0,  // distToClosestEnemySecuredPoint
		274: 0.0,  // ownerStatus * isSecured
		275: 0.0,  // ownerStatus * numBeyond2
		276: 0.0,  // ownerStatus * numEnemyChexInFront
		277: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		278: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint

		// Point "x"
		279: -1.0, // ownerStatus
		280: 0.0,  // numBeyond2
		281: 1.0,  // isSecured
		282: 0.0,  // oppositeDiff
		283: 0.0,  // numEnemyChexInFront
		284: 0.0,  // distToClosestEnemyBlotPoint
		285: 0.0,  // distToClosestEnemySecuredPoint
		286: -1.0, // ownerStatus * isSecured
		287: 0.0,  // ownerStatus * numBeyond2
		288: 0.0,  // ownerStatus * numEnemyChexInFront
		289: 0.0,  // ownerStatus * distToClosestEnemyBlotPoint
		290: 0.0,  // ownerStatus * distToClosestEnemySecuredPoint
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
