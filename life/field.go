package life

import (
	"math/rand"
	"strconv"
)

type point struct {
	x, y int
}

type field struct {
	board []*cell
	dim   int
}

func newField(dim int) *field {
	var board []*cell
	for i := 0; i < dim*dim; i++ {
		cell := newCell()
		board = append(board, cell)
	}
	return &field{
		dim:   dim,
		board: board,
	}
}

func (s *field) shuffle() {
	for i := 0; i < int(float64(s.dim*s.dim)*0.6); i++ {
		for {
			x, y := rand.Intn(s.dim), rand.Intn(s.dim)
			if s.getCellStatus(s.getIdx(x, y)) == blank {
				s.setCellAlive(s.getIdx(x, y))
				break
			}
		}
	}
}

func (s *field) turn() {
	var board []cellStatus
	for idx := range s.board {
		pos := s.getPos(int(idx))
		count := s.countNeighbours(pos.x, pos.y)
		// fmt.Println(idx, cell, pos, count)
		if count == 3 && s.isCellBlank(idx) {
			board = append(board, born)
			// log.Println("ячека жива вокруг 3 соседа", pos, count, cell)
		} else if count < 2 && s.isCellAlive(idx) || count > 3 && s.isCellAlive(idx) {
			board = append(board, die)
			// log.Println("вокруг ячейки меньше двух соседей, умираем от одиночества или больше трех соседей, умираем от тесноты", pos, count, cell)
		} else if s.isCellAlive(idx) {
			board = append(board, alive)
		} else {
			board = append(board, blank)
		}
	}
	for idx := range board {
		if board[idx] == born {
			s.board[idx].setAlive()
		} else if board[idx] == die {
			s.board[idx].setBlank()
		} else if board[idx] == alive {
			s.board[idx].setAlive()
		}
	}
}

func (s *field) getBoard() []*cell {
	return s.board
}

func (s *field) getCellStatus(idx int) cellStatus {
	return s.board[idx].getStatus()
}

func (s *field) isCellBlank(idx int) bool {
	return s.board[idx].getStatus() == blank
}

func (s *field) isCellAlive(idx int) bool {
	return s.board[idx].getStatus() == alive
}

func (s *field) isCellBorn(idx int) bool {
	return s.board[idx].getStatus() == born
}

func (s *field) isCellDie(idx int) bool {
	return s.board[idx].getStatus() == die
}

func (s *field) setCellAlive(idx int) {
	s.board[idx].setAlive()
}

func (s *field) setCellBlank(idx int) {
	s.board[idx].setBlank()
}

func (s *field) setCellBorn(idx int) {
	s.board[idx].setBorn()
}

func (s *field) setCellDie(idx int) {
	s.board[idx].setDie()
}

func (s *field) isEdge(x, y int) bool {
	return !(x >= 0 && x < s.dim && y >= 0 && y < s.dim)
}

func (s *field) countNeighbours(x0, y0 int) (count int) {
	arr := s.getNeighbour(x0, y0)
	for _, cell := range arr {
		if cell.getStatus() == alive {
			count++
		}
	}
	return count
}

func (s *field) getNeighbour(x0, y0 int) (neighbour []*cell) {
	if !s.isEdge(x0, y0) {
		for y := int(-1); y < 2; y++ {
			for x := int(-1); x < 2; x++ {
				var nx, ny int
				nx = x0 + x
				ny = y0 + y
				if nx >= 0 && nx < s.dim && ny >= 0 && ny < s.dim && !(x == 0 && y == 0) {
					nCell := s.board[s.getIdx(nx, ny)]
					neighbour = append(neighbour, nCell)
				}
			}
		}
	}
	return neighbour
}

func (s *field) getPos(idx int) point {
	x, y := idx%s.dim, idx/s.dim
	return point{x, y}
}

func (s *field) getIdx(x, y int) int {
	return y*s.dim + x
}

func (s *field) String() string {
	str := "field:" + strconv.Itoa(int(s.dim)) + "\n"
	for i, cell := range s.board {
		if i%int(s.dim) == 0 && i > 0 {
			str += "|\n-"
			for i := int(0); i < s.dim; i++ {
				str += "--"
			}
			str += "\n"
		}
		str += "|"
		str += cell.String()
	}
	str += "|\n"
	return str
}
