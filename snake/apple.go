package snake

import (
	"math/rand"

	"github.com/t0l1k/sdl2/sdl2/ui"

	"github.com/veandco/go-sdl2/sdl"
)

type apple struct {
	dim      int32 // dimension
	cellSize int32
	renderer *sdl.Renderer
	rect     sdl.Rect
	texture  *sdl.Texture
	fg, bg   sdl.Color
	sprites  []ui.Sprite
}

func newApple(dim int32, cellSize int32, rect sdl.Rect, renderer *sdl.Renderer) *apple {
	fg := sdl.Color{255, 0, 0, 255}
	bg := sdl.Color{0, 0, 64, 255}
	return &apple{
		dim:      dim,
		cellSize: cellSize,
		renderer: renderer,
		rect: sdl.Rect{
			int32(rand.Intn(int(dim)))*cellSize + rect.X,
			int32(rand.Intn(int(dim)))*cellSize + rect.Y,
			cellSize,
			cellSize,
		},
		fg:      fg,
		bg:      bg,
		texture: newAppleTex(renderer, sdl.Rect{0, 0, cellSize, cellSize}, fg, bg),
	}
}

func (s *apple) next(x, y int32) {
	s.rect.X = x
	s.rect.Y = y
}

func (s *apple) Render(renderer *sdl.Renderer) {
	if err := renderer.Copy(s.texture, nil, &s.rect); err != nil {
		panic(err)
	}
	for _, sprite := range s.sprites {
		sprite.Render(s.renderer)
	}
}

func (s *apple) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

func (s *apple) Event(event sdl.Event) {
	for _, sprite := range s.sprites {
		sprite.Event(event)
	}
}

func (s *apple) Destroy() {
	s.texture.Destroy()
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
}

func newAppleTex(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) *sdl.Texture {
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texture)
	renderer.SetDrawColor(bg.R, bg.G, bg.B, bg.A)
	renderer.Clear()
	renderer.SetDrawColor(fg.R, fg.G, fg.B, fg.A)
	renderer.DrawLine(0, 0, rect.W, rect.H)
	renderer.DrawLine(rect.W, 0, 0, rect.H)
	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.DrawRect(&sdl.Rect{0, 0, rect.W, rect.H})
	renderer.SetRenderTarget(nil)
	return texture
}
