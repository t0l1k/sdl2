package sdl2

import (
	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type StatusLine struct {
	rect    sdl.Rect
	fg, bg  sdl.Color
	sprites []ui.Sprite
}

func NewStatusLine(rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font, fnClock, fnLife, fnFifteen, fnMines func()) *StatusLine {
	var sprites []ui.Sprite
	btnClock := ui.NewButton(renderer, "Clock", sdl.Rect{rect.X, rect.Y, rect.H * 4, rect.H}, fg, bg, font, fnClock)
	sprites = append(sprites, btnClock)
	btnLife := ui.NewButton(renderer, "Conway's Life", sdl.Rect{rect.X + rect.H*4, rect.Y, rect.H * 5, rect.H}, fg, bg, font, fnLife)
	sprites = append(sprites, btnLife)
	btnFifteen := ui.NewButton(renderer, "Fifteen", sdl.Rect{rect.X + rect.H*9, rect.Y, rect.H * 4, rect.H}, fg, bg, font, fnFifteen)
	sprites = append(sprites, btnFifteen)
	btnMines := ui.NewButton(renderer, "Mines", sdl.Rect{rect.X + rect.H*13, rect.Y, rect.H * 4, rect.H}, fg, bg, font, fnMines)
	sprites = append(sprites, btnMines)
	return &StatusLine{
		rect:    rect,
		fg:      fg,
		bg:      bg,
		sprites: sprites,
	}
}

func (s *StatusLine) Render(renderer *sdl.Renderer) {
	ui.SetColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	for _, sprite := range s.sprites {
		sprite.Render(renderer)
	}
}

func (s *StatusLine) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

func (s *StatusLine) Event(e sdl.Event) {
	for _, sprite := range s.sprites {
		sprite.Event(e)
	}
}

func (s *StatusLine) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
}
