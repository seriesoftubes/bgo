// Package nnperf stores performance info on the neural net.
package nnperf

import (
	"sync"
)

var (
	gameAverageVariances container = container{}
	gameTotalVariances   container = container{}
)

type container struct {
	sync.RWMutex
	data []float32
}

func GameAverageVariances(lastN int, waitForWrites bool) []float32 {
	if waitForWrites {
		gameAverageVariances.Lock()
		defer gameAverageVariances.Unlock()
	} else {
		gameAverageVariances.RLock()
		defer gameAverageVariances.RUnlock()
	}

	sz := len(gameAverageVariances.data)
	firstIndex := sz - lastN
	if lastN < 0 || firstIndex < 0 {
		firstIndex = 0
	}
	return append([]float32{}, gameAverageVariances.data[firstIndex:sz]...)
}

func GameTotalVariances(lastN int, waitForWrites bool) []float32 {
	if waitForWrites {
		gameTotalVariances.Lock()
		defer gameTotalVariances.Unlock()
	} else {
		gameTotalVariances.RLock()
		defer gameTotalVariances.RUnlock()
	}

	sz := len(gameTotalVariances.data)
	firstIndex := sz - lastN
	if lastN < 0 {
		firstIndex = 0
	}
	return append([]float32{}, gameTotalVariances.data[firstIndex:sz]...)
}

func AppendGameData(totalVariance float32, numTrainings uint32) {
	appendAverageVariance(totalVariance / float32(numTrainings))
	appendTotalVariance(totalVariance)
}

func appendAverageVariance(av float32) {
	gameAverageVariances.Lock()
	defer gameAverageVariances.Unlock()
	gameAverageVariances.data = append(gameAverageVariances.data, av)
}

func appendTotalVariance(tv float32) {
	gameTotalVariances.Lock()
	gameTotalVariances.data = append(gameTotalVariances.data, tv)
	gameTotalVariances.Unlock()
}
