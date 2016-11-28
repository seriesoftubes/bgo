package game

import (
	"reflect"
	"testing"
)

func TestTurnEquals(t *testing.T) {
	cases := []struct {
		t1   Turn
		t2   Turn
		want bool
	}{
		{
			t1:   Turn{},
			t2:   Turn{},
			want: true,
		},
		{
			t1:   Turn{{PC, "a", 1}},
			t2:   Turn{},
			want: false,
		},
		{
			t1:   Turn{{PCC, "a", 1}, {PCC, "a", 2}},
			t2:   Turn{{PCC, "a", 1}, {PCC, "a", 2}},
			want: true,
		},
		{
			t1:   Turn{{PCC, "a", 1}, {PCC, "a", 2}},
			t2:   Turn{{PCC, "a", 2}, {PCC, "a", 1}},
			want: true,
		},
		{
			t1:   Turn{{PCC, "a", 1}, {PCC, "b", 2}},
			t2:   Turn{{PCC, "b", 2}, {PCC, "a", 1}},
			want: true,
		},
		{
			t1:   Turn{{PCC, "c", 1}, {PCC, "b", 2}},
			t2:   Turn{{PCC, "b", 2}, {PCC, "a", 1}},
			want: false,
		},
		{
			t1:   Turn{{PC, "a", 1}, {PC, "b", 2}},
			t2:   Turn{{PCC, "a", 1}, {PCC, "b", 2}},
			want: false,
		},
		{
			t1:   Turn{{PC, "a", 1}, {PC, "b", 2}, {PC, "c", 3}},
			t2:   Turn{{PC, "a", 1}, {PC, "b", 2}},
			want: false,
		},
		{
			t1:   Turn{{PC, "a", 1}, {PC, "b", 2}, {PC, "c", 3}},
			t2:   Turn{{PC, "a", 1}, {PC, "b", 2}, {PC, "c", 3}},
			want: true,
		},
		{
			t1:   Turn{{PC, "a", 1}, {PC, "b", 2}, {PC, "c", 3}},
			t2:   Turn{{PC, "a", 3}, {PC, "b", 2}, {PC, "c", 1}},
			want: false,
		},
		{
			t1:   Turn{{PC, "a", 1}, {PC, "b", 2}, {PC, "c", 3}},
			t2:   Turn{{PC, "c", 3}, {PC, "b", 2}, {PC, "a", 1}},
			want: true,
		},
	}
	for _, c := range cases {
		res1, res2 := c.t1.Equals(c.t2), c.t2.Equals(c.t1)
		if res1 != c.want {
			t.Errorf("expected t1 %v == t2 %v to be %v but got %v", c.t1, c.t2, c.want, res1)
		} else if res2 != c.want {
			t.Errorf("expected t2 %v == t1 %v to be %v but got %v", c.t2, c.t1, c.want, res2)
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
				"a4;e5",
				"a4;l5",
				"a4;q5",
				"l5;a4", // dupe! how to de-dupe this? or maybe not necessary.
				"l5;l4",
				"l5;q4",
				"l5;s4",
				"l4;p5",
				"l4;l5",
				"l4;q5",
				"q5;a4",
				"q5;l4",
				"q5;q4",
				"q5;s4",
				"q4;l5",
				"q4;q5",
				"s4;l5",
				"s4;q5",
			},
		},
	}
	for _, c := range cases {
		wants := map[string]bool{}
		for _, permString := range c.want {
			wants[permString] = true
		}

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

		perms, permStrings := TurnPerms(b, &c.roll, c.player), map[string]bool{}
		for _, p := range perms {
			permStrings[p.String()] = true
		}

		if !reflect.DeepEqual(permStrings, wants) {
			t.Errorf("TestTurnPerms bug for roll %v and player %s. expected %v got %v", c.roll, *c.player, wants, permStrings)
		}
	}
}
