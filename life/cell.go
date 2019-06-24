package life

type cellStatus byte

const (
	blank cellStatus = iota
	alive
	born
	die
)

type cell struct {
	status cellStatus
}

func newCell() *cell {
	return &cell{status: blank}
}

func (s *cell) getStatus() cellStatus {
	return s.status
}

func (s *cell) setAlive() {
	s.status = alive
}

func (s *cell) setBlank() {
	s.status = blank
}

func (s *cell) setBorn() {
	s.status = born
}

func (s *cell) setDie() {
	s.status = die
}

func (s cell) String() string {
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
