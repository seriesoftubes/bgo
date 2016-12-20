package turngen

import (
	"reflect"
	"sort"
	"testing"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
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

func TestValidTurns(t *testing.T) {
	cases := []struct {
		player plyr.Player
		roll   game.Roll
		want   []string // List of stringified turns
	}{
		{
			plyr.PCC, game.Roll{5, 4},
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
			plyr.PCC, game.Roll{4, 5},
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
			plyr.PCC, game.Roll{5, 5},
			[]string{
				"X;l5;l5;l5;l5",
				"X;l5;l5;l5;q5",
				"X;l5;l5;q5;q5",
				"X;l5;q5;q5;q5",
			},
		},
		{
			plyr.PC, game.Roll{5, 5},
			[]string{
				"O;m5;m5;m5;m5",
				"O;h5;h5;h5;m5",
				"O;h5;h5;m5;m5",
				"O;h5;m5;m5;m5",
			},
		},
	}
	for _, c := range cases {
		b := &game.Board{}
		b.SetUp()
		/*
		   game.Board is:
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
		turns := ValidTurns(b, c.roll, c.player)
		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, string(c.player), missingWants, extraWants)
		}
	}
}

func TestValidTurnsTwoOnTheBar(t *testing.T) {
	cases := []struct {
		player plyr.Player
		roll   game.Roll
		want   []string // List of stringified turns
	}{
		{
			plyr.PCC, game.Roll{6, 6}, // Can only land on spaces a-e with 2 on the bar.
			[]string{},
		},
		{
			plyr.PCC, game.Roll{1, 6},
			[]string{
				"X;y1",
			},
		},
		{
			plyr.PCC, game.Roll{2, 6},
			[]string{
				"X;y2",
			},
		},
		{
			plyr.PCC, game.Roll{2, 3},
			[]string{
				"X;y2;y3",
			},
		},
	}
	for _, c := range cases {
		b := &game.Board{}
		b.SetUp()
		b.Points[constants.Alpha2Num['a']].NumCheckers = 0
		b.BarCC = 2
		/*
		    game.Board is:
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
		turns := ValidTurns(b, c.roll, c.player)
		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, string(c.player), missingWants, extraWants)
		}
	}
}

func TestValidTurnsOneOnTheBar(t *testing.T) {
	cases := []struct {
		player plyr.Player
		roll   game.Roll
		want   []string // List of stringified turns
	}{
		{
			plyr.PCC, game.Roll{6, 6}, // Can only land on spaces a-e.
			[]string{},
		},
		{
			plyr.PCC, game.Roll{1, 6},
			[]string{
				"X;a6;y1",
				"X;l6;y1",
				"X;q6;y1",
			},
		},
		{
			plyr.PCC, game.Roll{2, 6},
			[]string{
				"X;a6;y2",
				"X;l6;y2",
				"X;q6;y2",
			},
		},
		{
			plyr.PCC, game.Roll{2, 3},
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
		b := &game.Board{}
		b.SetUp()
		b.Points[constants.Alpha2Num['a']].NumCheckers = 1
		b.BarCC = 1
		/*
		    game.Board is:
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
		turns := ValidTurns(b, c.roll, c.player)
		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, string(c.player), missingWants, extraWants)
		}
	}
}

// Tests that, with all chex in your home, you can bear stuff off.
func TestValidTurnsBearOff(t *testing.T) {
	cases := []struct {
		player plyr.Player
		roll   game.Roll
		want   []string // List of stringified turns
	}{
		{
			plyr.PC, game.Roll{6, 6},
			[]string{"O;f6;f6;f6;f6"},
		},
		{
			plyr.PC, game.Roll{6, 5},
			[]string{"O;e5;f6"},
		},
		{
			plyr.PC, game.Roll{6, 4},
			[]string{
				"O;d4;f6",
				"O;f4;f6",
			},
		},
		{
			plyr.PC, game.Roll{5, 5},
			[]string{"O;e5;e5;e5;e5"},
		},
	}
	for _, c := range cases {
		b := &game.Board{}
		b.Points = &[constants.NUM_BOARD_POINTS]*game.BoardPoint{
			// counter-clockwise player is in bottom-left.
			{plyr.PCC, 2}, {plyr.PC, 1}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 5}, {plyr.PC, 5}, {}, {}, {}, {}, {}, {plyr.PCC, 5},
			{}, {}, {}, {}, {plyr.PCC, 3}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {},
			//                                                        clockwise player in top-left.
		}
		/* game.Board looks like:
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
		wants := newStringSet(c.want)
		turns := ValidTurns(b, c.roll, c.player)
		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, string(c.player), missingWants, extraWants)
		}
	}
}

// Tests that, even when you start with 1 chex outside your home, you can bear stuff off.
func TestValidTurnsBearOffStartingFromOutside(t *testing.T) {
	cases := []struct {
		player plyr.Player
		roll   game.Roll
		want   []string // List of stringified turns
	}{
		{
			plyr.PC, game.Roll{6, 6},
			[]string{"O;f6;f6;f6;h6"},
		},
		{
			plyr.PC, game.Roll{6, 5},
			[]string{
				"O;e5;h6",
				"O;f6;h5",
			},
		},
		{
			plyr.PC, game.Roll{6, 4},
			[]string{
				"O;d4;h6",
				"O;f4;h6",
				"O;f6;h4",
			},
		},
		{
			plyr.PC, game.Roll{5, 5},
			[]string{"O;e5;e5;e5;h5"},
		},
	}
	for _, c := range cases {
		b := &game.Board{}
		b.Points = &[constants.NUM_BOARD_POINTS]*game.BoardPoint{
			// counter-clockwise player is in bottom-left.
			{plyr.PCC, 2}, {plyr.PC, 1}, {plyr.PC, 2}, {plyr.PC, 2}, {plyr.PC, 5}, {plyr.PC, 4}, {}, {plyr.PC, 1}, {}, {}, {}, {plyr.PCC, 5},
			{}, {}, {}, {}, {plyr.PCC, 3}, {}, {plyr.PCC, 5}, {}, {}, {}, {}, {},
			//                                                        clockwise player in top-left.
		}
		/* game.Board looks like:
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
		// plyr.PCC == "X", plyr.PC = O
		wants := newStringSet(c.want)
		turns := ValidTurns(b, c.roll, c.player)
		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, string(c.player), missingWants, extraWants)
		}
	}
}

func TestWeirdTurn(t *testing.T) {
	/*
	   game.Board
	     Player: O  game.Rolled: [2 2]

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
	b := &game.Board{}
	b.Points = &[constants.NUM_BOARD_POINTS]*game.BoardPoint{
		// counter-clockwise player is in bottom-left.
		{plyr.PC, 2}, {plyr.PC, 1}, {}, {plyr.PCC, 1}, {}, {plyr.PC, 2}, {}, {}, {plyr.PC, 4}, {}, {}, {},
		{plyr.PC, 1}, {plyr.PCC, 1}, {plyr.PCC, 1}, {}, {}, {}, {plyr.PCC, 4}, {plyr.PCC, 1}, {plyr.PC, 2}, {plyr.PCC, 5}, {plyr.PC, 1}, {plyr.PCC, 2},
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
	roll := game.Roll{2, 2}
	turns := ValidTurns(b, roll, plyr.PC)
	gots := stringSet{}
	for _, t := range turns {
		gots[t.String()] = true
	}

	if !reflect.DeepEqual(gots, wants) {
		extraWants := wants.subtract(gots).values()
		missingWants := gots.subtract(wants).values()
		t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", roll, plyr.PC, missingWants, extraWants)
	}
}
