package game

type Game struct {
	Board         *Board
	CurrentPlayer *Player
	CurrentRoll   *Roll
}

func CopyGame(g Game) Game { return g }  // is this a deep copy though?
func GamePointer(g Game) *Game { return &g }

func NewGame() *Game {
	b := &Board{}
	b.setUp()
	// TODO: select the initial player randomly.
	return &Game{b, PCC, newRoll()}
}
