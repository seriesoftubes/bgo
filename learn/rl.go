// Package learn contains utilities for powering an AI agent.
package learn

import (
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
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
	player                *plyr.Player
	numObservations       uint64
	qs                    *QContainer
}

func (a *Agent) SetPlayer(p *plyr.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)     { a.game = g }
func NewAgent(qvals *QContainer, alpha, gamma, epsilon float64) *Agent {
	return &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon, qs: qvals}
}

type StateActionPair struct {
	state  state.State
	action string // serialized, valid Turn
}

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(st state.State, validTurnsForState map[string]turn.Turn) string {
	possibleActions := make([]string, 0, len(validTurnsForState))
	for svt := range validTurnsForState {
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

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	s, _ := state.DetectState(a.player, a.game, maxChexToConsider)
	return s
}

func (a *Agent) StopLearning() { a.epsilon = 0 }
func (a *Agent) Learn(state1 state.State, action string, state2 state.State, rewardForState2 game.WinKind, validTurnsInState2 map[string]turn.Turn) {
	defer a.qs.Unlock()
	a.qs.Lock()

	oldStateAction := StateActionPair{state1, action}
	oldQ := a.qs.qvals[oldStateAction]

	var bestPossibleFutureQ float64
	for serializedTurn := range validTurnsInState2 {
		if q, ok := a.qs.qvals[StateActionPair{state2, serializedTurn}]; ok && q > bestPossibleFutureQ {
			bestPossibleFutureQ = q
		}
	}
	a.qs.qvals[oldStateAction] = oldQ + a.alpha*(float64(rewardForState2)+(a.gamma*bestPossibleFutureQ)-oldQ)

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
