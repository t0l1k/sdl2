package fifteen

import (
	"strconv"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {
	renderer   *sdl.Renderer
	rect       sdl.Rect
	dim        int
	sprites    []ui.Sprite
	fg, bg     sdl.Color
	field      *field
	messageBox *ui.MessageBox
	show       bool
}

func NewGame(renderer *sdl.Renderer, rect sdl.Rect) *Game {
	dim := 4
	fg, bg := sdl.Color{255, 255, 0, 255}, sdl.Color{0, 72, 0, 0}
	return &Game{
		dim:      dim,
		renderer: renderer,
		rect:     rect,
		fg:       fg,
		bg:       bg,
	}
}

func (s *Game) Setup() {
	var (
		cellSize int32
		fnSize   int
	)
	if s.rect.W < s.rect.H {
		cellSize = s.rect.W / int32(s.dim)
		fnSize = fontSize(s.rect.H)
	} else {
		cellSize = s.rect.H / int32(s.dim)
		fnSize = fontSize(s.rect.W)
	}
	font, err := ttf.OpenFont("assets/Roboto-Regular.ttf", fnSize)
	if err != nil {
		panic(err)
	}
	s.field = newField(s.dim)
	marginX := (s.rect.W - cellSize*int32(s.dim)) / 2
	marginY := (s.rect.H - cellSize*int32(s.dim)) / 2
	for i, value := range s.field.getBoard() {
		var x, y int32
		x = s.rect.X + cellSize*int32(i%s.dim) + marginX
		y = s.rect.Y + cellSize*int32(i/s.dim) + marginY
		if value.getNumber() > 0 {
			bnt := ui.NewButton(s.renderer, strconv.Itoa(value.getNumber()), sdl.Rect{x, y, cellSize, cellSize}, s.fg, s.bg, font, s.checkBtn)
			s.sprites = append(s.sprites, bnt)
		} else {
			bnt := ui.NewButton(s.renderer, " ", sdl.Rect{x, y, cellSize, cellSize}, s.bg, s.bg, font, s.checkBtn)
			s.sprites = append(s.sprites, bnt)
		}
	}
	s.messageBox = ui.NewMessageBox("Fifteen", "Start New Game", s.rect, s.fg, s.bg, s.renderer, s.reset)
	s.messageBox.SetShow(true)
}

func (s *Game) checkBtn() {
	for i, sprite := range s.sprites {
		if sprite.(*ui.Button).IsPressed() {
			s.field.moves(i)
			for i, value := range s.field.getBoard() {
				switch s.sprites[i].(type) {
				case *ui.Button:
					if value.getNumber() > 0 {
						s.sprites[i].(*ui.Button).SetFGBG(s.fg, s.bg)
						s.sprites[i].(*ui.Button).SetText(strconv.Itoa(value.getNumber()))
					} else {
						s.sprites[i].(*ui.Button).SetFGBG(s.bg, s.bg)
						s.sprites[i].(*ui.Button).SetText(" ")
					}
				}
			}
			if s.field.win() {
				s.messageBox.SetMessage("Board Solved")
				s.messageBox.SetShow(true)
			}
		}
	}
}

func (s *Game) reset() {
	s.field.shuffle(100)
	for i, value := range s.field.getBoard() {
		switch s.sprites[i].(type) {
		case *ui.Button:
			if value.getNumber() > 0 {
				s.sprites[i].(*ui.Button).SetFGBG(s.fg, s.bg)
				s.sprites[i].(*ui.Button).SetText(strconv.Itoa(value.getNumber()))
			} else {
				s.sprites[i].(*ui.Button).SetFGBG(s.bg, s.bg)
				s.sprites[i].(*ui.Button).SetText(" ")
			}
		}
	}
	s.messageBox.SetShow(false)
}

func (s *Game) GetShow() bool {
	return s.show
}

func (s *Game) SetShow(show bool) {
	s.show = show
}

func (s *Game) Render(renderer *sdl.Renderer) {
	if s.show {
		for _, sprite := range s.sprites {
			sprite.Render(s.renderer)
		}
		s.messageBox.Render(renderer)
	}
}

func (s *Game) Update() {
	if s.show {
		for _, sprite := range s.sprites {
			sprite.Update()
		}
		s.messageBox.Update()
	}
}
func (s *Game) Event(event sdl.Event) {
	if s.show {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_RETURN && t.State == sdl.RELEASED {
				s.reset()
			}
		}
		for _, sprite := range s.sprites {
			sprite.Event(event)
		}
		s.messageBox.Event(event)
	}
}
func (s *Game) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
}

func fontSize(height int32) int {
	return int(float64(height) * 0.08) // Главная константа перерисовки экрана размер шрифта
}

func lineHeight(fontSize int) int32 {
	return int32(float64(fontSize) * 1.5) // Высота строки меню и строки статуса
}
