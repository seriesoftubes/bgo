package game

type Game struct {
  Board         *Board
  CurrentPlayer *Player
  CurrentRoll   *Roll
}

func NewGame() *Game {
  b := &Board{}
  b.setUp()
  return &Game{b, PCC, newRoll()}
}
