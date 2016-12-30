package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
	"github.com/seriesoftubes/bgo/learn/nnet"
)

const (
	trainings                    = 200
	defaultNeuralNetworkFileName = "/Users/bweidenbaum/Desktop/bgo_nnet.json"
)

func loadNeuralNetwork(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("could not open neural network file %q. Skipping. %v\n", filePath, err)
		return
	}
	defer f.Close()

	if err := nnet.Load(f); err != nil {
		panic("could not deserialize neural network: " + err.Error())
	}
}

func saveNeuralNetwork(filePath string) {
	f, err := os.Create(filePath) // always overwrites the existing file.
	if err != nil {
		panic("could not create file: " + err.Error())
	}
	defer f.Close()

	if err := nnet.Save(f); err != nil {
		panic("couldnt save neural network: " + err.Error())
	}
}

func main() {
	loadNeuralNetwork(defaultNeuralNetworkFileName) // TODO: use flag instead of the default name.

	start := time.Now()
	gamesPlayed := uint64(0)
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			mgr := ctrl.New(false)
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

	fmt.Println("trained", gamesPlayed, "times in", time.Since(start))

	saveNeuralNetwork(defaultNeuralNetworkFileName)

	fmt.Println("GamesPlayed\tAvgVariance")
	for i, v := range ctrl.NNVariances {
		fmt.Printf("%v\t%v\n", i+1, v)
	}

	mgr := ctrl.New(true)
	mgr.PlayOneGame(1, true /* stopLearning=true */)
}
