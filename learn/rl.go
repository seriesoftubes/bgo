// Package learn contains utilities for powering an AI agent.
package learn

import (
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/learn/nnet"
	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

type Agent struct {
	// Alpha = learning rate
	// Gamma = discount rate for future rewards
	// Epsilon = probability of choosing a random action (at least at first until annealing kicks in)
	// TODO: annealing rate?
	alpha, gamma, epsilon           float32
	game                            *game.Game
	player                          plyr.Player
	numTrainings                    uint32
	totalVarianceAcrossAllTrainings float32
}

func NewAgent(alpha, gamma, epsilon float32) *Agent {
	return &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon}
}

func (a *Agent) ResetNeuralNetworkStats() { a.numTrainings, a.totalVarianceAcrossAllTrainings = 0, 0.0 }
func (a *Agent) AverageNeuralNetworkVariance() float32 {
	return a.totalVarianceAcrossAllTrainings / float32(a.numTrainings)
}
func (a *Agent) SetPlayer(p plyr.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)    { a.game = g }

func (a *Agent) EpsilonGreedyAction(b *game.Board, validTurnsForState map[turn.TurnArray]turn.Turn) turn.Turn {
	if len(validTurnsForState) == 0 {
		panic("should have prevented this function from being called!")
	}

	if random.Float32Between(0, 1) < a.epsilon { // exploration mode.
		for _, t := range validTurnsForState {
			return t
		}
		panic("should have prevented this function from being called!")
	}

	// Use 1-ply lookahead to get the best turn. TODO: use 3-ply.
	var bestTurn turn.Turn
	worstValForEnemy := float32(3e38)
	for _, t := range validTurnsForState {
		bcop := b.Copy()
		bcop.MustExecuteTurn(t, false)
		if val, _, _ := nnet.ValueEstimate(state.DetectState(a.player.Enemy(), bcop)); val < worstValForEnemy {
			worstValForEnemy = val
			bestTurn = t
		}
	}
	return bestTurn
}

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	return state.DetectState(a.player, a.game.Board)
}

func (a *Agent) StopLearning() { a.epsilon = 0 }

func (a *Agent) LearnNonFinalState(previousBoard, currentBoard *game.Board) {
	est, _, _ := nnet.ValueEstimate(state.DetectState(a.player, currentBoard))
	a.numTrainings++
	a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(state.DetectState(a.player, previousBoard), est*a.gamma, a.alpha)
}

func (a *Agent) LearnFinal(previousHeroBoard, previousEnemyBoard *game.Board, rewardForNextState game.WinKind) {
	discountedEstimate := float32(rewardForNextState) * a.gamma

	a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(state.DetectState(a.player, previousHeroBoard), discountedEstimate, a.alpha)
	a.numTrainings++

	if previousEnemyBoard != nil { // If the enemy is human, their previous boards aren't saved, which is fine.
		a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(state.DetectState(a.player.Enemy(), previousEnemyBoard), -discountedEstimate, a.alpha)
		a.numTrainings++
	}
}
