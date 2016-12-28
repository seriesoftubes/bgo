// Package nnet contains a neural network that estimates the game equity given the current state of the game.
// The network itself consists of the input layer, 1 fully-connected ReLU layer.
// TODO: maybe add 2nd hidden ReLU layer.
package nnet

import (
	"sync"

	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

const (
	bias                 = float32(3.0)
	numInputs            = len(state.State{}) + 1
	numIn2FhConnections  = numInputs * numInputs // the "FH" aka firsthidden layer is fully-connected with all inputs.
	numFhNodes           = numInputs + 1         // the "FH" layer has 1 node per input, plus one more bias.
	numFh2OutConnections = numFhNodes            // Each "FH" node connects directly to the 1 output signal.
)

var (
	// Weights that connect IN to FH.
	in2fhWeights [numIn2FhConnections]float32 = (func() [numIn2FhConnections]float32 {
		out := [numIn2FhConnections]float32{}
		for i := 0; i < numIn2FhConnections; i++ {
			out[i] = random.Float32Between(-1, 1)
		}
		return out
	})()
	in2fhWeightsMu sync.RWMutex

	// Weights that connect FH to OUT.
	fh2outWeights [numFh2OutConnections]float32 = (func() [numFh2OutConnections]float32 {
		out := [numFh2OutConnections]float32{}
		for i := 0; i < numFh2OutConnections; i++ {
			out[i] = random.Float32Between(-1, 1)
		}
		return out
	})()
	fh2outWeightsMu sync.RWMutex
)

func ValueEstimate(st state.State) (float32, [numFhNodes]float32, [numFhNodes]float32) {
	var (
		estimate       float32
		fhNodePreVals  [numFhNodes]float32
		fhNodePostVals [numFhNodes]float32
	)

	defer in2fhWeightsMu.RUnlock()
	in2fhWeightsMu.RLock()
	defer fh2outWeightsMu.RUnlock()
	fh2outWeightsMu.RLock()

	var fhNodeIdx, in2fhWeightIndex int
	for ; fhNodeIdx < numFhNodes-1; fhNodeIdx++ {
		// Important note: this loop only iterates over the first N-1 FH nodes-- the Nth one needs to be the new bias node!
		var fhNodeSum float32
		for _, num := range st {
			fhNodeSum += num * in2fhWeights[in2fhWeightIndex]
			in2fhWeightIndex++
		}
		// Here we artificially add a bias to the input "state":
		fhNodeSum += bias * in2fhWeights[in2fhWeightIndex]
		in2fhWeightIndex++

		fhNodePreVals[fhNodeIdx] = fhNodeSum
		if fhNodeSum > 0 { // This `if` switch does the ReLU work of the hidden activation.
			fhNodePostVals[fhNodeIdx] = fhNodeSum
			estimate += fhNodeSum * fh2outWeights[fhNodeIdx]
		}
	}
	// Now we artificially add a bias node to the FH layer:
	fhNodeSum := bias
	fhNodePreVals[fhNodeIdx] = fhNodeSum
	if fhNodeSum > 0 { // This `if` switch does the ReLU work of the hidden activation.
		fhNodePostVals[fhNodeIdx] = fhNodeSum
		estimate += fhNodeSum * fh2outWeights[fhNodeIdx]
	}

	return estimate, fhNodePreVals, fhNodePostVals
}

func capBetween(v, min, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func weightGradients(st state.State, target float32) (float32, [numFh2OutConnections]float32, [numIn2FhConnections]float32) {
	var (
		in2fhWeightIndex      int // tracks which in2fhWeight we're analyzing.
		fh2outWeightsGradient [numFh2OutConnections]float32
		in2fhWeightsGradient  [numIn2FhConnections]float32
	)
	est, fhNodePreVals, fhNodePostVals := ValueEstimate(st)

	// Variance = (t - e)**2 = (t - e)*(t - e) = (t**2 -2et + e**2)
	dVariance_wrt_Estimate := 2*est - 2*target // The ^^above^^ equation, derived for `e`.

	defer in2fhWeightsMu.RUnlock()
	in2fhWeightsMu.RLock()
	defer fh2outWeightsMu.RUnlock()
	fh2outWeightsMu.RLock()

	for fhNodeIdx := 0; fhNodeIdx < numFhNodes; fhNodeIdx++ {
		dEstimate_wrt_fh2outWeight := fhNodePostVals[fhNodeIdx] // derive this wrt weight1: `(fhNodePostVal1*weight1 + fhNodePostVal2*weight2 + ...)`.
		dVariance_wrt_fh2outWeight := dVariance_wrt_Estimate * dEstimate_wrt_fh2outWeight
		fh2outWeightsGradient[fhNodeIdx] = dVariance_wrt_fh2outWeight // capBetween(dVariance_wrt_fh2outWeight, -8, 8)

		if fhNodeIdx == numFhNodes-1 {
			break
		}
		// From here down, we're talking about the in2fh weights that connect to this FH node.
		// Remember, the final FH node is a standalone bias node with no weights coming into it-- that's why we skip those in2fh weights for the final FH node.

		dEstimate_wrt_fhNodePostVal := fh2outWeights[fhNodeIdx] // derive this wrt fhPostVal1: `estimate = (w1*fhPostVal1 + w2*fhPostVal2 + ...)`.
		var dFhNodePostVal_wrt_fhNodePreVal float32             // derivative of ReLU(x) is 1.0 when x > 0, 0 otherwise.
		if preValWasPositive := fhNodePreVals[fhNodeIdx] > 0; preValWasPositive {
			dFhNodePostVal_wrt_fhNodePreVal = 1.0
		}
		for _, dFhNodePreVal_wrt_in2fhWeight := range st { // derive ths wrt w1: est = (w1*inp1 + w2*inp2 + ...)
			dVariance_wrt_in2fhWeight := dVariance_wrt_Estimate * dEstimate_wrt_fhNodePostVal * dFhNodePostVal_wrt_fhNodePreVal * dFhNodePreVal_wrt_in2fhWeight
			in2fhWeightsGradient[in2fhWeightIndex] = dVariance_wrt_in2fhWeight // capBetween(dVariance_wrt_in2fhWeight, -8, 8)
			in2fhWeightIndex++
		}
		// and 1 more for the bias one.
		dFhNodePreVal_wrt_in2fhWeight := bias
		dVariance_wrt_in2fhWeight := dVariance_wrt_Estimate * dEstimate_wrt_fhNodePostVal * dFhNodePostVal_wrt_fhNodePreVal * dFhNodePreVal_wrt_in2fhWeight
		in2fhWeightsGradient[in2fhWeightIndex] = dVariance_wrt_in2fhWeight // capBetween(dVariance_wrt_in2fhWeight, -8, 8)
		in2fhWeightIndex++
	}

	return est, fh2outWeightsGradient, in2fhWeightsGradient
}

// TrainWeights back-propagates the error of an estimate against a target.
func TrainWeights(st state.State, target, learningRate float32) float32 {
	est, fh2outWeightsGradient, in2fhWeightsGradient := weightGradients(st, target)

	defer in2fhWeightsMu.Unlock()
	in2fhWeightsMu.Lock()
	defer fh2outWeightsMu.Unlock()
	fh2outWeightsMu.Lock()

	for i, in2fhWeightDerivative := range in2fhWeightsGradient {
		in2fhWeights[i] -= learningRate * in2fhWeightDerivative
	}
	for i, fh2outWeightDerivative := range fh2outWeightsGradient {
		fh2outWeights[i] -= learningRate * fh2outWeightDerivative
	}

	deviation := est - target
	return deviation * deviation // The variance before adjusting the weights.
}
