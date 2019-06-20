package life

import "fmt"

type CellStatus byte

const (
	blank CellStatus = iota
	alive
	born
	die
)

type Cell struct {
	status CellStatus
}

func NewCell() *Cell {
	return &Cell{status: blank}
}

func (s *Cell) GetStatus() CellStatus {
	return s.status
}

func (s *Cell) SetAlive() {
	s.status = alive
}

func (s *Cell) SetBlank() {
	s.status = blank
}

func (s *Cell) SetBorn() {
	s.status = born
}

func (s *Cell) SetDie() {
	s.status = die
}

func (s Cell) String() string {
	return fmt.Sprintf("Cell:%v", s.status)
}

func (s Cell) Print() string {
	str := " "
	switch s.status {
	case alive:
		str = "*"
	case blank:
		str = " "
	case born:
		str = "&"
	case die:
		str = "#"
	}
	return str
}
