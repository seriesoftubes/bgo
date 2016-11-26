package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Player uint8

func (p *Player) symbol() string {
	if p == pC {
		return PLAYER_C_SYMBOL
	} else if p == pCC {
		return PLAYER_CC_SYMBOL
	}
	panic("Invalid player")
}

const (
	MIN_DICE_AMT = 1
	MAX_DICE_AMT = 6
	NUM_POINTS                    uint8  = 24
	PLAYER_C                      Player = 1 // Clockwise
	PLAYER_C_SYMBOL                      = "O"
	PLAYER_CC                     Player = 2 // Counter-clockwise
	PLAYER_CC_SYMBOL                     = "X"
	MAX_SYMBOLS_TO_PRINT          uint8  = 5
	BOT_LEFT_POINT_IDX            uint8  = 0
	TOP_LEFT_POINT_IDX            uint8  = NUM_POINTS - 1
	BOT_RIGHT_POINT_IDX           uint8  = 11
	TOP_RIGHT_POINT_IDX           uint8  = 12
	TOP_POINT_TO_THE_RIGHT_OF_MID uint8  = 17
	BOT_POINT_TO_THE_RIGHT_OF_MID uint8  = 6
	BORDER                               = "==="
	TOP_MID_BORDER                       = "|m|"
	BOT_MID_BORDER                       = "|w|"
	EMPTY_CHECKERS                       = " - "
	BLANK_SPACE                          = "   "
)

func playerPointer(p Player) *Player { return &p }

var (
	pC        *Player = playerPointer(PLAYER_C)
	pCC       *Player = playerPointer(PLAYER_CC)
	num2alpha         = map[uint8]string{
		0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h", 8: "i", 9: "j",
		10: "k", 11: "l", 12: "m", 13: "n", 14: "o", 15: "p", 16: "q", 17: "r", 18: "s", 19: "t",
		20: "u", 21: "v", 22: "w", 23: "x", 24: "y", 25: "z",
	}
	alpha2num = map[string]uint8{
		"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7, "i": 8, "j": 9,
		"k": 10, "l": 11, "m": 12, "n": 13, "o": 14, "p": 15, "q": 16, "r": 17, "s": 18, "t": 19,
		"u": 20, "v": 21, "w": 22, "x": 23, "y": 24, "z": 25,
	}
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type BoardPoint struct {
	owner       *Player
	numCheckers uint8
}

func (p *BoardPoint) symbol() string {
	if p.owner == nil {
		panic("No owner of this point")
	}

	return p.owner.symbol()
}

type Board struct {
	points      [NUM_POINTS]*BoardPoint
	barCC, barC uint8 // # of checkers on each player's bar
}

func (b *Board) setUp() {
	b.points = [NUM_POINTS]*BoardPoint{
		// counter-clockwise player is in bottom-left.
		{pCC, 2}, {}, {}, {}, {}, {pC, 5}, {}, {pC, 3}, {}, {}, {}, {pCC, 5},
		{pC, 5}, {}, {}, {}, {pCC, 3}, {}, {pCC, 5}, {}, {}, {}, {}, {pC, 2},
		//                                                        clockwise player in top-left.
	}
}

func (b *Board) pipCounts() (int, int) {
	var pipC, pipCC int

	for i, p := range b.points {
		basePips, chex := i+1, int(p.numCheckers)
		if p.owner == pC {
			// the clockwise player's closest checker is at points[0].
			pipC += chex * basePips
		} else if p.owner == pCC {
			// the counter-clockwise player's furthest checker is at points[0].
			pipCC += chex * (int(NUM_POINTS) - basePips + 1)
		}
	}

	return pipC, pipCC
}

func (b *Board) print() {
	var topRows, botRows []string

	// Letters above the top border
	row := ""
	for pointIdx := TOP_LEFT_POINT_IDX; pointIdx >= TOP_RIGHT_POINT_IDX; pointIdx-- {
		prefix := " "
		if pointIdx == TOP_POINT_TO_THE_RIGHT_OF_MID {
			prefix = "   " + " "
		}
		row += prefix + num2alpha[pointIdx] + " "
	}
	topRows = append(topRows, row)

	// Top border
	row = ""
	for pointIdx := TOP_LEFT_POINT_IDX + 1; pointIdx >= TOP_RIGHT_POINT_IDX; pointIdx-- {
		row += BORDER
	}
	topRows = append(topRows, row)

	// Checkers, up to the max height.
	for height := uint8(0); height < MAX_SYMBOLS_TO_PRINT; height++ {
		row = ""
		for pointIdx := TOP_LEFT_POINT_IDX; pointIdx >= TOP_RIGHT_POINT_IDX; pointIdx-- {
			prefix := ""
			if pointIdx == TOP_POINT_TO_THE_RIGHT_OF_MID {
				prefix = TOP_MID_BORDER
			}

			p := b.points[pointIdx]
			if height == 0 && p.numCheckers == 0 {
				row += prefix + EMPTY_CHECKERS
			} else {
				if p.numCheckers > height {
					row += prefix + " " + p.symbol() + " "
				} else {
					row += prefix + BLANK_SPACE
				}
			}
		}
		topRows = append(topRows, row)
	}

	// Bottom row of the top: a number, if numCheckers > MAX_SYMBOLS_TO_PRINT
	row = ""
	for pointIdx := TOP_LEFT_POINT_IDX; pointIdx >= TOP_RIGHT_POINT_IDX; pointIdx-- {
		p := b.points[pointIdx]
		prefix := ""
		if pointIdx == TOP_POINT_TO_THE_RIGHT_OF_MID {
			prefix = TOP_MID_BORDER
		}

		if p.numCheckers > MAX_SYMBOLS_TO_PRINT {
			pads := " "
			spaces := " "
			if p.numCheckers > 9 {
				pads = ""
				spaces = " "
			}
			row += prefix + pads + strconv.Itoa(int(p.numCheckers)) + spaces
		} else {
			row += prefix + BLANK_SPACE
		}
	}
	topRows = append(topRows, row)

	// Top row of the bottom: a number, if numCheckers > MAX_SYMBOLS_TO_PRINT
	row = ""
	for pointIdx := BOT_LEFT_POINT_IDX; pointIdx <= BOT_RIGHT_POINT_IDX; pointIdx++ {
		p := b.points[pointIdx]
		prefix := ""
		if pointIdx == BOT_POINT_TO_THE_RIGHT_OF_MID {
			prefix = BOT_MID_BORDER
		}

		if p.numCheckers > MAX_SYMBOLS_TO_PRINT {
			pads := " "
			spaces := " "
			if p.numCheckers > 9 {
				pads = ""
				spaces = " "
			}
			row += prefix + pads + strconv.Itoa(int(p.numCheckers)) + spaces
		} else {
			row += prefix + BLANK_SPACE
		}
	}
	botRows = append(botRows, row)

	// Checkers, from the max height, down to 0.
	for height := MAX_SYMBOLS_TO_PRINT; height > uint8(0); height-- {
		row = ""
		for pointIdx := BOT_LEFT_POINT_IDX; pointIdx <= BOT_RIGHT_POINT_IDX; pointIdx++ {
			prefix := ""
			if pointIdx == BOT_POINT_TO_THE_RIGHT_OF_MID {
				prefix = BOT_MID_BORDER
			}

			p := b.points[pointIdx]
			if height == 1 && p.numCheckers == 0 {
				row += prefix + EMPTY_CHECKERS
			} else {
				if p.numCheckers >= height {
					row += prefix + " " + p.symbol() + " "
				} else {
					row += prefix + BLANK_SPACE
				}
			}
		}
		botRows = append(botRows, row)
	}

	// Bottom border.
	row = ""
	for pointIdx := BOT_LEFT_POINT_IDX; pointIdx <= BOT_RIGHT_POINT_IDX+1; pointIdx++ {
		row += BORDER
	}
	botRows = append(botRows, row)

	// Letters below the bottom border.
	row = ""
	for pointIdx := BOT_LEFT_POINT_IDX; pointIdx <= BOT_RIGHT_POINT_IDX; pointIdx++ {
		prefix := " "
		if pointIdx == BOT_POINT_TO_THE_RIGHT_OF_MID {
			prefix = "   " + " "
		}
		row += prefix + num2alpha[pointIdx] + " "
	}
	botRows = append(botRows, row)

	// Print the whole board.
	prefix := "\t"
	fmt.Println(prefix + "\n")
	for _, row := range topRows {
		fmt.Println(prefix + row)
	}
	fmt.Println(prefix + "\n")
	for _, row := range botRows {
		fmt.Println(prefix + row)
	}
	fmt.Println(prefix + "\n")
	fmt.Println(prefix + "The bar")
	fmt.Println(prefix + "y\t" + renderBar(pCC, b.barCC)) // character "y" is reserved for the CC bar.
	fmt.Println(prefix + "z\t" + renderBar(pC, b.barC))   // character "z" is reserved for the C bar.
	fmt.Println(prefix)
	fmt.Println(prefix + "Pipcounts")
	pipC, pipCC := b.pipCounts()
	fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", pCC.symbol(), pipCC))
	fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", pC.symbol(), pipC))
	fmt.Println(prefix + "\n")
}

func renderBar(p *Player, numOnBar uint8) string {
	bar := fmt.Sprintf("%s's: ", p.symbol())
	for i := uint8(0); i < numOnBar && i < MAX_SYMBOLS_TO_PRINT; i++ {
		bar += p.symbol()
	}
	if numOnBar == 0 {
		bar += strings.TrimSpace(EMPTY_CHECKERS)
	} else if numOnBar > MAX_SYMBOLS_TO_PRINT {
		bar += " " + strconv.Itoa(int(numOnBar))
	}
	return bar
}

type Roll [2]uint8

func (r *Roll) reverse() *Roll {
	r[1], r[0] = r[0], r[1]
	return r
}

func randBetween(min, max int) uint8 { return uint8(rnd.Intn(max-min) + min) }
func newRoll() *Roll                 { return &Roll{randBetween(MIN_DICE_AMT, MAX_DICE_AMT), randBetween(MIN_DICE_AMT, MAX_DICE_AMT)} }

type Game struct {
	board         *Board
	currentPlayer *Player
	currentRoll   *Roll
}

func (g *Game) print() {
	fmt.Println("\n\tCurrent player:\t", g.currentPlayer.symbol())
	fmt.Println("\n\tDice:\t\t", *g.currentRoll)
	g.board.print()
}

func main() {
	b := &Board{}
	b.setUp()
	g := &Game{board: b, currentPlayer: pCC, currentRoll: newRoll()}
	g.print()
}
