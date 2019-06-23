package mines

import (
	"testing"
)

func TestCell(t *testing.T) {
	t.Run("test init", func(t *testing.T) {
		c := newCell(0, 0)
		got := c.getState()
		want := closed
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
	})
	t.Run("test mine", func(t *testing.T) {
		c := newCell(0, 0)
		c.setMine()
		got := c.getMined()
		want := true
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
		c.open()
	})
	t.Run("test counter", func(t *testing.T) {
		c := newCell(0, 0)
		c.setCounter(1)
		got := c.getCount()
		want := byte(1)
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
		c.open()
	})
}
