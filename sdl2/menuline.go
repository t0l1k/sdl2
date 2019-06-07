package sdl2

import (
	"fmt"

	"github.com/t0l1k/sdl2/clock"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MenuLine struct {
	renderer           *sdl.Renderer
	rect               sdl.Rect
	fg, bg             sdl.Color
	lblTitle, lblClock *Label
	btnQuit            *Button
	sprites            []Sprite
}

func NewMenuLine(title string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font, fn func()) *MenuLine {
	var sprites []Sprite
	lblTitle := NewLabel(title, sdl.Point{rect.H + 3, rect.Y}, fg, renderer, font)
	sprites = append(sprites, lblTitle)
	clockStr := "00:00:00"
	w, _, _ := font.SizeUTF8(clockStr)
	lblClock := NewLabel(clockStr, sdl.Point{rect.W - int32(w) - 3, rect.Y}, fg, renderer, font)
	sprites = append(sprites, lblClock)
	btnQuit := NewButton(renderer, "<-", sdl.Rect{rect.X, rect.Y, rect.H, rect.H}, fg, bg, font, fn)
	sprites = append(sprites, btnQuit)
	return &MenuLine{
		renderer: renderer,
		rect:     rect,
		fg:       fg,
		bg:       bg,
		lblTitle: lblTitle,
		lblClock: lblClock,
		btnQuit:  btnQuit,
		sprites:  sprites,
	}
}

func (s *MenuLine) Render(renderer *sdl.Renderer) {
	setColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	for _, sprite := range s.sprites {
		sprite.Render(s.renderer)
	}
}

func (s *MenuLine) Update() {
	_, sec, minute, hour := clock.Get()
	strClock := fmt.Sprintf("%02v:%02v:%02v", hour, minute, sec)
	s.lblClock.SetText(strClock)
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

func (s *MenuLine) Event(e sdl.Event) {
	for _, sprite := range s.sprites {
		sprite.Event(e)
	}
}

func (s *MenuLine) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
}
