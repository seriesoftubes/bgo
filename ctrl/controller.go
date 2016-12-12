// Package ctrl contains a game controller, which moves a game along through its turns
package ctrl

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/game/turngen"
	"github.com/seriesoftubes/bgo/learn"
	"github.com/seriesoftubes/bgo/render"
	"github.com/seriesoftubes/bgo/state"
)

func readTurnFromStdin(validTurns map[turn.TurnArray]turn.Turn) turn.Turn {
	for {
		var supposedlySerializedTurn string
		fmt.Scanln(&supposedlySerializedTurn)

		t, err := turn.DeserializeTurn(supposedlySerializedTurn)
		if err != nil {
			fmt.Println("could not read your instructions, please try again: " + err.Error())
			continue
		}

		if _, ok := validTurns[t.Arrayify()]; ok {
			return t
		} else {
			fmt.Println("invalid turn entered, please try again")
		}
	}
}

// The best AI ever built.
func randomlyChooseValidTurn(validTurns map[turn.TurnArray]turn.Turn) turn.Turn {
	for _, t := range validTurns {
		return t
	}
	panic("no turns to choose. you should've prevented this line from being reached")
}

type GameController struct {
	g                *game.Game
	debug            bool
	agent            *learn.Agent
	prevStateActions map[*plyr.Player]*stateAction
}

func New(qs *learn.QContainer, debug bool) *GameController {
	learningRateAkaAlpha := 0.25
	rewardsDiscountRateAkaGamma := 0.5
	initialExplorationRateAkaEpsilon := 1.0
	prevStateActions := make(map[*plyr.Player]*stateAction, 2)
	agent := learn.NewAgent(qs, learningRateAkaAlpha, rewardsDiscountRateAkaGamma, initialExplorationRateAkaEpsilon)
	return &GameController{agent: agent, prevStateActions: prevStateActions, debug: debug}
}

func (gc *GameController) PlayOneGame(numHumanPlayers uint8, stopLearning bool) (*plyr.Player, game.WinKind) {
	gc.g = game.NewGame(numHumanPlayers)

	if stopLearning {
		gc.agent.StopLearning()
	}
	gc.agent.SetGame(gc.g)

	gc.maybePrint("Welcome to backgammon. Good luck and have fun!")
	var done bool
	for !done {
		done = gc.playOneTurn()
	}

	if gc.g.HasAnyHumans() || gc.debug {
		render.PrintGame(gc.g)
		fmt.Println("\tDONE WITH GAME!")
	}

	return gc.g.Board.Winner(), gc.g.Board.WinKind()
}

func (gc *GameController) maybePrint(s ...interface{}) {
	if gc.g.HasAnyHumans() || gc.debug {
		fmt.Println(s...)
	}
}

type stateAction struct {
	state  state.State
	action learn.PlayerAgnosticTurn
}

func (gc *GameController) chooseTurn(validTurns map[turn.TurnArray]turn.Turn, currentState state.State) turn.Turn {
	if gc.g.IsCurrentPlayerHuman() {
		return readTurnFromStdin(validTurns)
	}

	pat := gc.agent.EpsilonGreedyAction(currentState, validTurns)
	gc.prevStateActions[gc.g.CurrentPlayer] = &stateAction{currentState, pat}

	return learn.ConvertAgnosticTurn(pat, gc.g.CurrentPlayer)
}

// playOneTurn plays through one turn, and returns whether the game is finished after the turn executes.
func (gc *GameController) playOneTurn() bool {
	g := gc.g

	if g.HasAnyHumans() || gc.debug {
		render.PrintGame(g)
	}

	validTurns := turngen.ValidTurns(g.Board, g.CurrentRoll, g.CurrentPlayer)

	isComputer := !gc.g.IsCurrentPlayerHuman()
	var currentState state.State
	if isComputer {
		gc.agent.SetPlayer(g.CurrentPlayer)
		currentState = gc.agent.DetectState()
		if prevSA, ok := gc.prevStateActions[g.CurrentPlayer]; ok {
			gc.agent.Learn(prevSA.state, prevSA.action, currentState, game.WinKindNotWon, validTurns)
		}
	}

	var chosenTurn turn.Turn
	if len(validTurns) == 0 {
		gc.maybePrint("\tcan't do anything this turn, sorry!")
	} else if len(validTurns) == 1 {
		gc.maybePrint("\tthis turn only has 1 option, forcing!")
		chosenTurn = randomlyChooseValidTurn(validTurns)
		if isComputer {
			gc.prevStateActions[g.CurrentPlayer] = &stateAction{currentState, learn.AgnosticizeTurn(chosenTurn, g.CurrentPlayer)}
		}
	} else {
		gc.maybePrint(fmt.Sprintf("\tYour move, %q:", *g.CurrentPlayer))
		chosenTurn = gc.chooseTurn(validTurns, currentState)
	}
	gc.maybePrint("\tChose move:", chosenTurn)

	if gc.debug {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			fmt.Println("\n\n\n\t\tRecovering in playOneTurn() from", r)
			fmt.Println("\tWTF move:", chosenTurn)
			render.PrintGame(g)
			panic("ok, panic again because this should kill the whole game")
		}()
	}

	gc.g.Board.MustExecuteTurn(chosenTurn, gc.debug)
	winner, winAmt := gc.g.Board.Winner(), gc.g.Board.WinKind()

	if winner != nil {
		if isComputer {
			// special case: a computer won, and they need to learn from that without having to re-run this playOneTurn method.
			prevSA := gc.prevStateActions[g.CurrentPlayer]
			gc.agent.Learn(prevSA.state, prevSA.action, gc.agent.DetectState(), winAmt, nil) // There are 0 valid turns after a game has been won.
		}

		return true
	} else {
		gc.g.NextPlayersTurn()
		return false
	}
}
