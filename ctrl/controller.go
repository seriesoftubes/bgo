// Package ctrl contains a game controller, which moves a game along through its turns
package ctrl

import (
	"fmt"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/learn"
	"github.com/seriesoftubes/bgo/render"
)

func readTurnFromStdin(validTurns map[string]game.Turn) game.Turn {
	for {
		var supposedlySerializedTurn string
		fmt.Scanln(&supposedlySerializedTurn)

		t, err := game.DeserializeTurn(supposedlySerializedTurn)
		if err != nil {
			fmt.Println("could not read your instructions, please try again: " + err.Error())
			continue
		}

		if _, ok := validTurns[t.String()]; ok {
			return t
		} else {
			fmt.Println("invalid turn entered, please try again")
		}
	}
}

// The best AI ever built.
func randomlyChooseValidTurn(validTurns map[string]game.Turn) game.Turn {
	for _, t := range validTurns {
		return t
	}
	panic("no turns to choose. you should've prevented this line from being reached")
}

type GameController struct {
	g      *game.Game
	debug  bool
	agent  *learn.Agent
	state1 learn.State
	action string
}

func New(qs *learn.QContainer, stc *learn.SerializedTurnsCache, debug bool) *GameController {
	learningRateAkaAlpha := 0.1
	rewardsDiscountRateAkaGamma := 0.1
	initialExplorationRateAkaEpsilon := 1.0
	agent := learn.NewAgent(qs, stc, learningRateAkaAlpha, rewardsDiscountRateAkaGamma, initialExplorationRateAkaEpsilon)
	return &GameController{agent: agent, debug: debug}
}

func (gc *GameController) PlayOneGame(numHumanPlayers uint8, stopLearning bool) (*game.Player, game.WinKind) {
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

func (gc *GameController) chooseTurn(validTurns map[string]game.Turn) (game.Turn, bool) {
	if gc.g.IsCurrentPlayerHuman() {
		return readTurnFromStdin(validTurns), false
	} else {
		gc.agent.SetPlayer(gc.g.CurrentPlayer)
		gc.state1 = gc.agent.DetectState()
		var turn game.Turn
		gc.action, turn = gc.agent.EpsilonGreedyAction(gc.state1)
		return turn, true
	}
}

// playOneTurn plays through one turn, and returns whether the game is finished after the turn executes.
func (gc *GameController) playOneTurn() bool {
	g := gc.g

	if g.HasAnyHumans() || gc.debug {
		render.PrintGame(g)
	}

	validTurns := game.ValidTurns(g.Board, g.CurrentRoll, g.CurrentPlayer)
	g.SetValidTurns(validTurns) // IMPORTANT: this must be set for the current turn, always!

	var wasTurnPickedByAI bool
	var chosenTurn game.Turn
	if len(validTurns) == 0 {
		gc.maybePrint("\tcan't do anything this turn, sorry!")
	} else if len(validTurns) == 1 {
		gc.maybePrint("\tthis turn only has 1 option, forcing!")
		chosenTurn = randomlyChooseValidTurn(validTurns)
	} else {
		gc.maybePrint(fmt.Sprintf("\tYour move, %q:", *g.CurrentPlayer))
		chosenTurn, wasTurnPickedByAI = gc.chooseTurn(validTurns)
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

	if wasTurnPickedByAI {
		gc.agent.Learn(gc.state1, gc.action, gc.agent.DetectState(), winAmt)
	}

	if winner != nil {
		return true
	} else {
		gc.g.NextPlayersTurn()
		return false
	}
}