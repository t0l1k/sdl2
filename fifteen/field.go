package fifteen

import (
	"math/rand"
)

const (
	up = iota
	down
	left
	right
)

type point struct {
	x, y int
}

type field struct {
	size  int
	board []*cell
	blank point
}

func newField(size int) *field {
	var (
		board []*cell
		blank point
	)
	for i := 0; i < size*size; i++ {
		if i >= size*size-1 {
			board = append(board, newCell(0))
			blank = point{i % size, i / size}
		} else {
			board = append(board, newCell(i+1))
		}
	}
	return &field{
		size:  size,
		board: board,
		blank: blank,
	}
}

func (s *field) getBoard() []*cell {
	return s.board
}

func (s *field) win() bool {
	i := 0
	for _, cell := range s.board {
		if cell.getNumber() == i+1 {
			i++
		} else if cell.getNumber() != i+1 && cell.getNumber() != 0 {
			return false
		}
	}
	if s.board[i].getNumber() == 0 && i == len(s.board)-1 {
		return true
	}
	return false
}

func (s *field) shuffle(count int) {
	for i := 0; i < count; i++ {
		direction := rand.Intn(4)
		s.move(direction)
	}
}

func (s *field) moves(idx int) bool {
	x, y := s.getPos(idx)
	if s.blank.x == x {
		if s.blank.y > y {
			i := s.blank.y
			for s.blank.y <= i && y != i {
				s.move(up)
				i--
			}
		} else if s.blank.y < y {
			i := s.blank.y
			for s.blank.y >= i && y != i {
				s.move(down)
				i++
			}
		}
		return true
	} else if s.blank.y == y {
		if s.blank.x > x {
			i := s.blank.x
			for s.blank.x <= i && x != i {
				s.move(left)
				i--
			}
		} else if s.blank.x < x {
			i := s.blank.x
			for s.blank.x >= i && x != i {
				s.move(right)
				i++
			}
		}
		return true
	}
	return false
}

func (s *field) move(direction int) {
	x, y := s.blank.x, s.blank.y
	switch direction {
	case up:
		y--
		if s.isEdge(x, y) {
			y++
		}
	case down:
		y++
		if s.isEdge(x, y) {
			y--
		}
	case left:
		x--
		if s.isEdge(x, y) {
			x++
		}
	case right:
		x++
		if s.isEdge(x, y) {
			x--
		}
	}
	s.swap(s.blank, point{x, y})

}

func (s *field) swap(a, b point) {
	tmp := s.board[s.getIdx(a.x, a.y)]
	s.board[s.getIdx(a.x, a.y)] = s.board[s.getIdx(b.x, b.y)]
	s.board[s.getIdx(b.x, b.y)] = tmp
	s.blank = point{b.x, b.y}
}

func (s field) isEdge(x, y int) bool {
	return !(x >= 0 && x < s.size && y >= 0 && y < s.size)
}

func (s field) getIdx(x, y int) int {
	return y*s.size + x
}

func (s field) getPos(idx int) (int, int) {
	return idx % s.size, idx / s.size
}

func (s field) String() string {
	str := ""
	for y := 0; y < s.size; y++ {
		for x := 0; x < s.size; x++ {
			str += s.board[s.getIdx(x, y)].String()
		}
		str += "\n"
	}
	return str
}
