// Package learn contains utilities for powering an AI agent.
package learn

import (
	"math"
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

const maxChexToConsider uint8 = 7

type QContainer struct {
	sync.Mutex
	qvals map[StateActionPair]float64
}

func NewQContainer() *QContainer {
	return &QContainer{qvals: make(map[StateActionPair]float64, 12888444)}
}

type Agent struct {
	// Alpha = learning rate
	// Gamma = discount rate for future rewards
	// Epsilon = probability of choosing a random action (at least at first until annealing kicks in)
	// TODO: annealing rate?
	alpha, gamma, epsilon float64
	game                  *game.Game
	player                *game.Player
	numObservations       uint64
	qs                    *QContainer
}

func (a *Agent) SetPlayer(p *game.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)     { a.game = g }
func NewAgent(qvals *QContainer, alpha, gamma, epsilon float64) *Agent {
	return &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon, qs: qvals}
}

type StateActionPair struct {
	state  state.State
	action string // serialized, valid Turn
}

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(st state.State) string {
	validTurns := a.game.ValidTurns()
	possibleActions := make([]string, 0, len(validTurns))
	for svt := range validTurns {
		possibleActions = append(possibleActions, svt)
	}

	var idx int
	if random.Float64() < a.epsilon {
		idx = random.IntBetween(0, len(possibleActions)-1)
	} else {
		var bestQ float64
		var bestQIndices []int
		defer a.qs.Unlock()
		a.qs.Lock()
		for idx, action := range possibleActions {
			if q, ok := a.qs.qvals[StateActionPair{st, action}]; ok && q >= bestQ {
				bestQ = q
				bestQIndices = append(bestQIndices, idx)
			}
		}
		if len(bestQIndices) > 0 {
			idxWithinBestQIndices := random.IntUpTo(len(bestQIndices))
			idx = bestQIndices[idxWithinBestQIndices]
		} else {
			idx = random.IntUpTo(len(possibleActions))
		}
	}

	return possibleActions[idx]
}

// Use after training
func (a *Agent) BestAction() string {
	validTurns := a.game.ValidTurns()
	if len(validTurns) == 0 {
		return ""
	}

	state := a.DetectState()

	bestQ := math.Inf(-1)
	var bestAction string
	defer a.qs.Unlock()
	a.qs.Lock()
	for serializedTurn := range validTurns {
		if q, ok := a.qs.qvals[StateActionPair{state, serializedTurn}]; ok && q > bestQ {
			bestQ, bestAction = q, serializedTurn
		}
	}

	return bestAction
}

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	s, _ := state.DetectState(a.player, a.game, maxChexToConsider)
	return s
}

func (a *Agent) StopLearning() { a.epsilon = 0 }
func (a *Agent) Learn(state1 state.State, action string, state2 state.State, reward game.WinKind) {
	defer a.qs.Unlock()
	a.qs.Lock()
	var oldQ float64 // By default, assume Q is zero.
	sa := StateActionPair{state1, action}
	if q, ok := a.qs.qvals[sa]; ok {
		oldQ = q
	}
	var bestPossibleFutureQ float64
	for serializedTurn := range a.game.ValidTurns() {
		if q, ok := a.qs.qvals[StateActionPair{state1, serializedTurn}]; ok && q > bestPossibleFutureQ {
			bestPossibleFutureQ = q
		}
	}
	a.qs.qvals[sa] = oldQ + a.alpha*(float64(reward)+(a.gamma*bestPossibleFutureQ)-oldQ)

	a.numObservations++
	if obs := a.numObservations; obs == 80100 && a.epsilon > 0.6 {
		a.epsilon = 0.6
	} else if obs == 200100 && a.epsilon > 0.5 {
		a.epsilon = 0.5
	} else if obs == 500100 && a.epsilon > 0.4 {
		a.epsilon = 0.4
	} else if obs == 1200100 && a.epsilon > 0.3 {
		a.epsilon = 0.3
	} else if obs == 5200100 && a.epsilon > 0.2 {
		a.epsilon = 0.2
	} else if obs == 15200100 && a.epsilon > 0.1 {
		a.epsilon = 0.1
	} else if obs == 55200100 && a.epsilon > 0.01 {
		a.epsilon = 0.01
	}
}
