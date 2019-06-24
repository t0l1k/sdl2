package fifteen

import (
	"strconv"
)

type cell struct {
	number int
}

func newCell(number int) *cell {
	return &cell{
		number: number,
	}
}

func (s *cell) getNumber() int {
	return s.number
}

func (s *cell) isBlank() bool {
	return s.number == 0
}

func (s *cell) String() (str string) {
	if s.isBlank() {
		str = "*"
	} else {
		str = strconv.Itoa(int(s.number))
	}
	return str
}
