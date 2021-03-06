package turn

import (
	"reflect"
	"testing"

	"github.com/seriesoftubes/bgo/game/plyr"
)

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
			Turn{Move{plyr.PCC, 'j', 5}: 1, Move{plyr.PCC, 'a', 1}: 1},
			Turn{Move{plyr.PCC, 'j', 5}: 1, Move{plyr.PCC, 'a', 1}: 1},
		},
		{
			Turn{Move{plyr.PCC, 'j', 5}: 2, Move{plyr.PCC, 'a', 1}: 1},
			Turn{Move{plyr.PCC, 'j', 5}: 2, Move{plyr.PCC, 'a', 1}: 1},
		},
	}
	for _, c := range cases {
		got := c.turn.Copy()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("expected copying turn %v to produce %v but got %v", c.turn, c.want, got)
		}
	}
}

func TestSerdeTurn(t *testing.T) {
	cases := []struct {
		turn    Turn
		want    string
		wantArr TurnArray
	}{
		// abcdefghijklmnopqrstuvwxyz
		// 01234567890123456789012345
		// 0         10        20
		{
			Turn{Move{plyr.PCC, 'j', 5}: 1, Move{plyr.PCC, 'a', 1}: 1},
			"X;a1;j5",
			TurnArray{MoveArray{1, 0, 1, 1}, MoveArray{1, 9, 5, 1}},
		},
		{
			Turn{Move{plyr.PCC, 'j', 1}: 4},
			"X;j1;j1;j1;j1",
			TurnArray{MoveArray{1, 9, 1, 4}},
		},
		{
			Turn{Move{plyr.PC, 'j', 1}: 4},
			"O;j1;j1;j1;j1",
			TurnArray{MoveArray{0, 9, 1, 4}},
		},
		{
			Turn{Move{plyr.PC, 'a', 2}: 2, Move{plyr.PC, 'b', 2}: 2},
			"O;a2;a2;b2;b2",
			TurnArray{MoveArray{0, 0, 2, 2}, MoveArray{0, 1, 2, 2}},
		},
		{
			Turn{Move{plyr.PC, 't', 5}: 2, Move{plyr.PC, 'h', 5}: 2},
			"O;h5;h5;t5;t5",
			TurnArray{MoveArray{0, 7, 5, 2}, MoveArray{0, 19, 5, 2}},
		},
	}
	for _, c := range cases {
		got := c.turn.String()
		if got != c.want {
			t.Errorf("turn %v not serialized as %v; got %v", c.turn, c.want, got)
		}
		gotArr := c.turn.Arrayify()
		if !reflect.DeepEqual(gotArr, c.wantArr) {
			t.Errorf("turn %v not Arrayifyed as %v; got %v", c.turn, c.wantArr, gotArr)
		}

		if deser, err := DeserializeTurn(got); err != nil {
			t.Errorf("could not deserialize the serialized turn %s: %v", got, err)
		} else if !reflect.DeepEqual(deser, c.turn) {
			t.Errorf("unexpected deserialized turn %v for string %s", deser, got)
		}
	}
}

// test for hitting enemy checker and possibly moving on afterwards
// test for u can only move 1 dice amt but in different places
