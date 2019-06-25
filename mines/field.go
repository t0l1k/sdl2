package mines

import (
	"math/rand"
	"strconv"
	"time"
)

type point struct {
	x, y int
}

type gameState string

const (
	gameStart gameState = "start"
	gamePlay  gameState = "play"
	gamePause gameState = "pause"
	gameWin   gameState = "win"
	gameOver  gameState = "game over"
)

type minesField struct {
	field              []*cell
	savedField         []cellState
	state              gameState
	row, column, mines int
	firstMove          point
}

func newMinesField(row, column, mines int) *minesField {
	var field []*cell
	for y := 0; y < column; y++ {
		for x := 0; x < row; x++ {
			cell := newCell(x, y)
			field = append(field, cell)
		}
	}
	return &minesField{
		row:    row,
		column: column,
		mines:  mines,
		field:  field,
		state:  gameStart,
	}
}

func (s *minesField) getBoard() []*cell {
	return s.field
}

func (s *minesField) open(x, y int) {
	if s.isFieldEdge(x, y) {
		return
	}
	cell := s.getCell(x, y)
	if cell.getState() == flagged || cell.getState() == opened {
		return
	}
	cell.open()
	if cell.getMined() {
		cell.setState(firstMined)
		s.state = gameOver
		return
	}
	if cell.getCount() > 0 {
		return
	}
	for _, newCell := range s.getNeighbours(x, y) {
		s.open(newCell.getX(), newCell.getY())
	}
}

func (s *minesField) autoMarkFlags(x, y int) {
	var countClosed, countFlags byte
	cellValue := s.getCell(x, y).getCount()
	neighbours := s.getNeighbours(x, y)
	if s.getCell(x, y).getState() == opened {
		for _, value := range neighbours {
			if value.getState() == flagged {
				countFlags++
			} else if value.getState() == closed || value.getState() == questioned {
				countClosed++
			}
		}
	}
	if countClosed+countFlags == cellValue {
		for _, value := range neighbours {
			if value.getState() == closed || value.getState() == questioned {
				value.setState(flagged)
			}
		}
	} else if countFlags == cellValue {
		for _, value := range neighbours {
			s.open(value.getX(), value.getY())
		}
	}
}

func (s *minesField) markFlag(x, y int) {
	s.field[s.getIdx(x, y)].markFlag()
}

func (s *minesField) shuffle(fX, fY int) {
	if s.state == gameStart {
		s.firstMove.x = fX
		s.firstMove.y = fY
		rand.Seed(time.Now().UTC().UnixNano())
		var mines int
		for mines < s.mines {
			x, y := rand.Intn(s.row), rand.Intn(s.column)
			if x != fX && y != fY {
				cell := s.field[s.getIdx(x, y)]
				if !cell.getMined() {
					cell.setMine()
					mines++
				}
			}
		}
		for idx, cell := range s.field {
			var count byte
			if !cell.getMined() {
				x, y := s.getPos(idx)
				neighbours := s.getNeighbours(x, y)
				for _, newCell := range neighbours {
					if newCell.getMined() {
						count++
					}
				}
				s.field[idx].setCounter(count)
			}
		}
	}
	s.state = gamePlay
}

func (s *minesField) isWin() bool {
	var count int
	for _, cell := range s.field {
		if cell.getState() == opened {
			count++
		}
	}
	if count+s.mines == s.row*s.column {
		s.state = gameWin
		for idx, cell := range s.field {
			if cell.getMined() {
				s.field[idx].setState(saved)
			}
		}
		return true
	}
	return false
}

func (s *minesField) isGameOver() bool {
	if s.state == gameOver {
		for idx, cell := range s.field {
			if cell.getMined() && cell.getState() == closed {
				s.field[idx].open()
				s.field[idx].setState(blown)
			} else if cell.getState() == flagged && cell.getMined() {
				s.field[idx].setState(saved)
			}
		}
	} else {
		return false
	}
	return true
}

func (s *minesField) saveLastMove() {
	s.savedField = s.savedField[:0]
	for idx := range s.field {
		s.savedField = append(s.savedField, s.field[idx].getState())
	}
}

func (s *minesField) restoreLastMove() {
	for idx := range s.field {
		s.field[idx].setState(s.savedField[idx])
	}
	if s.state == gameOver || s.state == gameWin {
		s.state = gamePlay
	}
}

func (s *minesField) reset() {
	for idx := range s.field {
		s.field[idx].reset()
	}
	s.state = gameStart
	s.open(s.firstMove.x, s.firstMove.y)
}

func (s *minesField) isFieldEdge(x, y int) bool {
	return x < 0 || x > s.row-1 || y < 0 || y > s.column-1
}

func (s *minesField) getIdx(x, y int) int {
	return y*s.row + x
}

func (s *minesField) getPos(idx int) (int, int) {
	return idx % s.row, idx / s.row // x, y
}

func (s *minesField) getCell(x, y int) *cell {
	return s.field[s.getIdx(x, y)]
}

func (s *minesField) getNeighbours(x, y int) (cells []*cell) {
	for dy := -1; dy < 2; dy++ {
		for dx := -1; dx < 2; dx++ {
			nx := x + dx
			ny := y + dy
			if !s.isFieldEdge(nx, ny) {
				newCell := s.getCell(nx, ny)
				cells = append(cells, newCell)
			}
		}
	}
	return cells
}

func (s *minesField) getState() gameState {
	return s.state
}

func (s *minesField) setState(state gameState) {
	s.state = state
}

func (s *minesField) getStatus() (flags int) {
	for _, cell := range s.field {
		if cell.getState() == flagged || cell.getState() == saved {
			flags++
		}
	}
	return s.mines - flags
}

func (s *minesField) destroy() {
	s.field = nil
}

func (s *minesField) String() string {
	str := strconv.Itoa(s.row) + ","
	str += strconv.Itoa(s.column) + ","
	str += strconv.Itoa(s.mines) + "\n"
	for y := 0; y < s.column; y++ {
		for x := 0; x < s.row; x++ {
			str += s.field[y*s.row+x].String()
		}
		str += "\n"
	}
	return str
}
