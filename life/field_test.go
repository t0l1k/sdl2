package life

import (
	"fmt"
	"testing"
)

func TestField(t *testing.T) {

	t.Run("test turn ***", func(t *testing.T) {
		f := NewField(3)
		f.SetCellAlive(3)
		f.SetCellAlive(4)
		f.SetCellAlive(5)
		gotCount := f.CountNeighbours(1, 0)
		wantCount := 3
		if gotCount != wantCount {
			t.Errorf("got:%v want:%v", gotCount, wantCount)
		}
		fmt.Println(f)
		f.Turn()
		fmt.Println(f)
	})

	t.Run("test turn *****", func(t *testing.T) {
		f := NewField(11)
		f.SetCellAlive(f.GetIdx(3, 5))
		f.SetCellAlive(f.GetIdx(4, 5))
		f.SetCellAlive(f.GetIdx(5, 5))
		f.SetCellAlive(f.GetIdx(6, 5))
		f.SetCellAlive(f.GetIdx(7, 5))
		fmt.Println(f)
		for i := 0; i < 8; i++ {
			f.Turn()
			fmt.Println(f)
		}
	})
}
