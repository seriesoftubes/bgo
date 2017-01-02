package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seriesoftubes/bgo/ctrl"
	"github.com/seriesoftubes/bgo/learn/nnet"
	"github.com/seriesoftubes/bgo/learn/nnet/nnperf"
)

const (
	cmdAverageVariance                           = "avgvar"
	cmdTotalVariance                             = "ttlvar"
	cmdCFG                                       = "cfg"
	cmdprefixMultiplyLearningRate                = "mulr_"
	cmdprefixChangeLearningRateReducerInterval   = "interval_"
	cmdprefixChangeLearningRateReducerMultiplier = "multiplier_"
	cmdHelp                                      = "help"
)

var (
	totalGamesToPlayPtr = flag.Uint64("total_games_to_play", 2000, "The total number of games to play across all goroutines")
	numGoroutinesPtr    = flag.Uint64("goroutines", uint64(runtime.NumCPU()/2), "The number of goroutines to run on")
	epsilonPtr          = flag.Float64("epsilon", 1.0, "The chance (number between 0 and 1.0) that an agent picks a random move instead of an optimal one")
	inFilePathPtr       = flag.String("config_infile", "", "The file that contains the initial neural net config")
	outFilePathPtr      = flag.String("config_outfile", "", "The file that will contain the updated neural net config")
)

var (
	learningRateReducerInterval     = uint64(42000)
	learningRateReductionMultiplier = float32(0.5)
	lrConfigMu                      sync.RWMutex
	gamesPlayed                     uint64
)

func loadNeuralNetwork(filePath string) {
	fmt.Println("loading neural network config from", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("could not open neural network file %q. Skipping. %v\n", filePath, err)
		return
	}
	defer f.Close()

	if existingGamesPlayed, err := nnet.Load(f); err != nil {
		panic("could not deserialize neural network: " + err.Error())
	} else {
		gamesPlayed = existingGamesPlayed
	}

	fmt.Println("neural net loaded!")
}

func saveNeuralNetwork(filePath string) {
	fmt.Println("saving neural net config to", filePath)

	f, err := os.Create(filePath) // always overwrites the existing file.
	if err != nil {
		panic("could not create file: " + err.Error())
	}
	defer f.Close()

	if err := nnet.Save(f, gamesPlayed); err != nil {
		panic("couldnt save neural network: " + err.Error())
	}

	fmt.Println("neural net config saved!")
}

func incrementGamesPlayed() {
	atomic.AddUint64(&gamesPlayed, 1)
	ct := atomic.LoadUint64(&gamesPlayed)

	if ct%500 == 0 {
		fmt.Println(time.Now(), "trained on", ct, "games")
	}

	lrConfigMu.RLock()
	my_interval, my_multiplier := learningRateReducerInterval, learningRateReductionMultiplier
	lrConfigMu.RUnlock()

	if ct%my_interval == 0 {
		fmt.Println("multiplying the neural net's learning rate by", my_multiplier)
		nnet.MultiplyLearningRate(my_multiplier)
	}
}

func filePathFromFlag(fp *string) string {
	if fp == nil || *fp == "" {
		u, err := user.Current()
		if err != nil {
			panic("could not get current OS user: " + err.Error())
		}
		return fmt.Sprintf("%s/Desktop/bgo_nnet.json", u.HomeDir)
	}

	return *fp
}

func writeVarianceLogs(startGamesPlayed uint64, filePath string) {
	fmt.Println("saving variance data to", filePath)

	f, err := os.Create(filePath) // always overwrites the existing file.
	if err != nil {
		panic("could not create file: " + err.Error())
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	writeLine := func(ln string) {
		if _, err = w.WriteString(ln); err != nil {
			panic(fmt.Sprintf("w.WriteString(%q) error: %v", ln, err))
		}
	}

	writeLine("GamesPlayed\tAvgVariance\n")
	for i, v := range nnperf.GameAverageVariances(-1, true) {
		writeLine(fmt.Sprintf("%v\t%v\n", uint64(i)+1+startGamesPlayed, v))
		if i%1000 == 0 {
			w.Flush()
		}
	}
	w.Flush()

	fmt.Println("done saving variance data!")
}

func onHelpCmd() {
	fmt.Println("valid commands are:")
	fmt.Println(cmdAverageVariance)
	fmt.Println(cmdTotalVariance)
	fmt.Println(cmdCFG)
	fmt.Println(cmdprefixMultiplyLearningRate, "(plus a number, like mulr_1.23)")
	fmt.Println(cmdprefixChangeLearningRateReducerInterval, "(plus a number like bla_123123)")
	fmt.Println(cmdprefixChangeLearningRateReducerMultiplier, "(plus a number like bla_0.8")
	fmt.Println(cmdHelp)
}

func float32FromCommand(cmd string) (float32, error) {
	split := strings.Split(cmd, "_")
	if len(split) != 2 {
		return 0, fmt.Errorf("invalid command (should look like 'blabla_5.5') got", cmd)
	}

	factor, err := strconv.ParseFloat(split[1], 32)
	if err != nil {
		return 0, fmt.Errorf("invalid number %s in command %q: %v", split[1], cmd, err)
	}

	return float32(factor), nil
}

func onMulrCmd(cmd string) {
	factor, err := float32FromCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("multiplying learning rate by", factor)
	nnet.MultiplyLearningRate(factor)
}

func onChangeLearningRateReducerIntervalCmd(cmd string) {
	defer lrConfigMu.Unlock()
	lrConfigMu.Lock()

	newInterval, err := float32FromCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("setting learningRateReducerInterval to", uint64(newInterval))
	learningRateReducerInterval = uint64(newInterval)
}

func onChangeLearningRateReducerMultiplierCmd(cmd string) {
	defer lrConfigMu.Unlock()
	lrConfigMu.Lock()

	newMultiplier, err := float32FromCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("setting learningRateReductionMultiplier to", newMultiplier)
	learningRateReductionMultiplier = newMultiplier
}

func onCfgCmd() {
	learningRate, decayRate := nnet.LearningParams()
	fmt.Println("learningRate", learningRate)
	fmt.Println("decayRate", decayRate)
	fmt.Println("learningRateReductionMultiplier", learningRateReductionMultiplier)
	fmt.Println("learningRateReducerInterval", learningRateReducerInterval)
}

func readCommands(doneChan chan bool) {
	var previousCmd string

	for {
		select {
		case <-doneChan:
			fmt.Println("done reading model training commands")
			return
		default:
			fmt.Println("\n\nenter a command\n")
		}

		var rawCmd string
		fmt.Scanln(&rawCmd)
		cmd := strings.ToLower(strings.TrimSpace(rawCmd))

		cmdWasToRepeatPreviousCommand := cmd == "d" || cmd == "r"
		if cmdWasToRepeatPreviousCommand {
			fmt.Println("repeating command", previousCmd)
			cmd = previousCmd
		}

		if cmd == cmdAverageVariance {
			for _, v := range nnperf.GameAverageVariances(30, false) {
				fmt.Println(v)
			}
		} else if cmd == cmdTotalVariance {
			for _, v := range nnperf.GameTotalVariances(30, false) {
				fmt.Println(v)
			}
		} else if cmd == cmdCFG {
			onCfgCmd()
		} else if strings.HasPrefix(cmd, cmdprefixMultiplyLearningRate) {
			onMulrCmd(cmd)
		} else if strings.HasPrefix(cmd, cmdprefixChangeLearningRateReducerInterval) {
			onChangeLearningRateReducerIntervalCmd(cmd)
		} else if strings.HasPrefix(cmd, cmdprefixChangeLearningRateReducerMultiplier) {
			onChangeLearningRateReducerMultiplierCmd(cmd)
		} else if cmd == cmdHelp {
			onHelpCmd()
		}

		if !cmdWasToRepeatPreviousCommand {
			previousCmd = cmd
		}
	}
}

func main() {
	flag.Parse()

	infilePath := filePathFromFlag(inFilePathPtr)
	outfilePath := filePathFromFlag(outFilePathPtr)
	if !strings.HasSuffix(infilePath, ".json") || !strings.HasSuffix(outfilePath, ".json") {
		panic("both the infile and outfile must have a .json suffix.")
	}
	varianceLogsFilePath := strings.Replace(outfilePath, ".json", "_variance_report.txt", 1)

	gamesToPlayPerGoroutine := *totalGamesToPlayPtr / *numGoroutinesPtr
	fmt.Printf("training on %d games (%d goroutines X %d games per goroutine)...\n", *numGoroutinesPtr*gamesToPlayPerGoroutine, *numGoroutinesPtr, gamesToPlayPerGoroutine)

	loadNeuralNetwork(infilePath)
	doneChan := make(chan bool, 1)
	go readCommands(doneChan)
	start := time.Now()
	startGamesPlayed := gamesPlayed

	var wg sync.WaitGroup
	for i := uint64(0); i < *numGoroutinesPtr; i++ {
		wg.Add(1)
		go func() {
			mgr := ctrl.New(false)
			for j := uint64(0); j < gamesToPlayPerGoroutine; j++ {
				mgr.PlayOneGame(0, false) // Play 1 game with 0 humans and don't stop learning!
				mgr.TransmitStatsFromMostRecentGame()
				incrementGamesPlayed()
			}
			mgr.WaitForStats()
			wg.Done()
		}()
	}
	wg.Wait()
	doneChan <- true
	close(doneChan)
	fmt.Println("trained", gamesPlayed-startGamesPlayed, "times in", time.Since(start))

	saveNeuralNetwork(outfilePath)
	writeVarianceLogs(startGamesPlayed, varianceLogsFilePath)

	mgr := ctrl.New(true)
	mgr.PlayOneGame(1, true /* stopLearning=true */)
}
