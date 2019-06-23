package mines

import (
	"fmt"
	"testing"
)

func TestField(t *testing.T) {
	t.Run("test cells", func(t *testing.T) {
		f := NewField(3, 3, 3)
		f.field[0].setMine()
		f.field[0].open()
		f.field[0].setState(firstMined)
		f.field[1].setCounter(1)
		f.field[1].open()
		f.field[2].setCounter(0)
		f.field[2].markFlag()
		f.field[2].markFlag()
		f.field[2].markFlag()
		f.field[2].open()
		f.field[2].reset()
		f.field[2].open()
		fmt.Println(f)
	})
	t.Run("test cells", func(t *testing.T) {
		f := NewField(3, 3, 3)
		f.Shuffle(0, 0)
		open(f)
		fmt.Println(f)
	})
}

func open(f *Field) {
	for _, cell := range f.field {
		fmt.Printf("%v", cell)
		cell.open()
		fmt.Printf("%v,%v,%v", cell.state, cell.mined, cell.counter)
	}
	fmt.Println()
}
