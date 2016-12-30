package main

import (
	"fmt"
	//"sync"
	"sync/atomic"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
)

const trainings = 2000

func main() {
	start := time.Now()
	gamesPlayed := uint64(0)
	// var wg sync.WaitGroup

	// for i := 0; i < 1; i++ {
	// 	wg.Add(1)
	// 	go func() {
	mgr := ctrl.New(false)
	for i := 0; i < trainings; i++ {
		mgr.PlayOneGame(0, false) // Play 1 game with 0 humans and don't stop learning!

		atomic.AddUint64(&gamesPlayed, 1)
		if ct := atomic.LoadUint64(&gamesPlayed); ct%500 == 0 {
			fmt.Println(time.Now(), "trained on", ct, "games")
		}
	}
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()

	fmt.Println("trained", gamesPlayed, "times in", time.Since(start))

	fmt.Println("GamesPlayed\tAvgVariance")
	for i, v := range ctrl.NNVariances {
		fmt.Printf("%v\t%v\n", i+1, v)
	}

	mgr = ctrl.New(true)
	mgr.PlayOneGame(1, true /* stopLearning=true */)
}
