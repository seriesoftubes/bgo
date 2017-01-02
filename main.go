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
	totalGamesToPlayPtr = flag.Uint64("total_games_to_play", 5000, "The total number of games to play across all goroutines")
	numGoroutinesPtr    = flag.Uint64("goroutines", uint64(runtime.NumCPU()/2), "The number of goroutines to run on")
	epsilonPtr          = flag.Float64("epsilon", 1.0, "The chance (number between 0 and 1.0) that an agent picks a random move instead of an optimal one")
	inFilePathPtr       = flag.String("config_infile", "", "The file that contains the initial neural net config")
	outFilePathPtr      = flag.String("config_outfile", "", "The file that will contain the updated neural net config")
)

type (
	learningRateManager struct {
		sync.RWMutex
		interval   uint64  // Every `interval` games, we update the learningRate
		multiplier float32 // Every `interval` games, we multiply the learningRate by `multiplier`
	}

	// Gotta train'em all! Poke-model!
	pokemodelTrainer struct {
		configInFile, configOutFile, varianceLogsFilePath string
		hasLoadedNN                                       bool
		startGamesPlayed, gamesPlayed                     uint64
		lrManager                                         *learningRateManager
	}
)

func newTrainer(configInFile, configOutFile string) *pokemodelTrainer {
	if !strings.HasSuffix(configInFile, ".json") || !strings.HasSuffix(configOutFile, ".json") {
		panic("both the infile and outfile must have a .json suffix.")
	}

	varianceLogsFilePath := strings.Replace(configOutFile, ".json", "_variance_report.txt", 1)
	return &pokemodelTrainer{
		configInFile:         configInFile,
		configOutFile:        configOutFile,
		varianceLogsFilePath: varianceLogsFilePath,
		lrManager:            &learningRateManager{interval: 42000, multiplier: 0.00001},
	}
}

func onHelpCmd() {
	fmt.Println("valid commands are:")
	fmt.Println("'d' or 'r' to repeat the previous command")
	fmt.Println(cmdAverageVariance)
	fmt.Println(cmdTotalVariance)
	fmt.Println(cmdCFG)
	fmt.Println(cmdprefixMultiplyLearningRate, "(plus a number, like mulr_1.23)")
	fmt.Println(cmdprefixChangeLearningRateReducerInterval, "(plus a number like bla_123123)")
	fmt.Println(cmdprefixChangeLearningRateReducerMultiplier, "(plus a number like bla_0.8")
	fmt.Println(cmdHelp)
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

func (lrm *learningRateManager) maybeChangeLearningRate(numGamesCompleted uint64) {
	if my_multiplier, my_interval := lrm.params(); numGamesCompleted%my_interval == 0 {
		fmt.Println("multiplying the neural net's learning rate by", my_multiplier)
		nnet.MultiplyLearningRate(my_multiplier)
	}
}

func (lrm *learningRateManager) setInterval(newInterval uint64) {
	defer lrm.Unlock()
	lrm.Lock()

	fmt.Println("setting learningRateReducerInterval to", newInterval)
	lrm.interval = newInterval
}

func (lrm *learningRateManager) setMultiplier(newMultiplier float32) {
	defer lrm.Unlock()
	lrm.Lock()

	fmt.Println("setting learningRateReductionMultiplier to", newMultiplier)
	lrm.multiplier = newMultiplier
}

func (lrm *learningRateManager) params() (float32, uint64) {
	defer lrm.RUnlock()
	lrm.RLock()
	return lrm.multiplier, lrm.interval
}

func (pt *pokemodelTrainer) onCfgCmd() {
	learningRate, decayRate := nnet.LearningParams()
	multiplier, interval := pt.lrManager.params()

	fmt.Println("learningRate", learningRate)
	fmt.Println("decayRate", decayRate)
	fmt.Println("learningRateReductionMultiplier", multiplier)
	fmt.Println("learningRateReductionInterval", interval)
	fmt.Println("gamesPlayed", atomic.LoadUint64(&pt.gamesPlayed))
}

func (pt *pokemodelTrainer) onGameCompleted() {
	atomic.AddUint64(&pt.gamesPlayed, 1)
	ct := atomic.LoadUint64(&pt.gamesPlayed)
	if ct%500 == 0 {
		fmt.Println(time.Now(), "trained on", ct, "games")
	}
	pt.lrManager.maybeChangeLearningRate(ct)
}

func (pt *pokemodelTrainer) loadNeuralNetwork() {
	filePath := pt.configInFile
	fmt.Println("loading neural network config from", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("could not open neural network file %q. Skipping. %v\n", filePath, err)
		pt.hasLoadedNN = true
		return
	}
	defer f.Close()

	if existingGamesPlayed, err := nnet.Load(f); err != nil {
		panic("could not deserialize neural network: " + err.Error())
	} else {
		atomic.StoreUint64(&pt.gamesPlayed, existingGamesPlayed)
		atomic.StoreUint64(&pt.startGamesPlayed, existingGamesPlayed)
	}

	pt.hasLoadedNN = true
	fmt.Println("neural net loaded!")
}

func (pt *pokemodelTrainer) saveNeuralNetwork(waitForWrites bool) {
	filePath := pt.configOutFile
	fmt.Println("saving neural net config to", filePath)

	f, err := os.Create(filePath) // always overwrites the existing file.
	if err != nil {
		panic("could not create file: " + err.Error())
	}
	defer f.Close()

	if err := nnet.Save(f, atomic.LoadUint64(&pt.gamesPlayed), waitForWrites); err != nil {
		panic("couldnt save neural network: " + err.Error())
	}

	fmt.Println("neural net config saved!")
}

func (pt *pokemodelTrainer) writeVarianceLogs(waitForWrites bool) {
	filePath := pt.varianceLogsFilePath
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
	startGamesPlayed := atomic.LoadUint64(&pt.startGamesPlayed)
	for i, v := range nnperf.GameAverageVariances(-1, waitForWrites) {
		writeLine(fmt.Sprintf("%v\t%v\n", uint64(i)+1+startGamesPlayed, v))
		if i%1000 == 0 {
			w.Flush()
		}
	}
	w.Flush()

	fmt.Println("done saving variance data!")
}

func (pt *pokemodelTrainer) onChangeLearningRateReducerIntervalCmd(cmd string) {
	newInterval, err := float32FromCommand(cmd)
	if err != nil {
		fmt.Println("could not parse number from cmd", cmd, err.Error())
		return
	}
	pt.lrManager.setInterval(uint64(newInterval))
}

func (pt *pokemodelTrainer) onChangeLearningRateReducerMultiplierCmd(cmd string) {
	newMultiplier, err := float32FromCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pt.lrManager.setMultiplier(newMultiplier)
}

// TODO: split this into its own struct.
func (pt *pokemodelTrainer) readCommands(doneChan chan bool) {
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
			fmt.Printf("repeating command: %q\n", previousCmd)
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
			pt.onCfgCmd()
		} else if strings.HasPrefix(cmd, cmdprefixMultiplyLearningRate) {
			onMulrCmd(cmd)
		} else if strings.HasPrefix(cmd, cmdprefixChangeLearningRateReducerInterval) {
			pt.onChangeLearningRateReducerIntervalCmd(cmd)
		} else if strings.HasPrefix(cmd, cmdprefixChangeLearningRateReducerMultiplier) {
			pt.onChangeLearningRateReducerMultiplierCmd(cmd)
		} else if cmd == cmdHelp {
			onHelpCmd()
		}

		if !cmdWasToRepeatPreviousCommand {
			previousCmd = cmd
		}
	}
}

func (pt *pokemodelTrainer) train(numGoroutines, gamesToPlayPerGoroutine uint64) {
	if !pt.hasLoadedNN {
		pt.loadNeuralNetwork()
	}

	fmt.Printf("training on %d games (%d goroutines X %d games per goroutine)...\n", numGoroutines*gamesToPlayPerGoroutine, numGoroutines, gamesToPlayPerGoroutine)
	start := time.Now()

	doneChan := make(chan bool, 1)
	go pt.readCommands(doneChan)

	var wg sync.WaitGroup
	for i := uint64(0); i < numGoroutines; i++ {
		wg.Add(1)

		go func() {
			mgr := ctrl.New(false)
			for j := uint64(0); j < gamesToPlayPerGoroutine; j++ {
				mgr.PlayOneGame(0, false) // Play 1 game with 0 humans and don't stop learning!
				mgr.TransmitStatsFromMostRecentGame()
				pt.onGameCompleted()
			}
			mgr.WaitForStats()
			wg.Done()
		}()
	}
	wg.Wait()

	doneChan <- true
	close(doneChan)
	fmt.Printf("trained %d times in %v\n", atomic.LoadUint64(&pt.gamesPlayed)-atomic.LoadUint64(&pt.startGamesPlayed), time.Since(start))
}

func main() {
	flag.Parse()

	trainer := newTrainer(filePathFromFlag(inFilePathPtr), filePathFromFlag(outFilePathPtr))
	trainer.loadNeuralNetwork()
	trainer.train(*numGoroutinesPtr, *totalGamesToPlayPtr / *numGoroutinesPtr)
	trainer.saveNeuralNetwork(true /* waitForWrites=true*/)
	trainer.writeVarianceLogs(true /* waitForWrites=true*/)

	mgr := ctrl.New(true /* debug=true*/)
	mgr.PlayOneGame(1, true /* stopLearning=true */)
}
