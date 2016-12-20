package plyr

import (
	"github.com/seriesoftubes/bgo/constants"
)

var (
	PC  Player = 'O' // Clockwise player
	PCC Player = 'X' // Counter-clockwise player
)

type Player byte

func (p Player) Enemy() Player {
	if p == PCC {
		return PC
	}
	return PCC
}

func (p Player) HomePointIndices() (uint8, uint8) {
	endIdx := constants.NUM_POINTS_IN_HOME_BOARD - 1
	if p == PCC {
		endIdx = constants.NUM_BOARD_POINTS - 1
	}
	startIdx := endIdx - constants.NUM_POINTS_IN_HOME_BOARD + 1
	return startIdx, endIdx
}

func (p Player) Symbol() string {
	return string(p)
}
