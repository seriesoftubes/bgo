// Package learn contains utilities for powering an AI agent.
package learn

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/random"
	"github.com/seriesoftubes/bgo/state"
)

const (
	maxChexToConsider      uint8 = 7 // Because does it really matter if you have 7, 8, or 10 chex on a single point? Best dimensionality reduction ever.
	pamIdxStartPointIndex        = 0
	pamIdxFowardDistance         = 1
	pamIdxNumTimes               = 2
	agnosticIndexOfHeroBar       = constants.NUM_BOARD_POINTS
)

type (
	PlayerAgnosticMove [3]uint8 // 0: the agnostic start idx, 1: forward dist, and 2: numTimes.
	sortablePAT        []PlayerAgnosticMove
	PlayerAgnosticTurn [constants.MAX_MOVES_PER_TURN]PlayerAgnosticMove // up to 4 PAMoves, sorted according to sortablePAT logic.

	StateActionPair struct {
		State  state.State
		Action PlayerAgnosticTurn
	}

	QContainer struct {
		sync.RWMutex
		qvals map[StateActionPair]float64
	}

	QcontainerJsonRow struct {
		S state.StateArray
		A PlayerAgnosticTurn
		Q float64
	}

	Agent struct {
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
)

func NewQContainer() *QContainer {
	return &QContainer{qvals: make(map[StateActionPair]float64, 12888444)}
}

func (qc *QContainer) String() string {
	var out []string
	for sa, q := range qc.qvals {
		if q != 0 {
			out = append(out, fmt.Sprintf("%v: %v", q, sa))
		}
	}
	return strings.Join(out, "\n\n")
}

func (qc *QContainer) Serialize(w io.Writer) error {
	enc := json.NewEncoder(w)
	for sap, qval := range qc.qvals {
		if qval == 0 {
			continue
		}

		row := QcontainerJsonRow{sap.State.AsArray(), sap.Action, qval}
		if err := enc.Encode(row); err != nil {
			return fmt.Errorf("JSON enc.Encode(row) error: %v", err)
		}
	}
	return nil
}

func DeserializeQContainer(r io.Reader) (*QContainer, error) {
	out := NewQContainer()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var row QcontainerJsonRow
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return nil, fmt.Errorf("json.Unmarshal error: %v", err)
		}
		st := state.State{}
		st.InitFromArray(row.S)
		out.qvals[StateActionPair{st, PlayerAgnosticTurn(row.A)}] = row.Q
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err() error: %v", err)
	}

	return out, nil
}

func NewAgent(qvals *QContainer, alpha, gamma, epsilon float64) *Agent {
	return &Agent{alpha: alpha, gamma: gamma, epsilon: epsilon, qs: qvals}
}

func (a *Agent) SetPlayer(p *plyr.Player) { a.player = p }
func (a *Agent) SetGame(g *game.Game)     { a.game = g }

// Choose an action that helps with training
func (a *Agent) EpsilonGreedyAction(st state.State, validTurnsForState map[turn.TurnArray]turn.Turn) (PlayerAgnosticTurn, bool) {
	possibleActions := make([]PlayerAgnosticTurn, 0, len(validTurnsForState))
	for _, t := range validTurnsForState {
		possibleActions = append(possibleActions, AgnosticizeTurn(t, a.player))
	}

	var idx int
	var wasCacheHit bool
	if random.Float64() < a.epsilon {
		idx = random.IntBetween(0, len(possibleActions)-1)
	} else {
		var bestQ float64
		var bestQIndices []int
		defer a.qs.RUnlock()
		a.qs.RLock()
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
		wasCacheHit = bestQ > 0
	}

	return possibleActions[idx], wasCacheHit
}

func (a *Agent) DetectState() state.State {
	if a.game.CurrentPlayer != a.player {
		panic("shouldn't be detecting the state outside of the agent's own turn.")
	}
	s, _ := state.DetectState(a.player, a.game, maxChexToConsider)
	return s
}

func (a *Agent) StopLearning() { a.epsilon = 0 }

func (a *Agent) oldAndBestFutureQ(oldStateAction StateActionPair, state2 state.State, validTurnsInState2 map[turn.TurnArray]turn.Turn) (float64, float64) {
	defer a.qs.RUnlock()
	a.qs.RLock()

	oldQ := a.qs.qvals[oldStateAction]

	var bestPossibleFutureQ float64
	for _, t := range validTurnsInState2 {
		if q, ok := a.qs.qvals[StateActionPair{state2, AgnosticizeTurn(t, a.player)}]; ok && q > bestPossibleFutureQ {
			bestPossibleFutureQ = q
		}
	}

	return oldQ, bestPossibleFutureQ
}

func (a *Agent) Learn(state1 state.State, action PlayerAgnosticTurn, state2 state.State, rewardForState2 game.WinKind, validTurnsInState2 map[turn.TurnArray]turn.Turn) {
	oldStateAction := StateActionPair{state1, action}
	oldQ, bestPossibleFutureQ := a.oldAndBestFutureQ(oldStateAction, state2, validTurnsInState2)

	defer a.qs.Unlock()
	a.qs.Lock()
	if newQ := oldQ + a.alpha*(float64(rewardForState2)+(a.gamma*bestPossibleFutureQ)-oldQ); newQ != 0 {
		a.qs.qvals[oldStateAction] = newQ
	} else {
		delete(a.qs.qvals, oldStateAction)
	}

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

func (pam PlayerAgnosticMove) isEmpty() bool { return pam[pamIdxNumTimes] == 0 }
func (pam PlayerAgnosticMove) asMove(p *plyr.Player) turn.Move {
	var letter string
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

func AgnosticizeTurn(t turn.Turn, p *plyr.Player) PlayerAgnosticTurn {
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
		panic("turn was invalid (had > 4 moves) for player " + *p)
	}
	for i, pam := range spat {
		out[i] = pam
	}
	return out
}

func ConvertAgnosticTurn(paa PlayerAgnosticTurn, p *plyr.Player) turn.Turn {
	out := turn.Turn{}
	for _, pam := range paa {
		if pam.isEmpty() {
			break
		}
		out[pam.asMove(p)] = pam[pamIdxNumTimes]
	}
	return out
}
