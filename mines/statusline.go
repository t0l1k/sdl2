package mines

import (
	"sdl/mines/sdl2/ui"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type statusLine struct {
	renderer             *sdl.Renderer
	rect                 sdl.Rect
	fg, bg               sdl.Color
	font                 *ttf.Font
	lblFlags, lblStopper *ui.Label
	show                 bool
}

func newStatusLine(rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *statusLine {
	w, _, _ := font.SizeUTF8("00")
	x, y := rect.X+(rect.W/3)-int32(w), rect.Y
	lblFlags := ui.NewLabel("00", sdl.Point{x, y}, fg, renderer, font)
	w, _, _ = font.SizeUTF8("000")
	x, y = rect.X+(rect.W/3)*2-int32(w), rect.Y
	lblStopper := ui.NewLabel("000", sdl.Point{x, y}, fg, renderer, font)
	return &statusLine{
		rect:       rect,
		fg:         fg,
		bg:         bg,
		renderer:   renderer,
		font:       font,
		lblFlags:   lblFlags,
		lblStopper: lblStopper,
		show:       true,
	}
}

func (s *statusLine) setFlags(str string) {
	s.lblFlags.SetText(str)
}
func (s *statusLine) setStopper(str string) {
	s.lblStopper.SetText(str)
}

func (s *statusLine) Render(renderer *sdl.Renderer) {
	s.lblFlags.Render(renderer)
	s.lblStopper.Render(renderer)
}

func (s *statusLine) Update()         {}
func (s *statusLine) Event(sdl.Event) {}
func (s *statusLine) Destroy() {
	s.lblFlags.Destroy()
	s.lblStopper.Destroy()
}
