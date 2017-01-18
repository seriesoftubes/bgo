// Package learn contains utilities for powering an AI agent.
package learn

import (
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/learn/nnet"
	"github.com/seriesoftubes/bgo/learn/nnet/nnperf"
	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

type Agent struct {
	// Epsilon = probability of choosing a random action (at least at first until annealing kicks in)
	epsilon                         float32
	game                            *game.Game
	player                          plyr.Player
	numTrainings                    uint32
	totalVarianceAcrossAllTrainings float32
	statsWG                         sync.WaitGroup
}

func NewAgent(epsilon float32) *Agent {
	return &Agent{epsilon: epsilon}
}

func bestTurnOnePly(b *game.Board, bvt map[turn.TurnArray]turn.Turn, p plyr.Player) turn.Turn {
	var bestTurn turn.Turn
	enemy := p.Enemy()
	worstValForEnemy := float32(3e38)
	for _, t := range bvt {
		bcop := b.Copy()
		bcop.MustExecuteTurn(t, false)
		if val, _ := nnet.ValueEstimate(state.DetectState(enemy, bcop)); val < worstValForEnemy {
			worstValForEnemy = val
			bestTurn = t
		}
	}
	return bestTurn
}

func (a *Agent) WaitForStats() { a.statsWG.Wait() }

func (a *Agent) SetPlayer(p plyr.Player) {
	a.player = p
}

func (a *Agent) TransmitStats() {
	a.statsWG.Add(1)

	go func(tv float32, nt uint32) {
		nnperf.AppendGameData(tv, nt)
		a.statsWG.Done()
	}(a.totalVarianceAcrossAllTrainings, a.numTrainings)

	a.numTrainings = 0
	a.totalVarianceAcrossAllTrainings = 0.0
}

// It's assumed that this is only called when the Agent's GameController is starting a new game.
func (a *Agent) SetGame(g *game.Game) {
	if a.game != nil {
		go func(gid uint32) { nnet.RemoveUselessGameData(gid) }(a.game.ID)
	}

	a.game = g
}

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

	return bestTurnOnePly(b, validTurnsForState, a.player)
}

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	return state.DetectState(a.player, a.game.Board)
}

func (a *Agent) StopLearning() { a.epsilon = 0 }

func (a *Agent) LearnNonFinalState(previousBoard, currentBoard *game.Board) {
	// `currentBoard` is the state that is about to be played by the enemy.
	// `previousBoard` is the state that the hero made a move on which led to `currentBoard`.
	// so the value of `currentBoard` from the hero's POV == -1*(currentboard_value_from_enemyPOV).
	// a.player is the player who made the transition from previous to current board.
	newStateFromEnemyPOV := state.DetectState(a.player.Enemy(), currentBoard)
	previousStateHeroPOV := state.DetectState(a.player, previousBoard)
	enemyEst, _ := nnet.ValueEstimate(newStateFromEnemyPOV)

	a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(a.game.ID, previousStateHeroPOV, -enemyEst)
	a.numTrainings++
}

func (a *Agent) LearnFinal(preWinningMoveBoard, boardInWonState *game.Board, rewardForNextState game.WinKind) {
	actualReward := float32(rewardForNextState)
	previousStateHeroPOV := state.DetectState(a.player, preWinningMoveBoard)

	a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(a.game.ID, previousStateHeroPOV, actualReward)
	a.numTrainings++

	losingStateEnemyPOV := state.DetectState(a.player.Enemy(), boardInWonState)
	a.totalVarianceAcrossAllTrainings += nnet.TrainWeights(a.game.ID, losingStateEnemyPOV, -actualReward)
	a.numTrainings++
}
