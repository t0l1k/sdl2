package fifteen

import (
	"math/rand"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

type point struct {
	x, y int
}

type Field struct {
	size  int
	board []*Cell
	blank point
}

func NewField(size int) *Field {
	var (
		board []*Cell
		blank point
	)
	for i := 0; i < size*size; i++ {
		if i >= size*size-1 {
			board = append(board, NewCell(0))
			blank = point{i % size, i / size}
		} else {
			board = append(board, NewCell(i+1))
		}
	}
	return &Field{
		size:  size,
		board: board,
		blank: blank,
	}
}

func (s *Field) GetBoard() []*Cell {
	return s.board
}

func (s *Field) Win() bool {
	i := 0
	for _, cell := range s.board {
		if cell.GetNumber() == i+1 {
			i++
		} else if cell.GetNumber() != i+1 && cell.GetNumber() != 0 {
			return false
		}
	}
	if s.board[i].GetNumber() == 0 && i == len(s.board)-1 {
		return true
	}
	return false
}

func (s *Field) Shuffle(count int) {
	for i := 0; i < count; i++ {
		direction := rand.Intn(4)
		s.move(direction)
	}
}

func (s *Field) Move(idx int) bool {
	x, y := s.getPos(idx)
	if s.blank.x == x {
		if s.blank.y > y {
			i := s.blank.y
			for s.blank.y <= i && y != i {
				s.move(UP)
				i--
			}
		} else if s.blank.y < y {
			i := s.blank.y
			for s.blank.y >= i && y != i {
				s.move(DOWN)
				i++
			}
		}
		return true
	} else if s.blank.y == y {
		if s.blank.x > x {
			i := s.blank.x
			for s.blank.x <= i && x != i {
				s.move(LEFT)
				i--
			}
		} else if s.blank.x < x {
			i := s.blank.x
			for s.blank.x >= i && x != i {
				s.move(RIGHT)
				i++
			}
		}
		return true
	}
	return false
}

func (s *Field) move(direction int) {
	x, y := s.blank.x, s.blank.y
	switch direction {
	case UP:
		y--
		if s.isEdge(x, y) {
			y++
		}
	case DOWN:
		y++
		if s.isEdge(x, y) {
			y--
		}
	case LEFT:
		x--
		if s.isEdge(x, y) {
			x++
		}
	case RIGHT:
		x++
		if s.isEdge(x, y) {
			x--
		}
	}
	s.swap(s.blank, point{x, y})

}

func (s *Field) swap(a, b point) {
	tmp := s.board[s.getIdx(a.x, a.y)]
	s.board[s.getIdx(a.x, a.y)] = s.board[s.getIdx(b.x, b.y)]
	s.board[s.getIdx(b.x, b.y)] = tmp
	s.blank = point{b.x, b.y}
}

func (s Field) isEdge(x, y int) bool {
	return !(x >= 0 && x < s.size && y >= 0 && y < s.size)
}

func (s Field) getIdx(x, y int) int {
	return y*s.size + x
}

func (s Field) getPos(idx int) (int, int) {
	return idx % s.size, idx / s.size
}

func (s Field) String() string {
	str := ""
	for y := 0; y < s.size; y++ {
		for x := 0; x < s.size; x++ {
			str += s.board[s.getIdx(x, y)].String()
		}
		str += "\n"
	}
	return str
}
