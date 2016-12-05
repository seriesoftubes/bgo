// Package learn contains utilities for powering an AI agent.
package learn

import (
	"math"
	"sync"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
)

const maxChexToConsider uint8 = 7

type SerializedTurnsCache struct {
	sync.Mutex
	cache map[string]game.Turn
}

func NewSerializedTurnsCache() *SerializedTurnsCache {
	return &SerializedTurnsCache{cache: make(map[string]game.Turn, 193000)}
}

func (stc *SerializedTurnsCache) Get(s string) game.Turn {
	defer stc.Unlock()
	stc.Lock()
	return stc.cache[s]
}

func (stc *SerializedTurnsCache) Set(s string, t game.Turn) {
	defer stc.Unlock()
	stc.Lock()
	stc.cache[s] = t
}

type QContainer struct {
	sync.Mutex
	qvals map[StateActionPair]float64
}

func NewQContainer() *QContainer {
	return &QContainer{qvals: make(map[StateActionPair]float64, 12888444)}
}

func (qc *QContainer) GetQ(sa StateActionPair) (float64, bool) {
	defer qc.Unlock()
	qc.Lock()
	q, ok := qc.qvals[sa]
	return q, ok
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
	deserializedActions   *SerializedTurnsCache
	qs                    *QContainer
}

func (a *Agent) SetPlayer(p *game.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)     { a.game = g }
func NewAgent(qvals *QContainer, stc *SerializedTurnsCache, alpha, gamma, epsilon float64) *Agent {
	out := &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon}
	out.deserializedActions = stc
	out.qs = qvals
	return out
}

type boardPointState struct {
	isOwnedByMe bool
	numChex     uint8
}

// State which must be hashable.
type (
	State struct {
		boardPoints               [game.NUM_BOARD_POINTS]boardPointState
		numOnMyBar, numOnEnemyBar uint8
		myRoll                    game.Roll
	}

	StateActionPair struct {
		state  State
		action string // serialized, valid Turn
	}
)

func (a *Agent) updateQValueFromSars(state1 State, action string, state2 State, reward uint8) {
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
}

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(state State) (string, game.Turn) {
	validTurns := a.game.ValidTurns()
	possibleActions := make([]string, 0, len(validTurns))
	for svt, t := range validTurns {
		a.deserializedActions.Set(svt, t)
		possibleActions = append(possibleActions, svt)
	}

	var idx int
	if constants.Rand.Float64() < a.epsilon {
		idx = constants.Rand.Intn(len(possibleActions))
	} else {
		var bestQ float64
		var bestQIndices []int
		for idx, action := range possibleActions {
			if q, ok := a.qs.GetQ(StateActionPair{state, action}); ok && q >= bestQ {
				bestQ = q
				bestQIndices = append(bestQIndices, idx)
			}
		}
		if len(bestQIndices) > 0 {
			idxWithinBestQIndices := constants.Rand.Intn(len(bestQIndices))
			idx = bestQIndices[idxWithinBestQIndices]
		} else {
			idx = constants.Rand.Intn(len(possibleActions))
		}
	}

	action := possibleActions[idx]
	return action, a.deserializedActions.Get(action)
}

// Use after training
func (a *Agent) BestAction() (string, game.Turn) {
	validTurns := a.game.ValidTurns()
	if len(validTurns) == 0 {
		return "", nil
	}

	state := a.DetectState()

	bestQ := math.Inf(-1)
	var bestAction string
	for serializedTurn, turn := range validTurns {
		a.deserializedActions.Set(serializedTurn, turn)

		if q, ok := a.qs.GetQ(StateActionPair{state, serializedTurn}); ok && q > bestQ {
			bestQ, bestAction = q, serializedTurn
		}
	}

	return bestAction, a.deserializedActions.Get(bestAction)
}

func (a *Agent) DetectState() State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}

	out := State{}

	out.myRoll = a.game.CurrentRoll.Sorted()

	if a.player == game.PCC {
		out.numOnMyBar = a.game.Board.BarCC
	} else {
		out.numOnMyBar = a.game.Board.BarC
	}

	out.boardPoints = [game.NUM_BOARD_POINTS]boardPointState{}
	for i, p := range a.game.Board.Points {
		chex := p.NumCheckers
		if chex > maxChexToConsider {
			chex = maxChexToConsider
		}
		out.boardPoints[i] = boardPointState{p.Owner == a.player, chex}
	}

	return out
}

func (a *Agent) StopLearning() { a.epsilon = 0 }
func (a *Agent) Learn(state1 State, action string, state2 State, reward game.WinKind) {
	a.updateQValueFromSars(state1, action, state2, uint8(reward))

	a.numObservations++
	if obs := a.numObservations; obs == 800 && a.epsilon > 0.6 {
		a.epsilon = 0.6
	} else if obs == 2000 && a.epsilon > 0.5 {
		a.epsilon = 0.5
	} else if obs == 5000 && a.epsilon > 0.4 {
		a.epsilon = 0.4
	} else if obs == 15000 && a.epsilon > 0.3 {
		a.epsilon = 0.3
	} else if obs == 50000 && a.epsilon > 0.2 {
		a.epsilon = 0.2
	} else if obs == 150000 && a.epsilon > 0.1 {
		a.epsilon = 0.1
	} else if obs == 550500 && a.epsilon > 0.01 {
		a.epsilon = 0.01
	}
}
