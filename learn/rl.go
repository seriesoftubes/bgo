// Package learn contains utilities for powering an AI agent.
package learn

import (
	"math"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
)

const maxChexToConsider uint8 = 9

type Agent struct {
	// Alpha = learning rate
	// Gamma = discount rate for future rewards
	// Epsilon = probability of choosing a random action (at least at first until annealing kicks in)
	// TODO: annealing rate?
	alpha, gamma, epsilon float64
	game                  *game.Game
	player                *game.Player
	numObservations       uint64
	deserializedActions   map[string]game.Turn
	qByStateAction        map[StateActionPair]float64
}

func (a *Agent) SetPlayer(p *game.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)     { a.game = g }
func NewAgent(alpha, gamma, epsilon float64) *Agent {
	out := &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon}
	out.deserializedActions = make(map[string]game.Turn, 193000)
	out.qByStateAction = make(map[StateActionPair]float64, 12888444)
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
	var oldQ float64 // By default, assume Q is zero.
	sa := StateActionPair{state1, action}
	if q, ok := a.qByStateAction[sa]; ok {
		oldQ = q
	}

	var bestPossibleFutureQ float64
	for serializedTurn := range a.game.ValidTurns() {
		if q, ok := a.qByStateAction[StateActionPair{state1, serializedTurn}]; ok && q > bestPossibleFutureQ {
			bestPossibleFutureQ = q
		}
	}

	a.qByStateAction[sa] = oldQ + a.alpha*(float64(reward)+(a.gamma*bestPossibleFutureQ)-oldQ)
}

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(state State) (string, game.Turn) {
	validTurns := a.game.ValidTurns()
	possibleActions := make([]string, 0, len(validTurns))
	for svt, t := range validTurns {
		a.deserializedActions[svt] = t
		possibleActions = append(possibleActions, svt)
	}

	var idx int
	if constants.Rand.Float64() < a.epsilon {
		idx = constants.Rand.Intn(len(possibleActions))
	} else {
		var bestQ float64
		var bestQIndices []int
		for idx, action := range possibleActions {
			if q, ok := a.qByStateAction[StateActionPair{state, action}]; ok && q >= bestQ {
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
	return action, a.deserializedActions[action]
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
		a.deserializedActions[serializedTurn] = turn // memoize

		if q, ok := a.qByStateAction[StateActionPair{state, serializedTurn}]; ok && q > bestQ {
			bestQ, bestAction = q, serializedTurn
		}
	}

	return bestAction, a.deserializedActions[bestAction]
}

func (a *Agent) DetectState() State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}

	out := State{}

	out.myRoll = *a.game.CurrentRoll

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
	if a.numObservations == 800 && a.epsilon > 0.6 {
		a.epsilon = 0.6
	} else if a.numObservations == 2000 && a.epsilon > 0.5 {
		a.epsilon = 0.5
	} else if a.numObservations == 5000 && a.epsilon > 0.4 {
		a.epsilon = 0.4
	} else if a.numObservations == 15000 && a.epsilon > 0.3 {
		a.epsilon = 0.3
	} else if a.numObservations == 50000 && a.epsilon > 0.2 {
		a.epsilon = 0.2
	} else if a.numObservations == 150000 && a.epsilon > 0.1 {
		a.epsilon = 0.1
	} else if a.numObservations == 550500 && a.epsilon > 0.01 {
		a.epsilon = 0.01
	}
}
