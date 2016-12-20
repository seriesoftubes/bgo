// Package render renders stuff to the command line.
package render

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/seriesoftubes/bgo/constants"
	"github.com/seriesoftubes/bgo/game"
	"github.com/seriesoftubes/bgo/game/plyr"
)

const (
	maxCheckersToPrint      uint8 = 5 // E.g., "OOOOO"
	botLeftPointIdx         uint8 = 0
	topLeftPointIdx         uint8 = constants.NUM_BOARD_POINTS - 1
	botRightPointIdx        uint8 = 11
	topRightPointIdx        uint8 = 12
	topPointToTheRightOfMid uint8 = 17
	botPointToTheRightOfMid uint8 = 6
	border                        = "==="
	topMidBorder                  = "|m|"
	botMidBorder                  = "|w|"
	emptyCheckers                 = " - "
	blankSpace                    = "   "
)

func PrintBoard(b *game.Board) {
	if winner := b.Winner(); winner != 0 {
		fmt.Println(fmt.Sprintf("\n\n\t\tWINNER: %q (won %d points)", string(winner), b.WinKind()))
	}

	var topRows, botRows []string

	// Letters above the top border
	row := ""
	for pointIdx := topLeftPointIdx; pointIdx >= topRightPointIdx; pointIdx-- {
		prefix := " "
		if pointIdx == topPointToTheRightOfMid {
			prefix = "   " + " "
		}
		row += prefix + constants.Num2Alpha[pointIdx] + " "
	}
	topRows = append(topRows, row)

	// Top border
	row = ""
	for pointIdx := topLeftPointIdx + 1; pointIdx >= topRightPointIdx; pointIdx-- {
		row += border
	}
	topRows = append(topRows, row)

	// Checkers, up to the max height.
	for height := uint8(0); height < maxCheckersToPrint; height++ {
		row = ""
		for pointIdx := topLeftPointIdx; pointIdx >= topRightPointIdx; pointIdx-- {
			prefix := ""
			if pointIdx == topPointToTheRightOfMid {
				prefix = topMidBorder
			}

			p := b.Points[pointIdx]
			if height == 0 && p.NumCheckers == 0 {
				row += prefix + emptyCheckers
			} else {
				if p.NumCheckers > height {
					row += prefix + " " + p.Symbol() + " "
				} else {
					row += prefix + blankSpace
				}
			}
		}
		topRows = append(topRows, row)
	}

	// Bottom row of the top: a number, if NumCheckers > maxCheckersToPrint
	row = ""
	for pointIdx := topLeftPointIdx; pointIdx >= topRightPointIdx; pointIdx-- {
		p := b.Points[pointIdx]
		prefix := ""
		if pointIdx == topPointToTheRightOfMid {
			prefix = topMidBorder
		}

		if p.NumCheckers > maxCheckersToPrint {
			pads := " "
			spaces := " "
			if p.NumCheckers > 9 {
				pads = ""
				spaces = " "
			}
			row += prefix + pads + strconv.Itoa(int(p.NumCheckers)) + spaces
		} else {
			row += prefix + blankSpace
		}
	}
	topRows = append(topRows, row)

	// Top row of the bottom: a number, if NumCheckers > maxCheckersToPrint
	row = ""
	for pointIdx := botLeftPointIdx; pointIdx <= botRightPointIdx; pointIdx++ {
		p := b.Points[pointIdx]
		prefix := ""
		if pointIdx == botPointToTheRightOfMid {
			prefix = botMidBorder
		}

		if p.NumCheckers > maxCheckersToPrint {
			pads := " "
			spaces := " "
			if p.NumCheckers > 9 {
				pads = ""
				spaces = " "
			}
			row += prefix + pads + strconv.Itoa(int(p.NumCheckers)) + spaces
		} else {
			row += prefix + blankSpace
		}
	}
	botRows = append(botRows, row)

	// Checkers, from the max height, down to 0.
	for height := maxCheckersToPrint; height > uint8(0); height-- {
		row = ""
		for pointIdx := botLeftPointIdx; pointIdx <= botRightPointIdx; pointIdx++ {
			prefix := ""
			if pointIdx == botPointToTheRightOfMid {
				prefix = botMidBorder
			}

			p := b.Points[pointIdx]
			if height == 1 && p.NumCheckers == 0 {
				row += prefix + emptyCheckers
			} else {
				if p.NumCheckers >= height {
					row += prefix + " " + p.Symbol() + " "
				} else {
					row += prefix + blankSpace
				}
			}
		}
		botRows = append(botRows, row)
	}

	// Bottom border.
	row = ""
	for pointIdx := botLeftPointIdx; pointIdx <= botRightPointIdx+1; pointIdx++ {
		row += border
	}
	botRows = append(botRows, row)

	// Letters below the bottom border.
	row = ""
	for pointIdx := botLeftPointIdx; pointIdx <= botRightPointIdx; pointIdx++ {
		prefix := " "
		if pointIdx == botPointToTheRightOfMid {
			prefix = "   " + " "
		}
		row += prefix + constants.Num2Alpha[pointIdx] + " "
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
	fmt.Println(prefix + constants.LETTER_BAR_CC + "\t" + renderBar(plyr.PCC, b.BarCC))
	fmt.Println(prefix + constants.LETTER_BAR_C + "\t" + renderBar(plyr.PC, b.BarC))
	fmt.Println(prefix)
	fmt.Println(prefix + "Beared off")
	fmt.Println(prefix + fmt.Sprintf("\t%s's: %d\t\t%s's: %d", plyr.PCC.Symbol(), b.OffCC, plyr.PC.Symbol(), b.OffC))
	fmt.Println(prefix + "Pipcounts")
	pipC, pipCC := b.PipCounts()
	fmt.Println(prefix + fmt.Sprintf("\t%s's: %d\t%s's: %d", plyr.PCC.Symbol(), pipCC, plyr.PC.Symbol(), pipC))
	fmt.Println(prefix + "\n")
}

func renderBar(p plyr.Player, numOnBar uint8) string {
	bar := fmt.Sprintf("%s's: ", p.Symbol())
	for i := uint8(0); i < numOnBar && i < maxCheckersToPrint; i++ {
		bar += p.Symbol()
	}
	if numOnBar == 0 {
		bar += strings.TrimSpace(emptyCheckers)
	} else if numOnBar > maxCheckersToPrint {
		bar += " " + strconv.Itoa(int(numOnBar))
	}
	return bar
}
