// Package ctrl contains a game controller, which moves a game along through its turns
package ctrl

import (
	"fmt"
	"sync"

	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
	"github.com/seriesoftubes/bgo/game/turn"
	"github.com/seriesoftubes/bgo/game/turngen"
	"github.com/seriesoftubes/bgo/learn"
	"github.com/seriesoftubes/bgo/render"
)

const (
	msgWelcome      = "Welcome to backgammon. Good luck and have fun!"
	msgGameOver     = "\tDONE WITH GAME!"
	msgNoMovesAvail = "\tcan't do anything this turn, sorry!"
	msgForceMove    = "\tthis turn only has 1 option, forcing!"
	msgAskForMove   = "\tYour move, "
	msgChoseMove    = "\tChose move:"
)

var (
	NNVariances   []float32
	nnVariancesMu sync.Mutex
)

func readTurnFromStdin(p plyr.Player, validTurns map[turn.TurnArray]turn.Turn) turn.Turn {
	fmt.Println(msgAskForMove, string(p))
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
	prevBoards map[plyr.Player]*game.Board
}

func New(debug bool) *GameController {
	learningRateAkaAlpha := float32(0.0001)
	rewardsDiscountRateAkaGamma := float32(0.999)
	initialExplorationRateAkaEpsilon := float32(1.0)
	prevBoards := make(map[plyr.Player]*game.Board, 2)
	agent := learn.NewAgent(learningRateAkaAlpha, rewardsDiscountRateAkaGamma, initialExplorationRateAkaEpsilon)
	return &GameController{agent: agent, prevBoards: prevBoards, debug: debug}
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
	avgVariance := gc.agent.AverageNeuralNetworkVariance() // gets avg variance as of the end of a game.
	defer nnVariancesMu.Unlock()
	nnVariancesMu.Lock()
	NNVariances = append(NNVariances, avgVariance)
	gc.agent.ResetNeuralNetworkStats()

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

// playOneTurn plays through one turn, and returns whether the game is finished after the turn executes.
func (gc *GameController) playOneTurn() bool {
	g := gc.g

	if g.HasAnyHumans() || gc.debug {
		render.PrintGame(g)
	}

	validTurns := turngen.ValidTurns(g.Board, g.CurrentRoll, g.CurrentPlayer)

	isComputer := !g.IsCurrentPlayerHuman()
	var currentBoard *game.Board
	if isComputer {
		gc.agent.SetPlayer(g.CurrentPlayer)
		currentBoard = g.Board.Copy()
		if prevBoard, ok := gc.prevBoards[g.CurrentPlayer]; ok {
			gc.agent.LearnNonFinalState(prevBoard, currentBoard)
		}
	}

	var chosenTurn turn.Turn
	if len(validTurns) == 0 {
		gc.maybePrint(msgNoMovesAvail)
	} else if len(validTurns) == 1 {
		gc.maybePrint(msgForceMove)
		chosenTurn = randomlyChooseValidTurn(validTurns)
	} else {
		if gc.g.IsCurrentPlayerHuman() {
			chosenTurn = readTurnFromStdin(g.CurrentPlayer, validTurns)
		} else {
			chosenTurn = gc.agent.EpsilonGreedyAction(currentBoard, validTurns)
		}
	}
	gc.maybePrint(msgChoseMove, chosenTurn)

	gc.prevBoards[g.CurrentPlayer] = currentBoard
	gc.g.Board.MustExecuteTurn(chosenTurn, gc.debug)
	winner, winAmt := gc.g.Board.Winner(), gc.g.Board.WinKind()

	if winner != 0 {
		if isComputer { // special case: a computer won, and they need to learn from that without having to re-run this playOneTurn method.
			previousEnemyBoard := gc.prevBoards[g.CurrentPlayer.Enemy()]  // this can be null if the other player is human.
			gc.agent.LearnFinal(currentBoard, previousEnemyBoard, winAmt) // The `currentBoard` variable still reflects the state before the turn was executed.
		}
		return true
	} else {
		gc.g.NextPlayersTurn()
		return false
	}
}
