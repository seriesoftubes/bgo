// Package learn contains utilities for powering an AI agent.
package learn

import (
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/game/turngen"
	"github.com/seriesoftubes/bgo/learn/nnet"
	"github.com/seriesoftubes/bgo/learn/nnet/nnperf"
	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

const (
	numOneOfAKindRolls  = 6
	numRollPermutations = float32(32)
)

var uniqueRolls [21]game.Roll = [21]game.Roll{
	// one-of-a-kind rolls
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 4},
	{5, 5},
	{6, 6},
	// Below here, there are 2 of each kind of roll.
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 5},
	{1, 6},
	{2, 3},
	{2, 4},
	{2, 5},
	{2, 6},
	{3, 4},
	{3, 5},
	{3, 6},
	{4, 5},
	{4, 6},
	{5, 6},
}

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

func bestTurnOnePly(b *game.Board, bvt map[turn.TurnArray]turn.Turn, p plyr.Player) [3]*turn.Turn {
	var topVals [3]*float32 // the top-most negative values.
	var topTurns [3]*turn.Turn
	updateTop := func(val float32, t turn.Turn) {
		if valPtr0 := topVals[0]; valPtr0 == nil || val > *valPtr0 {
			topVals[0] = &val
			topTurns[0] = &t
			return
		} // From here down, val <= vals0

		if valPtr1 := topVals[1]; valPtr1 == nil || val > *valPtr1 {
			topVals[1] = &val
			topTurns[1] = &t
			return
		} // From here down, val <= vals1

		if valPtr2 := topVals[2]; valPtr2 == nil || val > *valPtr2 {
			topVals[2] = &val
			topTurns[2] = &t
		}
	}

	enemy := p.Enemy()
	for _, t := range bvt {
		bcop := b.Copy()
		bcop.MustExecuteTurn(t, false)
		val, _ := nnet.ValueEstimate(state.DetectState(enemy, bcop))
		updateTop(-1*val, t) // the top consists of the most negative values so multiply val by -1.
	}

	return topTurns
}

func updateRollAVG(rollIdx int, total *float32, newNum float32) {
	if rollIdx < numOneOfAKindRolls { // the first 6 rolls are 1 of a kind. otherwise there's 2 instances.
		*total += newNum
	} else {
		*total += newNum * 2
	}
}

func bestTurnTwoPly(b *game.Board, heroTurns0 map[turn.TurnArray]turn.Turn, hero plyr.Player) turn.Turn {
	enemy := hero.Enemy()
	bestVal := float32(-9e37)
	var bestTurn turn.Turn

	for _, heroTurn0 := range heroTurns0 { // these turns yield the game over to the enemy.
		enemyBoard1 := b.Copy()
		enemyBoard1.MustExecuteTurn(heroTurn0, false)

		var totalRollEquity1 float32
		for rollIdx1, r1 := range uniqueRolls {
			var totalTurnVal1, numTurns1 float32
			for _, enemyTurnPtr1 := range bestTurnOnePly(enemyBoard1, turngen.ValidTurns(enemyBoard1, r1, enemy), enemy) {
				if enemyTurnPtr1 == nil {
					break
				}
				heroBoard2 := enemyBoard1.Copy()
				heroBoard2.MustExecuteTurn(*enemyTurnPtr1, false)
				val, _ := nnet.ValueEstimate(state.DetectState(hero, heroBoard2))
				totalTurnVal1 += val
				numTurns1++
			}

			var avgTurnVal1 float32
			if numTurns1 > 0.0 {
				avgTurnVal1 = totalTurnVal1 / numTurns1
			} else {
				// there were no turns for the enemy, so go straight to hero turn
				val, _ := nnet.ValueEstimate(state.DetectState(hero, enemyBoard1))
				avgTurnVal1 = val
			}
			if rollIdx1 < numOneOfAKindRolls {
				totalRollEquity1 += avgTurnVal1
			} else {
				totalRollEquity1 += 2 * avgTurnVal1
			}
		}

		if totalRollEquity1 > bestVal {
			bestTurn = heroTurn0
			bestVal = totalRollEquity1
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

	return *(bestTurnOnePly(b, validTurnsForState, a.player)[0])
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
