package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
	"github.com/seriesoftubes/bgo/learn"
)

const trainings = 250500

func main() {
	// Shared resources across all goroutines.
	qs := learn.NewQContainer()
	start := time.Now()
	gamesPlayed := uint64(0)
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			mgr := ctrl.New(qs, false)
			for i := 0; i < trainings; i++ {
				mgr.PlayOneGame(0, false) // Play 1 game with 0 humans and don't stop learning!

				atomic.AddUint64(&gamesPlayed, 1)
				if ct := atomic.LoadUint64(&gamesPlayed); ct%500 == 0 {
					fmt.Println(time.Now(), "trained on", ct, "games")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("trained", trainings, "times in", time.Since(start))

	mgr := ctrl.New(qs, false)
	mgr.PlayOneGame(1 /* stopLearning=true */, true)
}
