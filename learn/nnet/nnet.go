// Package nnet contains a neural network that estimates the game equity given the current state of the game.
// The network itself consists of the input layer, 1 fully-connected Sigmoid layer.
// TODO: maybe add 2nd hidden ReLU layer.
package nnet

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"runtime"
	"sync"

	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

const (
	bias                 = float32(1.8)           // If you change the bias, saved weights will be erroneous and you need to retrain the network.
	numInputs            = len(state.State{}) + 1 // The `+1` is for the artificially-added bias.
	numIn2FhConnections  = numInputs * numInputs  // the "FH" aka firsthidden layer is fully-connected with all inputs.
	numFhNodes           = numInputs + 1          // the "FH" layer has 1 node per input, plus one more bias.
	numFh2OutConnections = numFhNodes             // Each "FH" node connects directly to the 1 output signal.
)

var (
	// TODO: automatically skip modifying certain weights to keep it closer to "all else equal".
	learningRate         = float32(0.00001)
	eligibilityDecayRate = float32(0.95)
	configMu             sync.RWMutex

	maxConcurrentGames int = runtime.NumCPU() * 2 // Assume a max of 2 goroutines training per CPU. This variable would be a const but that caused a compiler error.

	bigassArrays chan *[numIn2FhConnections]float32 = func() chan *[numIn2FhConnections]float32 {
		num := maxConcurrentGames * 3
		ch := make(chan *[numIn2FhConnections]float32, num)
		for i := 0; i < num; i++ {
			var arr [numIn2FhConnections]float32
			ch <- &arr
		}
		return ch
	}()

	// Weights that connect IN to FH.
	in2fhWeightsPtr *[numIn2FhConnections]float32 = (func() *[numIn2FhConnections]float32 {
		out := [numIn2FhConnections]float32{}
		for i := 0; i < numIn2FhConnections; i++ {
			out[i] = random.Float32Between(-1, 1)
		}
		return &out
	})()
	// each key in the map is a gameID, and each value is an array of previous eligibility traces-- one trace for each in2fh weight.
	in2fhWeightsPreviousEligibilityTracesByGameID map[uint32]*[numIn2FhConnections]float32 = make(map[uint32]*[numIn2FhConnections]float32, maxConcurrentGames)
	in2fhWeightsMu                                sync.RWMutex
	in2fhWeightsChunked                           []int = splitIntoChunkSizes(numIn2FhConnections, runtime.NumCPU()*3/2)

	// Weights that connect FH to OUT.
	fh2outWeights [numFh2OutConnections]float32 = (func() [numFh2OutConnections]float32 {
		out := [numFh2OutConnections]float32{}
		for i := 0; i < numFh2OutConnections; i++ {
			out[i] = random.Float32Between(-1, 1)
		}
		return out
	})()
	// each key in the map is a gameID, and each value is an array of previous eligibility traces-- one trace for each fh2out weight.
	fh2outWeightsPreviousEligibilityTracesByGameID map[uint32]*[numFh2OutConnections]float32 = make(map[uint32]*[numFh2OutConnections]float32, maxConcurrentGames)
	fh2outWeightsMu                                sync.RWMutex
)

type netConfig struct {
	GamesPlayedSoFar     uint64
	LearningRate         float32
	EligibilityDecayRate float32
	In2FhWeights         [numIn2FhConnections]float32
	Fh2OutWeights        [numFh2OutConnections]float32
}

func Save(w io.Writer, gamesPlayedSoFar uint64, waitForWrites bool) error {
	my_learningRate, my_eligibilityDecayRate := LearningParams()
	cfg := netConfig{
		GamesPlayedSoFar:     gamesPlayedSoFar,
		EligibilityDecayRate: my_eligibilityDecayRate,
		LearningRate:         my_learningRate,
	}

	if waitForWrites {
		fh2outWeightsMu.Lock()
		cfg.Fh2OutWeights = fh2outWeights
		fh2outWeightsMu.Unlock()
		in2fhWeightsMu.Lock()
		cfg.In2FhWeights = *in2fhWeightsPtr
		in2fhWeightsMu.Unlock()
	} else {
		fh2outWeightsMu.RLock()
		cfg.Fh2OutWeights = fh2outWeights
		fh2outWeightsMu.RUnlock()
		in2fhWeightsMu.RLock()
		cfg.In2FhWeights = *in2fhWeightsPtr
		in2fhWeightsMu.RUnlock()
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("JSON Encode error: %v", err)
	}
	return nil
}

func Load(r io.Reader) (uint64, error) {
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("ioutil.ReadAll(r) error: %v", err)
	}

	var cfg netConfig
	if err := json.Unmarshal(text, &cfg); err != nil {
		return 0, fmt.Errorf("json.Unmarshal error: %v", err)
	}

	configMu.Lock()
	learningRate = cfg.LearningRate
	eligibilityDecayRate = cfg.EligibilityDecayRate
	configMu.Unlock()

	if len(cfg.In2FhWeights) != numIn2FhConnections {
		return 0, fmt.Errorf("serialized network In2FH weights do not match dimensions of the one in this program. expected both to have length of %d", numIn2FhConnections)
	}
	if len(cfg.Fh2OutWeights) != numFh2OutConnections {
		return 0, fmt.Errorf("serialized network FH2Out weights do not match dimensions of the one in this program. expected both to have length of %d", numFh2OutConnections)
	}

	in2fhWeightsMu.Lock()
	for i, v := range cfg.In2FhWeights {
		(*in2fhWeightsPtr)[i] = v
	}
	in2fhWeightsMu.Unlock()

	fh2outWeightsMu.Lock()
	for i, v := range cfg.Fh2OutWeights {
		fh2outWeights[i] = v
	}
	fh2outWeightsMu.Unlock()

	return cfg.GamesPlayedSoFar, nil
}

func RemoveUselessGameData(gameID uint32) {
	in2fhWeightsMu.Lock()
	go recycleBigassArray(in2fhWeightsPreviousEligibilityTracesByGameID[gameID])
	delete(in2fhWeightsPreviousEligibilityTracesByGameID, gameID)
	in2fhWeightsMu.Unlock()

	fh2outWeightsMu.Lock()
	delete(fh2outWeightsPreviousEligibilityTracesByGameID, gameID)
	fh2outWeightsMu.Unlock()
}

func ValueEstimate(st state.State) (float32, [numFhNodes]float32) {
	var estimate float32
	var fhNodePostVals [numFhNodes]float32

	in2fhWeightsMu.RLock()
	my_in2fhWeights := *in2fhWeightsPtr
	in2fhWeightsMu.RUnlock()

	fh2outWeightsMu.RLock()
	my_fh2outWeights := fh2outWeights
	fh2outWeightsMu.RUnlock()

	var fhNodeIdx, in2fhWeightIndex int
	for ; fhNodeIdx < numFhNodes-1; fhNodeIdx++ {
		// Important note: this loop only iterates over the first N-1 FH nodes-- the Nth one needs to be the new bias node!
		var fhNodeSum float32
		for _, num := range st {
			fhNodeSum += num * my_in2fhWeights[in2fhWeightIndex]
			in2fhWeightIndex++
		}
		// Here we artificially add a bias to the input "state":
		fhNodeSum += bias * my_in2fhWeights[in2fhWeightIndex]
		in2fhWeightIndex++

		fhNodePostVal := tanh(fhNodeSum)
		fhNodePostVals[fhNodeIdx] = fhNodePostVal
		estimate += fhNodePostVal * my_fh2outWeights[fhNodeIdx]
	}
	// Now we artificially add a bias node to the FH layer:
	fhNodeSum := bias
	fhNodePostVal := tanh(fhNodeSum)
	fhNodePostVals[fhNodeIdx] = fhNodePostVal
	estimate += fhNodePostVal * my_fh2outWeights[fhNodeIdx]

	return estimate, fhNodePostVals
}

func MultiplyLearningRate(rateMultiplier float32) {
	configMu.Lock()
	defer configMu.Unlock()
	learningRate *= rateMultiplier
}

func LearningParams() (float32, float32) {
	configMu.RLock()
	defer configMu.RUnlock()
	return learningRate, eligibilityDecayRate
}

// TrainWeights back-propagates the error of an estimate against a target.
func TrainWeights(gameID uint32, st state.State, target float32) float32 {
	est, fh2outWeightsGradient, in2fhWeightsGradientPtr := weightGradients(st, target)
	valueEstimateDiff := target - est // If this diff is positive, we need to add the gradient in the positive direction. else in the negative direction.
	my_learningRate, my_eligibilityDecayRate := LearningParams()

	defer in2fhWeightsMu.Unlock()
	in2fhWeightsMu.Lock()
	defer fh2outWeightsMu.Unlock()
	fh2outWeightsMu.Lock()
	// Important: don't write to any of the global vars until these locks are acquired-- that's why very little processing could happen above this line.

	if _, ok := in2fhWeightsPreviousEligibilityTracesByGameID[gameID]; !ok {
		select {
		case in2fhWeightsPreviousEligibilityTracesByGameID[gameID] = <-bigassArrays:
		default:
			in2fhWeightsPreviousEligibilityTracesByGameID[gameID] = &([numIn2FhConnections]float32{})
		}
	}
	in2fhWeightsPreviousEligibilityTraces := in2fhWeightsPreviousEligibilityTracesByGameID[gameID]

	startIdx := 0
	var wg sync.WaitGroup
	for _, sz := range in2fhWeightsChunked { // using this global var is OK because everything is locked right now.
		wg.Add(1)
		go func(start, end int) {
			for i := start; i <= end; i++ {
				in2fhWeightDerivative := (*in2fhWeightsGradientPtr)[i]
				previousEligibilityTrace := (*in2fhWeightsPreviousEligibilityTraces)[i]
				eligibilityTrace := in2fhWeightDerivative + (my_eligibilityDecayRate * previousEligibilityTrace)
				(*in2fhWeightsPreviousEligibilityTraces)[i] = eligibilityTrace
				(*in2fhWeightsPtr)[i] += my_learningRate * valueEstimateDiff * eligibilityTrace
			}
			wg.Done()
		}(startIdx, startIdx+sz-1)
		startIdx += sz
	}

	if _, ok := fh2outWeightsPreviousEligibilityTracesByGameID[gameID]; !ok {
		fh2outWeightsPreviousEligibilityTracesByGameID[gameID] = &([numFh2OutConnections]float32{})
	}
	fh2outWeightsPreviousEligibilityTraces := fh2outWeightsPreviousEligibilityTracesByGameID[gameID]
	for i, fh2outWeightDerivative := range fh2outWeightsGradient {
		previousEligibilityTrace := (*fh2outWeightsPreviousEligibilityTraces)[i]
		eligibilityTrace := fh2outWeightDerivative + (my_eligibilityDecayRate * previousEligibilityTrace)
		(*fh2outWeightsPreviousEligibilityTraces)[i] = eligibilityTrace
		fh2outWeights[i] += my_learningRate * valueEstimateDiff * eligibilityTrace
	}

	wg.Wait()
	go recycleBigassArray(in2fhWeightsGradientPtr)

	return valueEstimateDiff * valueEstimateDiff // The variance before adjusting the weights.
}

// splitIntoStartEndIndices spits out ideal chunk sizes that would most evenly split an array into N (ideally) equally-sized chunks.
// maxChunks refers to the number of chunks you want to produce.
// It's possible that you will receive fewer chunks than you request, if there aren't enough elements in your array.
// It's impossible to receive more chunks than you request though.
func splitIntoChunkSizes(arrLen, maxChunks int) []int {
	if maxChunks < 1 {
		panic(fmt.Sprintf("`maxChunks` must be >= 1, but got %d", maxChunks))
	} else if arrLen < 0 {
		panic(fmt.Sprintf("`arrLen` must be >= 0, but got %d", arrLen))
	}

	var out []int

	if arrLen <= maxChunks {
		elementsPerChunk := 1
		for i := 0; i < arrLen; i++ {
			out = append(out, elementsPerChunk)
		}
		return out
	}

	if canBeEvenlyDivided := arrLen%maxChunks == 0; canBeEvenlyDivided {
		elementsPerChunk := arrLen / maxChunks
		for i := 0; i < maxChunks; i++ {
			out = append(out, elementsPerChunk)
		}
		return out
	}

	ratio := float64(arrLen) / float64(maxChunks)
	ceil, floor := int(math.Ceil(ratio)), int(math.Floor(ratio))
	dominantNumEls := ceil
	if math.Abs(float64(floor*maxChunks-arrLen)) < math.Abs(float64(ceil*maxChunks-arrLen)) {
		dominantNumEls = floor
	}

	out = append(out, dominantNumEls)
	return append(out, splitIntoChunkSizes(arrLen-dominantNumEls, maxChunks-1)...)
}

func weightGradients(st state.State, target float32) (float32, [numFh2OutConnections]float32, *[numIn2FhConnections]float32) {
	var (
		in2fhWeightIndex        int // tracks which in2fhWeight we're analyzing.
		fh2outWeightsGradient   [numFh2OutConnections]float32
		in2fhWeightsGradientPtr *[numIn2FhConnections]float32
	)

	select {
	case in2fhWeightsGradientPtr = <-bigassArrays:
	default:
		in2fhWeightsGradientPtr = &([numIn2FhConnections]float32{})
	}
	in2fhWeightsGradient := *in2fhWeightsGradientPtr

	est, fhNodePostVals := ValueEstimate(st)

	fh2outWeightsMu.RLock()
	my_fh2outWeights := fh2outWeights
	fh2outWeightsMu.RUnlock()

	for fhNodeIdx := 0; fhNodeIdx < numFhNodes; fhNodeIdx++ {
		dEstimate_wrt_fh2outWeight := fhNodePostVals[fhNodeIdx] // derive this wrt weight1: `(fhNodePostVal1*weight1 + fhNodePostVal2*weight2 + ...)`.
		fh2outWeightsGradient[fhNodeIdx] = dEstimate_wrt_fh2outWeight
		if fhNodeIdx == numFhNodes-1 {
			break // The final FH node is a standalone bias node with no weights coming into it-- that's why we skip those in2fh weights for the final FH node.
		}
		// From here down, we're talking about the incoming in2fh weights that connect to this FH node.

		dEstimate_wrt_fhNodePostVal := my_fh2outWeights[fhNodeIdx]                                   // derive this wrt fhPostVal1: `estimate = (w1*fhPostVal1 + w2*fhPostVal2 + ...)`.
		dFhNodePostVal_wrt_fhNodePreVal := 1 - dEstimate_wrt_fh2outWeight*dEstimate_wrt_fh2outWeight // derivative of tanh(x) is: 1 - tanh(x)*tanh(x)
		for _, dFhNodePreVal_wrt_in2fhWeight := range st {                                           // derive ths wrt w1: est = (w1*inp1 + w2*inp2 + ...)
			dEstimate_wrt_in2fhWeight := dEstimate_wrt_fhNodePostVal * dFhNodePostVal_wrt_fhNodePreVal * dFhNodePreVal_wrt_in2fhWeight
			in2fhWeightsGradient[in2fhWeightIndex] = dEstimate_wrt_in2fhWeight
			in2fhWeightIndex++
		}
		// and 1 more for the bias one.
		dFhNodePreVal_wrt_in2fhWeight := bias
		dEstimate_wrt_in2fhWeight := dEstimate_wrt_fhNodePostVal * dFhNodePostVal_wrt_fhNodePreVal * dFhNodePreVal_wrt_in2fhWeight
		in2fhWeightsGradient[in2fhWeightIndex] = dEstimate_wrt_in2fhWeight
		in2fhWeightIndex++
	}

	return est, fh2outWeightsGradient, in2fhWeightsGradientPtr
}

func sigmoid(x float32) float32 { return float32(1.0 / (1.0 + math.Exp(float64(-x)))) }
func tanh(x float32) float32    { return float32(math.Tanh(float64(x))) }

func recycleBigassArray(arrPtr *[numIn2FhConnections]float32) {
	arr := *arrPtr
	for i := range arr {
		arr[i] = 0.0 // Must be all zeroed out for reuse.
	}
	bigassArrays <- arrPtr
}
