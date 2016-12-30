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
	bias                 = float32(3.0) // If you change the bias, saved weights will be erroneous and you need to retrain the network.
	numInputs            = len(state.State{}) + 1
	numIn2FhConnections  = numInputs * numInputs // the "FH" aka firsthidden layer is fully-connected with all inputs.
	numFhNodes           = numInputs + 1         // the "FH" layer has 1 node per input, plus one more bias.
	numFh2OutConnections = numFhNodes            // Each "FH" node connects directly to the 1 output signal.
)

var (
	maxConcurrentGames int = runtime.NumCPU() * 2 // Assume a max of 2 goroutines training per CPU. This variable would be a const but that caused a compiler error.

	// Weights that connect IN to FH.
	in2fhWeights [numIn2FhConnections]float32 = (func() [numIn2FhConnections]float32 {
		out := [numIn2FhConnections]float32{}
		for i := 0; i < numIn2FhConnections; i++ {
			out[i] = random.Float32Between(-1, 1)
		}
		return out
	})()
	// each key in the map is a gameID, and each value is an array of previous eligibility traces-- one trace for each in2fh weight.
	in2fhWeightsPreviousEligibilityTracesByGameID map[uint32]*[numIn2FhConnections]float32 = make(map[uint32]*[numIn2FhConnections]float32, maxConcurrentGames)
	in2fhWeightsMu                                sync.RWMutex

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

type jsonifyable struct {
	In2FhWeights  [numIn2FhConnections]float32
	Fh2OutWeights [numFh2OutConnections]float32
}

func Save(w io.Writer) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(jsonifyable{in2fhWeights, fh2outWeights}); err != nil {
		return fmt.Errorf("JSON Encode error: %v", err)
	}
	return nil
}

func Load(r io.Reader) error {
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll(r) error: %v", err)
	}

	var j jsonifyable
	if err := json.Unmarshal(text, &j); err != nil {
		return fmt.Errorf("json.Unmarshal error: %v", err)
	}

	if len(j.In2FhWeights) != len(in2fhWeights) {
		return fmt.Errorf("serialized network does not match dimensions of the one in this program. expected both to have length of %d", len(in2fhWeights))
	}
	if len(j.Fh2OutWeights) != len(fh2outWeights) {
		return fmt.Errorf("serialized network does not match dimensions of the one in this program. expected both to have length of %d", len(fh2outWeights))
	}

	for i, v := range j.In2FhWeights {
		in2fhWeights[i] = v
	}
	for i, v := range j.Fh2OutWeights {
		fh2outWeights[i] = v
	}
	return nil
}

func RemoveUselessGameData(gameID uint32) {
	in2fhWeightsMu.Lock()
	delete(in2fhWeightsPreviousEligibilityTracesByGameID, gameID)
	in2fhWeightsMu.Unlock()

	fh2outWeightsMu.Lock()
	delete(fh2outWeightsPreviousEligibilityTracesByGameID, gameID)
	fh2outWeightsMu.Unlock()
}

func sigmoid(x float32) float32 { return float32(1.0 / (1.0 + math.Exp(float64(-x)))) }

func ValueEstimate(st state.State) (float32, [numFhNodes]float32, [numFhNodes]float32) {
	var (
		estimate       float32
		fhNodePreVals  [numFhNodes]float32
		fhNodePostVals [numFhNodes]float32
	)

	in2fhWeightsMu.RLock()
	my_in2fhWeights := in2fhWeights
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

		fhNodePreVals[fhNodeIdx] = fhNodeSum
		fhNodePostVal := sigmoid(fhNodeSum)
		fhNodePostVals[fhNodeIdx] = fhNodePostVal
		estimate += fhNodePostVal * my_fh2outWeights[fhNodeIdx]
	}
	// Now we artificially add a bias node to the FH layer:
	fhNodeSum := bias
	fhNodePreVals[fhNodeIdx] = fhNodeSum
	fhNodePostVal := sigmoid(fhNodeSum)
	fhNodePostVals[fhNodeIdx] = fhNodePostVal
	estimate += fhNodePostVal * my_fh2outWeights[fhNodeIdx]

	return estimate, fhNodePreVals, fhNodePostVals
}

func weightGradients(st state.State, target float32) (float32, [numFh2OutConnections]float32, [numIn2FhConnections]float32) {
	var (
		in2fhWeightIndex      int // tracks which in2fhWeight we're analyzing.
		fh2outWeightsGradient [numFh2OutConnections]float32
		in2fhWeightsGradient  [numIn2FhConnections]float32
	)
	est, _, fhNodePostVals := ValueEstimate(st)

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

		dEstimate_wrt_fhNodePostVal := my_fh2outWeights[fhNodeIdx]                                       // derive this wrt fhPostVal1: `estimate = (w1*fhPostVal1 + w2*fhPostVal2 + ...)`.
		dFhNodePostVal_wrt_fhNodePreVal := dEstimate_wrt_fh2outWeight * (1 - dEstimate_wrt_fh2outWeight) // derivative of sigmoid(x) is: sigmoid(x) * (1 - sigmoid(x))
		for _, dFhNodePreVal_wrt_in2fhWeight := range st {                                               // derive ths wrt w1: est = (w1*inp1 + w2*inp2 + ...)
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

	return est, fh2outWeightsGradient, in2fhWeightsGradient
}

// TrainWeights back-propagates the error of an estimate against a target.
func TrainWeights(gameID uint32, st state.State, target, learningRate, eligibilityDecayRate float32) float32 {
	est, fh2outWeightsGradient, in2fhWeightsGradient := weightGradients(st, target)
	valueEstimateDiff := target - est // If this diff is positive, we need to add the gradient in the positive direction. else in the negative direction.

	defer in2fhWeightsMu.Unlock()
	in2fhWeightsMu.Lock()
	defer fh2outWeightsMu.Unlock()
	fh2outWeightsMu.Lock()
	// Important: don't write to any of the global vars until these locks are acquired-- that's why very little processing could happen above this line.

	if _, ok := in2fhWeightsPreviousEligibilityTracesByGameID[gameID]; !ok {
		in2fhWeightsPreviousEligibilityTracesByGameID[gameID] = &([numIn2FhConnections]float32{})
	}
	in2fhWeightsPreviousEligibilityTraces := in2fhWeightsPreviousEligibilityTracesByGameID[gameID]
	for i, in2fhWeightDerivative := range in2fhWeightsGradient {
		previousEligibilityTrace := (*in2fhWeightsPreviousEligibilityTraces)[i]
		eligibilityTrace := in2fhWeightDerivative + (eligibilityDecayRate * previousEligibilityTrace)
		(*in2fhWeightsPreviousEligibilityTraces)[i] = eligibilityTrace
		in2fhWeights[i] += learningRate * valueEstimateDiff * eligibilityTrace
	}

	if _, ok := fh2outWeightsPreviousEligibilityTracesByGameID[gameID]; !ok {
		fh2outWeightsPreviousEligibilityTracesByGameID[gameID] = &([numFh2OutConnections]float32{})
	}
	fh2outWeightsPreviousEligibilityTraces := fh2outWeightsPreviousEligibilityTracesByGameID[gameID]
	for i, fh2outWeightDerivative := range fh2outWeightsGradient {
		previousEligibilityTrace := (*fh2outWeightsPreviousEligibilityTraces)[i]
		eligibilityTrace := fh2outWeightDerivative + (eligibilityDecayRate * previousEligibilityTrace)
		(*fh2outWeightsPreviousEligibilityTraces)[i] = eligibilityTrace
		fh2outWeights[i] += learningRate * valueEstimateDiff * eligibilityTrace
	}

	return valueEstimateDiff * valueEstimateDiff // The variance before adjusting the weights.
}
