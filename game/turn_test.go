package game

import (
	"fmt"
	"testing"
)

func TestTurnPerms(t *testing.T) {
	b := &Board{}
	b.setUp()
	r := &Roll{1, 6}
	perms := TurnPerms(b, r, PCC)
	for _, p := range perms {
		fmt.Println(p)
	}
}
