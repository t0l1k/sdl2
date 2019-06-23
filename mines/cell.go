package mines

import "strconv"

type cellState string

const (
	closed       cellState = "closed"        // начальное состояние
	flagged                = "flagged"       // отмечена флагом
	questioned             = "questioned"    // отмечена вопросом
	opened                 = "opened"        // открыта
	empty                  = "empty"         // помечена цифрой
	firstMined             = "first mined"   // взорваная первой
	saved                  = "saved"         // если помечена флагом
	blown                  = "blown"         // есть мина, не отмечена флагом
	wrongFlagged           = "wrong flagged" // помечена флагом, без мины
	cellEmpty              = 9
)

type cell struct {
	state   cellState
	counter byte
	mined   bool
	x, y    int
}

func newCell(x, y int) *cell {
	return &cell{
		state:   closed,
		counter: cellEmpty,
		x:       x,
		y:       y,
	}
}

func (s *cell) getX() int {
	return s.x
}

func (s *cell) getY() int {
	return s.y
}

func (s *cell) getState() cellState {
	return s.state
}

func (s *cell) setState(state cellState) {
	s.state = state
}

func (s *cell) getCount() byte {
	return s.counter
}

func (s *cell) setCounter(count byte) {
	s.counter = count
}

func (s *cell) getMined() bool {
	return s.mined && s.counter == cellEmpty
}

func (s *cell) setMine() {
	s.mined = true
}

func (s *cell) open() {
	if s.state == closed || s.state == questioned {
		s.state = opened
	}
}

func (s *cell) markFlag() {
	switch s.state {
	case closed:
		s.state = flagged
	case flagged:
		s.state = questioned
	case questioned:
		s.state = closed
	}
}

func (s *cell) reset() {
	s.state = closed
}

func (s cell) String() (str string) {
	switch s.state {
	case closed:
		str += " "
	case flagged:
		str += "F"
	case questioned:
		str += "Q"
	case firstMined:
		str += "f"
	case saved:
		str += "v"
	case blown:
		str += "b"
	case wrongFlagged:
		str += strconv.Itoa(int(s.counter))
	case opened:
		if !s.getMined() && s.getCount() != cellEmpty {
			switch s.getCount() {
			case 0:
				str += " "
			default:
				str += strconv.Itoa(int(s.counter))
			}
		} else if s.getMined() && s.getCount() == cellEmpty {
			str += "*"
		} else {
			str += "!"
		}
	default:
		str += "?"
	}
	return str
}
