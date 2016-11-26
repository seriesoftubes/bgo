package game

var (
  PC  *Player = playerPointer("O")
  PCC *Player = playerPointer("X")
)

type Player string

func playerPointer(p Player) *Player { return &p }

func (p *Player) Symbol() string {
  if p == nil {
    panic("Invalid player")
  }

  return string(*p)
}
