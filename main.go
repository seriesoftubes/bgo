package main

import (
	"fmt"

	"github.com/seriesoftubes/bgo/ctrl"
)

func main() {
	ctrl := ctrl.New(false)

	for i := 0; i < 1000; i++ {
		if i%10 == 0 {
			fmt.Println("trained on", i, "games")
		}
		ctrl.PlayOneGame(0, false)
	}

	ctrl.PlayOneGame(1, true)
}
