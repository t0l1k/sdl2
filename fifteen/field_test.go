package fifteen

import (
	"fmt"
	"testing"
)

func TestCell(t *testing.T) {
	cell := NewCell(12)
	got := cell.IsBlank()
	want := false
	if got != want {
		t.Errorf("got:%v want:%v", got, want)
	}
	gotN := cell.GetNumber()
	wantN := 12
	if gotN != wantN {
		t.Errorf("got:%v want:%v", gotN, wantN)
	}
}

func TestField(t *testing.T) {
	t.Run("test init board", func(t *testing.T) {
		f := NewField(3)
		got := len(f.board)
		want := 9
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotBlank := f.blank
		wantBlank := point{2, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		gotX, gotY := f.getPos(2)
		wantX, wantY := 2, 0
		if gotX != wantX && gotY != wantY {
			t.Errorf("gotX:%v wantX:%v gotY:%v wantY:%v", gotX, wantX, gotY, wantY)
		}
	})

	t.Run("test move", func(t *testing.T) {
		f := NewField(3)
		fmt.Println(f)
		f.move(UP)
		gotBlank := f.blank
		wantBlank := point{2, 1}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(UP)
		gotBlank = f.blank
		wantBlank = point{2, 0}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(UP)
		gotBlank = f.blank
		wantBlank = point{2, 0}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(DOWN)
		gotBlank = f.blank
		wantBlank = point{2, 1}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(DOWN)
		gotBlank = f.blank
		wantBlank = point{2, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(DOWN)
		gotBlank = f.blank
		wantBlank = point{2, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(LEFT)
		gotBlank = f.blank
		wantBlank = point{1, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(LEFT)
		gotBlank = f.blank
		wantBlank = point{0, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(LEFT)
		gotBlank = f.blank
		wantBlank = point{0, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(RIGHT)
		gotBlank = f.blank
		wantBlank = point{1, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(RIGHT)
		gotBlank = f.blank
		wantBlank = point{2, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
		f.move(RIGHT)
		gotBlank = f.blank
		wantBlank = point{2, 2}
		if gotBlank != wantBlank {
			t.Errorf("got:%v want:%v", gotBlank, wantBlank)
		}
		fmt.Println(f)
	})
	t.Run("test Move", func(t *testing.T) {
		f := NewField(3)
		got := f.Move(2)
		want := true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB := f.board[2].GetNumber()
		wantB := 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(8)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[8].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(5)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[5].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(2)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[2].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(5)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[5].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(3)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[3].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(5)
		want = true
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[5].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(6)
		want = false
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[5].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)

		got = f.Move(1)
		want = false
		if got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		gotB = f.board[5].GetNumber()
		wantB = 0
		if gotB != wantB {
			t.Errorf("got:%v want:%v", gotB, wantB)
		}
		fmt.Println(f)
	})
}
func TestWin(t *testing.T) {
	t.Run("test win", func(t *testing.T) {
		f := NewField(3)
		f.Move(7)
		got := f.Win()
		want := false
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
		fmt.Println(f, f.Win())
		f.Move(8)
		got = f.Win()
		want = true
		fmt.Println(f)
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
		fmt.Println(f, f.Win())
	})

}
