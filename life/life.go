package life

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Life struct {
	dim             int32
	renderer        *sdl.Renderer
	rect            sdl.Rect
	bg, fg          sdl.Color
	texBG, texFG    *sdl.Texture
	field           *Field
	DELAY, lastTime uint32
	show            bool
}

func NewLife(dim int32, delay uint32, renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) *Life {
	field := NewField(int(dim))
	field.Shuffle()
	texBG := NewBackground(renderer, rect, fg, bg, dim)
	texFG := NewForeground(renderer, rect, fg, bg, dim, field)
	return &Life{
		dim:      dim,
		renderer: renderer,
		rect:     rect,
		bg:       bg,
		fg:       fg,
		texBG:    texBG,
		texFG:    texFG,
		field:    field,
		DELAY:    delay,
		show:     false,
	}
}

func NewBackground(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color, dim int32) *sdl.Texture {
	texBackground, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texBackground)
	texBackground.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, sdl.Color{255, 0, 0, 255})
	renderer.DrawRect(&sdl.Rect{0, 0, rect.W, rect.H})
	setColor(renderer, fg)
	cellWidth, cellHeight := float64(rect.W)/float64(dim), float64(rect.H)/float64(dim)
	x0 := (float64(rect.W) - cellWidth*float64(dim)) / 2
	y0 := (float64(rect.H) - cellHeight*float64(dim)) / 2
	var x, y, w, h, dx, dy int32
	for dy = 0; dy < dim; dy++ {
		for dx = 0; dx < dim; dx++ {
			x = int32(x0 + float64(dx)*cellWidth)
			y = int32(y0 + float64(dy)*cellHeight)
			w = int32(float64(rect.W) - x0)
			h = int32(float64(rect.H) - y0)
			renderer.DrawLine(x, int32(y0), x, h)
			renderer.DrawLine(int32(x0), y, w, y)
		}
	}
	renderer.SetRenderTarget(nil)
	return texBackground
}

func NewForeground(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color, dim int32, f *Field) *sdl.Texture {
	texForeground, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texForeground)
	texForeground.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)

	cellWidth, cellHeight := float64(rect.W)/float64(dim), float64(rect.H)/float64(dim)
	x0 := (float64(rect.W) - cellWidth*float64(dim)) / 2
	y0 := (float64(rect.H) - cellHeight*float64(dim)) / 2
	w, h := float64(cellWidth)*0.5, float64(cellHeight)*0.5
	marginX, marginY := (cellWidth-w)/2, (cellHeight-h)/2
	for idx, cell := range f.GetBoard() {
		if cell.GetStatus() == alive {
			fPos := f.GetPos(idx)
			x, y := float64(fPos.X)*cellWidth+x0+marginX, float64(fPos.Y)*cellHeight+y0+marginY
			renderer.FillRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
		}
	}

	renderer.SetRenderTarget(nil)
	return texForeground
}

func (s *Life) GetShow() bool {
	return s.show
}

func (s *Life) SetShow(show bool) {
	s.show = show
}

func (s *Life) Event(e sdl.Event) {
	if s.show {
		switch t := e.(type) {
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_RETURN && t.State == sdl.RELEASED {
				s.field.Shuffle()
			}
		}
	}
}

func (s *Life) Update() {
	if s.show {
		if sdl.GetTicks()-s.lastTime > s.DELAY {
			s.field.Turn()
			s.texFG.Destroy()
			s.texFG = NewForeground(s.renderer, s.rect, s.fg, s.bg, s.dim, s.field)
			s.lastTime = sdl.GetTicks()
		}
	}
}

func (s *Life) Render(renderer *sdl.Renderer) {
	if s.show {
		renderer.Copy(s.texBG, nil, &s.rect)
		renderer.Copy(s.texFG, nil, &s.rect)
	}
}

func (s *Life) Destroy() {
	s.texBG.Destroy()
	s.texFG.Destroy()
}

func setColor(renderer *sdl.Renderer, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

func FillCircle(renderer *sdl.Renderer, x0, y0, radius int32, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(x0+x+radius, y0+y+radius)
			}
		}
	}
}
