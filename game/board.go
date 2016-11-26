package game


const (
  NUM_BOARD_POINTS                    uint8  = 24
)

type BoardPoint struct {
  Owner       *Player
  NumCheckers uint8
}

func (p *BoardPoint) Symbol() string {
  if p.Owner == nil {
    panic("No owner of this point")
  }

  return p.Owner.Symbol()
}

type Board struct {
  Points      [NUM_BOARD_POINTS]*BoardPoint
  BarCC, BarC uint8 // # of checkers on each player's bar
}

func (b *Board) setUp() {
  b.Points = [NUM_BOARD_POINTS]*BoardPoint{
    // counter-clockwise player is in bottom-left.
    {PCC, 2}, {}, {}, {}, {}, {PC, 5}, {}, {PC, 3}, {}, {}, {}, {PCC, 5},
    {PC, 5}, {}, {}, {}, {PCC, 3}, {}, {PCC, 5}, {}, {}, {}, {}, {PC, 2},
    //                                                        clockwise player in top-left.
  }
}

func (b *Board) PipCounts() (int, int) {
  var pipC, pipCC int

  for i, p := range b.Points {
    basePips, chex := i+1, int(p.NumCheckers)
    if p.Owner == PC {
      // the clockwise player's closest checker is at points[0].
      pipC += chex * basePips
    } else if p.Owner == PCC {
      // the counter-clockwise player's furthest checker is at points[0].
      pipCC += chex * (int(NUM_BOARD_POINTS) - basePips + 1)
    }
  }

  return pipC, pipCC
}
