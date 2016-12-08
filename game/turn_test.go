package game

import (
	"reflect"
	"sort"
	"testing"

	"github.com/seriesoftubes/bgo/constants"
)

type stringSet map[string]bool

func newStringSet(strs []string) stringSet {
	out := stringSet{}
	for _, s := range strs {
		out[s] = true
	}
	return out
}

func (ss stringSet) copy() stringSet {
	out := stringSet{}
	for s := range ss {
		out[s] = true
	}
	return out
}

func (ss stringSet) values() []string {
	var out []string
	for s := range ss {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func (ss stringSet) subtract(o stringSet) stringSet {
	orig := ss.copy()
	for s := range o {
		delete(orig, s)
	}
	return orig
}

func TestCopyTurn(t *testing.T) {
	cases := []struct {
		turn Turn
		want Turn
	}{
		{
			Turn{},
			Turn{},
		},
		{
			Turn{Move{PCC, "j", 5}: 1, Move{PCC, "a", 1}: 1},
			Turn{Move{PCC, "j", 5}: 1, Move{PCC, "a", 1}: 1},
		},
		{
			Turn{Move{PCC, "j", 5}: 2, Move{PCC, "a", 1}: 1},
			Turn{Move{PCC, "j", 5}: 2, Move{PCC, "a", 1}: 1},
		},
	}
	for _, c := range cases {
		got := c.turn.copy()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("expected copying turn %v to produce %v but got %v", c.turn, c.want, got)
		}
	}
}

func TestSerdeTurn(t *testing.T) {
	cases := []struct {
		turn Turn
		want string
	}{
		{
			Turn{Move{PCC, "j", 5}: 1, Move{PCC, "a", 1}: 1},
			"X;a1;j5",
		},
		{
			Turn{Move{PCC, "j", 1}: 4},
			"X;j1;j1;j1;j1",
		},
		{
			Turn{Move{PC, "j", 1}: 4},
			"O;j1;j1;j1;j1",
		},
		{
			Turn{Move{PC, "a", 2}: 2, Move{PC, "b", 2}: 2},
			"O;a2;a2;b2;b2",
		},
		{
			Turn{Move{PC, "t", 5}: 2, Move{PC, "h", 5}: 2},
			"O;h5;h5;t5;t5",
		},
	}
	for _, c := range cases {
		got := c.turn.String()
		if got != c.want {
			t.Errorf("turn %v not serialized as %v; got %v", c.turn, c.want, got)
		}

		if deser, err := DeserializeTurn(got); err != nil {
			t.Errorf("could not deserialize the serialized turn %s: %v", got, err)
		} else if !reflect.DeepEqual(deser, c.turn) {
			t.Errorf("unexpected deserialized turn %v for string %s", deser, got)
		}
	}
}

func TestValidTurns(t *testing.T) {
	cases := []struct {
		player *Player
		roll   Roll
		want   []string // List of stringified turns
	}{
		{
			PCC, Roll{5, 4},
			[]string{
				"X;a4;e5",
				"X;a4;l5",
				"X;a4;q5",
				"X;l5;q4",
				"X;l5;s4",
				"X;l4;p5",
				"X;l4;l5",
				"X;l4;q5",
				"X;q5;s4",
				"X;q4;q5",
			},
		},
		{
			PCC, Roll{4, 5},
			[]string{
				"X;a4;e5",
				"X;a4;l5",
				"X;a4;q5",
				"X;l5;q4",
				"X;l5;s4",
				"X;l4;p5",
				"X;l4;l5",
				"X;l4;q5",
				"X;q5;s4",
				"X;q4;q5",
			},
		},
		{
			PCC, Roll{5, 5},
			[]string{
				"X;l5;l5;l5;l5",
				"X;l5;l5;l5;q5",
				"X;l5;l5;q5;q5",
				"X;l5;q5;q5;q5",
			},
		},
		{
			PC, Roll{5, 5},
			[]string{
				"O;m5;m5;m5;m5",
				"O;h5;h5;h5;m5",
				"O;h5;h5;m5;m5",
				"O;h5;m5;m5;m5",
			},
		},
	}
	for _, c := range cases {
		b := &Board{}
		b.setUp()
		/*
		   Board is:
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
		wants := newStringSet(c.want)
		turns := ValidTurns(b, &c.roll, c.player)
		gots := stringSet{}
		for ts := range turns {
			gots[ts] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}

func TestValidTurnsTwoOnTheBar(t *testing.T) {
	cases := []struct {
		player *Player
		roll   Roll
		want   []string // List of stringified turns
	}{
		{
			PCC, Roll{6, 6}, // Can only land on spaces a-e with 2 on the bar.
			[]string{},
		},
		{
			PCC, Roll{1, 6},
			[]string{
				"X;y1",
			},
		},
		{
			PCC, Roll{2, 6},
			[]string{
				"X;y2",
			},
		},
		{
			PCC, Roll{2, 3},
			[]string{
				"X;y2;y3",
			},
		},
	}
	for _, c := range cases {
		b := &Board{}
		b.setUp()
		b.Points[alpha2Num["a"]].NumCheckers = 0
		b.BarCC = 2
		/*
		    Board is:
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
		                   O |w|    O           X
		    -  -  -  -  -  O |w| -  O  -  -  -  X
		   =======================================
		    a  b  c  d  e  f     g  h  i  j  k  l

		   The bar
		   y X's: XX
		   z O's: -
		*/
		wants := newStringSet(c.want)
		turns := ValidTurns(b, &c.roll, c.player)
		gots := stringSet{}
		for ts := range turns {
			gots[ts] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}

func TestValidTurnsOneOnTheBar(t *testing.T) {
	cases := []struct {
		player *Player
		roll   Roll
		want   []string // List of stringified turns
	}{
		{
			PCC, Roll{6, 6}, // Can only land on spaces a-e.
			[]string{},
		},
		{
			PCC, Roll{1, 6},
			[]string{
				"X;a6;y1",
				"X;l6;y1",
				"X;q6;y1",
			},
		},
		{
			PCC, Roll{2, 6},
			[]string{
				"X;a6;y2",
				"X;l6;y2",
				"X;q6;y2",
			},
		},
		{
			PCC, Roll{2, 3},
			[]string{
				"X;a2;y3",
				"X;a3;y2",
				"X;b3;y2",
				"X;c2;y3",
				"X;l2;y3",
				"X;l3;y2",
				"X;q2;y3",
				"X;q3;y2",
				"X;s2;y3",
				"X;s3;y2",
			},
		},
	}
	for _, c := range cases {
		b := &Board{}
		b.setUp()
		b.Points[alpha2Num["a"]].NumCheckers = 1
		b.BarCC = 1
		/*
		    Board is:
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
		                   O |w|    O           X
		    X  -  -  -  -  O |w| -  O  -  -  -  X
		   =======================================
		    a  b  c  d  e  f     g  h  i  j  k  l

		   The bar
		   y X's: X
		   z O's: -
		*/
		wants := newStringSet(c.want)
		turns := ValidTurns(b, &c.roll, c.player)
		gots := stringSet{}
		for ts := range turns {
			gots[ts] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}

// Tests that, with all chex in your home, you can bear stuff off.
func TestValidTurnsBearOff(t *testing.T) {
	cases := []struct {
		player *Player
		roll   Roll
		want   []string // List of stringified turns
	}{
		{
			PC, Roll{6, 6},
			[]string{"O;f6;f6;f6;f6"},
		},
		{
			PC, Roll{6, 5},
			[]string{"O;e5;f6"},
		},
		{
			PC, Roll{6, 4},
			[]string{
				"O;d4;f6",
				"O;f4;f6",
			},
		},
		{
			PC, Roll{5, 5},
			[]string{"O;e5;e5;e5;e5"},
		},
	}
	for _, c := range cases {
		b := &Board{}
		b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{PCC, 2}, {PC, 1}, {PC, 2}, {PC, 2}, {PC, 5}, {PC, 5}, {}, {}, {}, {}, {}, {PCC, 5},
			{}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {},
			//                                                        clockwise player in top-left.
		}
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
		wants := newStringSet(c.want)
		turns := ValidTurns(b, &c.roll, c.player)
		gots := stringSet{}
		for ts := range turns {
			gots[ts] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}

// Tests that, even when you start with 1 chex outside your home, you can bear stuff off.
func TestValidTurnsBearOffStartingFromOutside(t *testing.T) {
	cases := []struct {
		player *Player
		roll   Roll
		want   []string // List of stringified turns
	}{
		{
			PC, Roll{6, 6},
			[]string{"O;f6;f6;f6;h6"},
		},
		{
			PC, Roll{6, 5},
			[]string{
				"O;e5;h6",
				"O;f6;h5",
			},
		},
		{
			PC, Roll{6, 4},
			[]string{
				"O;d4;h6",
				"O;f4;h6",
				"O;f6;h4",
			},
		},
		{
			PC, Roll{5, 5},
			[]string{"O;e5;e5;e5;h5"},
		},
	}
	for _, c := range cases {
		b := &Board{}
		b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
			// counter-clockwise player is in bottom-left.
			{PCC, 2}, {PC, 1}, {PC, 2}, {PC, 2}, {PC, 5}, {PC, 4}, {}, {PC, 1}, {}, {}, {}, {PCC, 5},
			{}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {},
			//                                                        clockwise player in top-left.
		}
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
		                O    |w|                X
		                O  O |w|                X
		                O  O |w|                X
		    X     O  O  O  O |w|                X
		    X  O  O  O  O  O |w| -  O  -  -  -  X
		   =======================================
		    a  b  c  d  e  f     g  h  i  j  k  l
		*/
		// PCC == "X", PC = O
		wants := newStringSet(c.want)
		turns := ValidTurns(b, &c.roll, c.player)
		gots := stringSet{}
		for ts := range turns {
			gots[ts] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}

func TestWeirdTurn(t *testing.T) {
	/*
		Board
		  Player: O  Rolled: [2 2]

		   x  w  v  u  t  s     r  q  p  o  n  m
		  =======================================
		   X  O  X  O  X  X |m| -  -  -  X  X  O
		   X     X  O     X |m|
		         X        X |m|
		         X        X |m|
		         X          |m|
		                    |m|


		                    |w|
		                    |w|
		                    |w|       O
		                    |w|       O
		   O              O |w|       O
		   O  O  -  X  -  O |w| -  -  O  -  -  -
		  =======================================
		   a  b  c  d  e  f     g  h  i  j  k  l

		  The bar
		  y X's: -
		  z O's: OO
	*/
	b := &Board{}
	b.Points = &[constants.NUM_BOARD_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{PC, 2}, {PC, 1}, {}, {PCC, 1}, {}, {PC, 2}, {}, {}, {PC, 4}, {}, {}, {},
		{PC, 1}, {PCC, 1}, {PCC, 1}, {}, {}, {}, {PCC, 4}, {PCC, 1}, {PC, 2}, {PCC, 5}, {PC, 1}, {PCC, 2},
		//                                                        clockwise player in top-left.
	}
	b.BarC = 2

	wants := newStringSet([]string{
		"O;d2;f2;z2;z2",
		"O;f2;f2;z2;z2",
		"O;f2;i2;z2;z2",
		"O;f2;m2;z2;z2",
		"O;f2;w2;z2;z2",
		"O;g2;i2;z2;z2",
		"O;i2;i2;z2;z2",
		"O;i2;m2;z2;z2",
		"O;i2;w2;z2;z2",
		"O;k2;m2;z2;z2",
		"O;m2;w2;z2;z2",
		"O;w2;w2;z2;z2",
	})
	roll := &Roll{2, 2}
	turns := ValidTurns(b, roll, PC)
	gots := stringSet{}
	for ts := range turns {
		gots[ts] = true
	}

	if !reflect.DeepEqual(gots, wants) {
		extraWants := wants.subtract(gots).values()
		missingWants := gots.subtract(wants).values()
		t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", roll, PC, missingWants, extraWants)
	}
}

// test for hitting enemy checker and possibly moving on afterwards
// test for u can only move 1 dice amt but in different places
