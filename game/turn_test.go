package game

import (
	"reflect"
	"sort"
	"testing"
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
			Turn{},
			"",
		},
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

func TestTurnPerms(t *testing.T) {
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
		turns, err := TurnPerms(b, &c.roll, c.player)
		if err != nil {
			t.Errorf("error producing turnPerms for roll %v and player %s: %v", c.roll, *c.player, err)
		}

		gots := stringSet{}
		for _, t := range turns {
			gots[t.String()] = true
		}

		if !reflect.DeepEqual(gots, wants) {
			extraWants := wants.subtract(gots).values()
			missingWants := gots.subtract(wants).values()
			t.Errorf("TestTurnPerms bug for roll %v and player %s.\nwants is missing %v,\nwants has extra %v", c.roll, *c.player, missingWants, extraWants)
		}
	}
}
