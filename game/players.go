package game

var (
	PC  *Player = playerPointer("O") // Clockwise player
	PCC *Player = playerPointer("X") // Counter-clockwise player
)

const numPointsInHomeBoard uint8 = 6

type Player string

func playerPointer(p Player) *Player { return &p }

func (p *Player) homePointIndices() (uint8, uint8) {
	endIdx := numPointsInHomeBoard - 1
	if p == PCC {
		endIdx = NUM_BOARD_POINTS - 1
	}
	startIdx := endIdx - numPointsInHomeBoard + 1
	return startIdx, endIdx
}

func (p *Player) Symbol() string {
	return string(*p)
}
