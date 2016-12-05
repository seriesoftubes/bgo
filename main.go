package main

import (
	"fmt"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
)

const trainings = 1000

func main() {
	ctrl := ctrl.New(false)

	start := time.Now()
	for i := 0; i < trainings; i++ {
		if i%1000 == 0 {
			fmt.Println(time.Now(), "trained on", i, "games")
		}
		ctrl.PlayOneGame(0, false)
	}
	fmt.Println("trained", trainings, "times in", time.Since(start))

	ctrl.PlayOneGame(1, true)
}
