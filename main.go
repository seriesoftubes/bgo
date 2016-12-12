package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
	"github.com/seriesoftubes/bgo/learn"
)

const (
	trainings     = 250500
	qvalsFileName = "/tmp/bgo_qvals.gob"
)

func readQs() *learn.QContainer {
	f, err := os.Open(qvalsFileName)
	if err != nil {
		return learn.NewQContainer()
	}
	defer f.Close()

	qc, err := learn.DeserializeQContainer(f)
	if err != nil {
		panic("could not deserialize qs: " + err.Error())
	}
	return qc
}

func saveQs(qc *learn.QContainer) {
	f, err := os.Create(qvalsFileName) // always overwrites the existing file.
	if err != nil {
		panic("could not create file: " + err.Error())
	}
	defer f.Close()

	if err := qc.Serialize(f); err != nil {
		panic("couldnt Serialize qvals: " + err.Error())
	}
}

func main() {
	// Shared resources across all goroutines.
	qs := readQs()

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

	saveQs(qs)

	mgr := ctrl.New(qs, true)
	mgr.PlayOneGame(1, true /* stopLearning=true */)
}
