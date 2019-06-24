package life

import (
	"math/rand"
	"strconv"
)

type Point struct {
	X, Y int
}

type Field struct {
	board []*Cell
	dim   int
}

func NewField(dim int) *Field {
	var board []*Cell
	for i := 0; i < dim*dim; i++ {
		cell := NewCell()
		board = append(board, cell)
	}
	return &Field{
		dim:   dim,
		board: board,
	}
}

func (s *Field) Shuffle() {
	for i := 0; i < int(float64(s.dim*s.dim)*0.6); i++ {
		for {
			x, y := rand.Intn(s.dim), rand.Intn(s.dim)
			if s.GetCellStatus(s.GetIdx(x, y)) == blank {
				s.SetCellAlive(s.GetIdx(x, y))
				break
			}
		}
	}
}

func (s *Field) Turn() {
	var board []CellStatus
	for idx := range s.board {
		pos := s.GetPos(int(idx))
		count := s.CountNeighbours(pos.X, pos.Y)
		// fmt.Println(idx, cell, pos, count)
		if count == 3 && s.IsCellBlank(idx) {
			board = append(board, born)
			// log.Println("ячека жива вокруг 3 соседа", pos, count, cell)
		} else if count < 2 && s.IsCellAlive(idx) || count > 3 && s.IsCellAlive(idx) {
			board = append(board, die)
			// log.Println("вокруг ячейки меньше двух соседей, умираем от одиночества или больше трех соседей, умираем от тесноты", pos, count, cell)
		} else if s.IsCellAlive(idx) {
			board = append(board, alive)
		} else {
			board = append(board, blank)
		}
	}
	for idx := range board {
		if board[idx] == born {
			s.board[idx].SetAlive()
		} else if board[idx] == die {
			s.board[idx].SetBlank()
		} else if board[idx] == alive {
			s.board[idx].SetAlive()
		}
	}
}

func (s *Field) GetBoard() []*Cell {
	return s.board
}

func (s *Field) GetCellStatus(idx int) CellStatus {
	return s.board[idx].GetStatus()
}

func (s *Field) IsCellBlank(idx int) bool {
	return s.board[idx].GetStatus() == blank
}

func (s *Field) IsCellAlive(idx int) bool {
	return s.board[idx].GetStatus() == alive
}

func (s *Field) IsCellBorn(idx int) bool {
	return s.board[idx].GetStatus() == born
}

func (s *Field) IsCellDie(idx int) bool {
	return s.board[idx].GetStatus() == die
}

func (s *Field) SetCellAlive(idx int) {
	s.board[idx].SetAlive()
}

func (s *Field) SetCellBlank(idx int) {
	s.board[idx].SetBlank()
}

func (s *Field) SetCellBorn(idx int) {
	s.board[idx].SetBorn()
}

func (s *Field) SetCellDie(idx int) {
	s.board[idx].SetDie()
}

func (s *Field) IsEdge(x, y int) bool {
	return !(x >= 0 && x < s.dim && y >= 0 && y < s.dim)
}

func (s *Field) CountNeighbours(x0, y0 int) (count int) {
	arr := s.GetNeighbour(x0, y0)
	for _, cell := range arr {
		if cell.GetStatus() == alive {
			count++
		}
	}
	return count
}

func (s *Field) GetNeighbour(x0, y0 int) (neighbour []*Cell) {
	if !s.IsEdge(x0, y0) {
		for y := int(-1); y < 2; y++ {
			for x := int(-1); x < 2; x++ {
				var nx, ny int
				nx = x0 + x
				ny = y0 + y
				if nx >= 0 && nx < s.dim && ny >= 0 && ny < s.dim && !(x == 0 && y == 0) {
					nCell := s.board[s.GetIdx(nx, ny)]
					neighbour = append(neighbour, nCell)
				}
			}
		}
	}
	return neighbour
}

func (s *Field) GetPos(idx int) Point {
	x, y := idx%s.dim, idx/s.dim
	return Point{x, y}
}

func (s *Field) GetIdx(x, y int) int {
	return y*s.dim + x
}

func (s *Field) String() string {
	str := "Field:" + strconv.Itoa(int(s.dim)) + "\n"
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
