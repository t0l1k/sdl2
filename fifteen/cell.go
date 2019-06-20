package fifteen

import (
	"strconv"
)

type Cell struct {
	number int
}

func NewCell(number int) *Cell {
	return &Cell{
		number: number,
	}
}

func (s *Cell) GetNumber() int {
	return s.number
}

func (s *Cell) IsBlank() bool {
	return s.number == 0
}

func (s *Cell) String() (str string) {
	if s.IsBlank() {
		str = "*"
	} else {
		str = strconv.Itoa(int(s.number))
	}
	return str
}
