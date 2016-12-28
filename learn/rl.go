// Package learn contains utilities for powering an AI agent.
package learn

import (
	"sort"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/state"
)

const (
	pamIdxStartPointIndex  = 0
	pamIdxFowardDistance   = 1
	pamIdxNumTimes         = 2
	agnosticIndexOfHeroBar = constants.NUM_BOARD_POINTS
)

type (
	PlayerAgnosticMove [3]uint8 // 0: the agnostic start idx, 1: forward dist, and 2: numTimes.
	sortablePAT        []PlayerAgnosticMove
	PlayerAgnosticTurn [constants.MAX_MOVES_PER_TURN]PlayerAgnosticMove // up to 4 PAMoves, sorted according to sortablePAT logic.

	Agent struct {
		// Alpha = learning rate
		// Gamma = discount rate for future rewards
		// Epsilon = probability of choosing a random action (at least at first until annealing kicks in)
		// TODO: annealing rate?
		alpha, gamma, epsilon float64
		game                  *game.Game
		player                plyr.Player
	}
)

func NewAgent(alpha, gamma, epsilon float64) *Agent {
	return &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon}
}

func (a *Agent) SetPlayer(p plyr.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)    { a.game = g }

// TODO: Agent interface with DetectState and EpsilonGreedyAction?
func (a *Agent) EpsilonGreedyAction(st state.State, validTurnsForState map[turn.TurnArray]turn.Turn) PlayerAgnosticTurn {
	for _, t := range validTurnsForState {
		return AgnosticizeTurn(t, a.player)
	}
	panic("should have prevented this function from being called!")
}

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	return state.DetectState(a.player, a.game)
}

func (a *Agent) StopLearning() { a.epsilon = 0 }

func (a *Agent) Learn(state1 state.State, state2 state.State, rewardForState2 game.WinKind) {}

func (pam PlayerAgnosticMove) isEmpty() bool { return pam[pamIdxNumTimes] == 0 }
func (pam PlayerAgnosticMove) asMove(p plyr.Player) turn.Move {
	var letter byte
	if bpi := pam[pamIdxStartPointIndex]; p == plyr.PCC {
		if bpi == agnosticIndexOfHeroBar {
			letter = constants.LETTER_BAR_CC
		} else {
			letter = constants.Num2Alpha[bpi]
		}
	} else {
		if bpi == agnosticIndexOfHeroBar {
			letter = constants.LETTER_BAR_C
		} else {
			letter = constants.Num2Alpha[constants.FINAL_BOARD_POINT_INDEX-bpi]
		}
	}

	return turn.Move{Requestor: p, FowardDistance: pam[pamIdxFowardDistance], Letter: letter}
}

func (sp sortablePAT) Len() int      { return len(sp) }
func (sp sortablePAT) Swap(i, j int) { sp[j], sp[i] = sp[i], sp[j] }
func (sp sortablePAT) Less(i, j int) bool {
	left, right := sp[i], sp[j]

	if leftIdx, rightIdx := left[pamIdxStartPointIndex], right[pamIdxStartPointIndex]; leftIdx < rightIdx {
		return true
	} else if leftIdx > rightIdx {
		return false
	}

	if leftDist, rightDist := left[pamIdxFowardDistance], right[pamIdxFowardDistance]; leftDist < rightDist {
		return true
	} else if leftDist > rightDist {
		return false
	}

	return left[pamIdxNumTimes] < right[pamIdxNumTimes]
}

func AgnosticizeTurn(t turn.Turn, p plyr.Player) PlayerAgnosticTurn {
	var spat sortablePAT

	if p == plyr.PC {
		for mov, times := range t {
			var idx uint8
			if mov.Letter == constants.LETTER_BAR_C {
				idx = agnosticIndexOfHeroBar
			} else if mov.Letter == constants.LETTER_BAR_CC {
				panic("can't agnosticize an invalid turn (included moving something off the enemy's bar")
			} else {
				idx = constants.FINAL_BOARD_POINT_INDEX - constants.Alpha2Num[mov.Letter]
			}
			spat = append(spat, PlayerAgnosticMove{idx, mov.FowardDistance, times})
		}
	} else {
		for mov, times := range t {
			var idx uint8
			if mov.Letter == constants.LETTER_BAR_CC {
				idx = agnosticIndexOfHeroBar
			} else if mov.Letter == constants.LETTER_BAR_C {
				panic("can't agnosticize an invalid turn (included moving something off the enemy's bar")
			} else {
				idx = constants.Alpha2Num[mov.Letter]
			}
			spat = append(spat, PlayerAgnosticMove{idx, mov.FowardDistance, times})
		}
	}

	sort.Sort(spat)
	out := PlayerAgnosticTurn{}
	if len(spat) > constants.MAX_MOVES_PER_TURN {
		panic("turn was invalid (had > 4 moves) for player " + string(p))
	}
	for i, pam := range spat {
		out[i] = pam
	}
	return out
}

func ConvertAgnosticTurn(paa PlayerAgnosticTurn, p plyr.Player) turn.Turn {
	out := turn.Turn{}
	for _, pam := range paa {
		if pam.isEmpty() {
			break
		}
		out[pam.asMove(p)] = pam[pamIdxNumTimes]
	}
	return out
}
