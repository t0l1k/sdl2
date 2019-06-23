package mines

import (
	"fmt"
	"time"
)

type gameConfig struct {
	row, column, mines int
	begin, end         time.Time
	timer              time.Duration
}

func newGameConfig(row, column, mines int) *gameConfig {
	return &gameConfig{
		row:    row,
		column: column,
		mines:  mines,
		end:    time.Now(),
	}
}

func (s *gameConfig) getGameData() (int, int, int) {
	return s.row, s.column, s.mines
}

func (s *gameConfig) start() {
	s.begin = time.Now()
}

func (s *gameConfig) stop() {
	s.end = time.Now()
}

func (s *gameConfig) getStopper() int {
	if s.end.After(s.begin) {
		s.timer = s.end.Sub(s.begin)
		fmt.Println("set endper")
	} else {
		s.timer = time.Now().Sub(s.begin)
	}
	return int(s.timer.Seconds())
}
