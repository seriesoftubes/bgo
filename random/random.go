// Package random generates random numbers in a threadsafe way.
package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	mu  sync.Mutex
	gen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Float32Between(min, max float32) float32 {
	defer mu.Unlock()
	mu.Lock()
	return gen.Float32()*(max-min) + min
}

func Float64() float64 {
	defer mu.Unlock()
	mu.Lock()
	return gen.Float64()
}

func IntBetween(min, max int) int {
	defer mu.Unlock()
	mu.Lock()
	return gen.Intn(max-min+1) + min
}
func IntUpTo(exclusiveMax int) int    { return IntBetween(0, exclusiveMax-1) }
func Uint8Between(min, max int) uint8 { return uint8(IntBetween(min, max)) }
