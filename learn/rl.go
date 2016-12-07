// Package learn contains utilities for powering an AI agent.
package learn

import (
	"math"
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/random"
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

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(state State) string {
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
			if q, ok := a.qs.qvals[StateActionPair{state, action}]; ok && q >= bestQ {
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

func uint8Ceiling(x, max uint8) uint8 {
	if x > max {
		return max
	}
	return x
}

func (a *Agent) DetectState() State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}

	out := State{}

	out.myRoll = a.game.CurrentRoll.Sorted()

	isPCC := a.player == game.PCC
	if isPCC {
		out.numOnMyBar = uint8Ceiling(a.game.Board.BarCC, maxChexToConsider)
		out.numOnEnemyBar = uint8Ceiling(a.game.Board.BarC, maxChexToConsider)
	} else {
		out.numOnMyBar = uint8Ceiling(a.game.Board.BarC, maxChexToConsider)
		out.numOnEnemyBar = uint8Ceiling(a.game.Board.BarCC, maxChexToConsider)
	}

	out.boardPoints = [game.NUM_BOARD_POINTS]boardPointState{}
	lastPointIndex := int(game.NUM_BOARD_POINTS - 1)
	for ptIdx, pt := range a.game.Board.Points {
		chex := uint8Ceiling(pt.NumCheckers, maxChexToConsider)
		// fill them in order of distance from enemy home. so PCC starts as normal
		translatedPtIdx := lastPointIndex - ptIdx
		if isPCC {
			translatedPtIdx = ptIdx
		}
		out.boardPoints[translatedPtIdx] = boardPointState{pt.Owner == a.player, chex}
	}

	return out
}

func (a *Agent) StopLearning() { a.epsilon = 0 }
func (a *Agent) Learn(state1 State, action string, state2 State, reward game.WinKind) {
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
	if obs := a.numObservations; obs == 80000 && a.epsilon > 0.6 {
		a.epsilon = 0.6
	} else if obs == 200000 && a.epsilon > 0.5 {
		a.epsilon = 0.5
	} else if obs == 500000 && a.epsilon > 0.4 {
		a.epsilon = 0.4
	} else if obs == 1500800 && a.epsilon > 0.3 {
		a.epsilon = 0.3
	} else if obs == 5200800 && a.epsilon > 0.2 {
		a.epsilon = 0.2
	} else if obs == 15300500 && a.epsilon > 0.1 {
		a.epsilon = 0.1
	} else if obs == 50300800 && a.epsilon > 0.01 {
		a.epsilon = 0.01
	}
}
