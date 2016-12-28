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

const (
	msgWelcome      = "Welcome to backgammon. Good luck and have fun!"
	msgGameOver     = "\tDONE WITH GAME!"
	msgNoMovesAvail = "\tcan't do anything this turn, sorry!"
	msgForceMove    = "\tthis turn only has 1 option, forcing!"
	msgAskForMove   = "\tYour move, "
	msgChoseMove    = "\tChose move:"
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
	g          *game.Game
	debug      bool
	agent      *learn.Agent
	prevStates map[plyr.Player]state.State
}

func New(debug bool) *GameController {
	learningRateAkaAlpha := 0.5
	rewardsDiscountRateAkaGamma := 0.9999999999
	initialExplorationRateAkaEpsilon := 1.0
	prevStates := make(map[plyr.Player]state.State, 2)
	agent := learn.NewAgent(learningRateAkaAlpha, rewardsDiscountRateAkaGamma, initialExplorationRateAkaEpsilon)
	return &GameController{agent: agent, prevStates: prevStates, debug: debug}
}

func (gc *GameController) PlayOneGame(numHumanPlayers uint8, stopLearning bool) (plyr.Player, game.WinKind) {
	gc.g = game.NewGame(numHumanPlayers)

	if stopLearning {
		gc.agent.StopLearning()
	}
	gc.agent.SetGame(gc.g)

	gc.maybePrint(msgWelcome)
	var done bool
	for !done {
		done = gc.playOneTurn()
	}

	if gc.g.HasAnyHumans() || gc.debug {
		render.PrintGame(gc.g)
		fmt.Println(msgGameOver)
	}

	return gc.g.Board.Winner(), gc.g.Board.WinKind()
}

func (gc *GameController) maybePrint(s ...interface{}) {
	if gc.g.HasAnyHumans() || gc.debug {
		fmt.Println(s...)
	}
}

func (gc *GameController) chooseTurn(validTurns map[turn.TurnArray]turn.Turn, currentState state.State) turn.Turn {
	if gc.g.IsCurrentPlayerHuman() {
		return readTurnFromStdin(validTurns)
	}

	pat := gc.agent.EpsilonGreedyAction(currentState, validTurns)
	gc.prevStates[gc.g.CurrentPlayer] = currentState

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
		if prevState, ok := gc.prevStates[g.CurrentPlayer]; ok {
			gc.agent.Learn(prevState, currentState, game.WinKindNotWon)
		}
	}

	var chosenTurn turn.Turn
	if len(validTurns) == 0 {
		gc.maybePrint(msgNoMovesAvail)
	} else if len(validTurns) == 1 {
		gc.maybePrint(msgForceMove)
		chosenTurn = randomlyChooseValidTurn(validTurns)
		if isComputer {
			gc.prevStates[g.CurrentPlayer] = currentState
		}
	} else {
		gc.maybePrint(msgAskForMove, string(g.CurrentPlayer))
		chosenTurn = gc.chooseTurn(validTurns, currentState)
	}
	gc.maybePrint(msgChoseMove, chosenTurn)

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

	if winner != 0 {
		if isComputer {
			// special case: a computer won, and they need to learn from that without having to re-run this playOneTurn method.
			prevState := gc.prevStates[g.CurrentPlayer]
			gc.agent.Learn(prevState, gc.agent.DetectState(), winAmt)
		}

		return true
	} else {
		gc.g.NextPlayersTurn()
		return false
	}
}
