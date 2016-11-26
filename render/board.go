// Package render renders stuff to the command line.
package render

import (
  "fmt"
  "strconv"
  "strings"

  "github.com/seriesoftubes/bgo/game"
)

const (
  MAX_SYMBOLS_TO_PRINT          uint8  = 5
  BOT_LEFT_POINT_IDX            uint8  = 0
  TOP_LEFT_POINT_IDX            uint8  = game.NUM_BOARD_POINTS - 1
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

var num2alpha         = map[uint8]string{
    0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h", 8: "i", 9: "j",
    10: "k", 11: "l", 12: "m", 13: "n", 14: "o", 15: "p", 16: "q", 17: "r", 18: "s", 19: "t",
    20: "u", 21: "v", 22: "w", 23: "x", 24: "y", 25: "z",
  }

func printBoard(b *game.Board) {
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

      p := b.Points[pointIdx]
      if height == 0 && p.NumCheckers == 0 {
        row += prefix + EMPTY_CHECKERS
      } else {
        if p.NumCheckers > height {
          row += prefix + " " + p.Symbol() + " "
        } else {
          row += prefix + BLANK_SPACE
        }
      }
    }
    topRows = append(topRows, row)
  }

  // Bottom row of the top: a number, if NumCheckers > MAX_SYMBOLS_TO_PRINT
  row = ""
  for pointIdx := TOP_LEFT_POINT_IDX; pointIdx >= TOP_RIGHT_POINT_IDX; pointIdx-- {
    p := b.Points[pointIdx]
    prefix := ""
    if pointIdx == TOP_POINT_TO_THE_RIGHT_OF_MID {
      prefix = TOP_MID_BORDER
    }

    if p.NumCheckers > MAX_SYMBOLS_TO_PRINT {
      pads := " "
      spaces := " "
      if p.NumCheckers > 9 {
        pads = ""
        spaces = " "
      }
      row += prefix + pads + strconv.Itoa(int(p.NumCheckers)) + spaces
    } else {
      row += prefix + BLANK_SPACE
    }
  }
  topRows = append(topRows, row)

  // Top row of the bottom: a number, if NumCheckers > MAX_SYMBOLS_TO_PRINT
  row = ""
  for pointIdx := BOT_LEFT_POINT_IDX; pointIdx <= BOT_RIGHT_POINT_IDX; pointIdx++ {
    p := b.Points[pointIdx]
    prefix := ""
    if pointIdx == BOT_POINT_TO_THE_RIGHT_OF_MID {
      prefix = BOT_MID_BORDER
    }

    if p.NumCheckers > MAX_SYMBOLS_TO_PRINT {
      pads := " "
      spaces := " "
      if p.NumCheckers > 9 {
        pads = ""
        spaces = " "
      }
      row += prefix + pads + strconv.Itoa(int(p.NumCheckers)) + spaces
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

      p := b.Points[pointIdx]
      if height == 1 && p.NumCheckers == 0 {
        row += prefix + EMPTY_CHECKERS
      } else {
        if p.NumCheckers >= height {
          row += prefix + " " + p.Symbol() + " "
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
  fmt.Println(prefix + "y\t" + renderBar(game.PCC, b.BarCC)) // character "y" is reserved for the CC bar.
  fmt.Println(prefix + "z\t" + renderBar(game.PC, b.BarC))   // character "z" is reserved for the C bar.
  fmt.Println(prefix)
  fmt.Println(prefix + "Beared off")
  fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", game.PCC.Symbol(), b.OffCC))
  fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", game.PC.Symbol(), b.OffC))
  fmt.Println(prefix)
  fmt.Println(prefix + "Pipcounts")
  pipC, pipCC := b.PipCounts()
  fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", game.PCC.Symbol(), pipCC))
  fmt.Println(prefix + fmt.Sprintf("\t%s's: %d", game.PC.Symbol(), pipC))
  fmt.Println(prefix + "\n")
}

func renderBar(p *game.Player, numOnBar uint8) string {
  bar := fmt.Sprintf("%s's: ", p.Symbol())
  for i := uint8(0); i < numOnBar && i < MAX_SYMBOLS_TO_PRINT; i++ {
    bar += p.Symbol()
  }
  if numOnBar == 0 {
    bar += strings.TrimSpace(EMPTY_CHECKERS)
  } else if numOnBar > MAX_SYMBOLS_TO_PRINT {
    bar += " " + strconv.Itoa(int(numOnBar))
  }
  return bar
}
