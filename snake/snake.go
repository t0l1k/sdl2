package snake

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	left  = "left"
	right = "right"
	up    = "up"
	down  = "down"
)

type point struct {
	x, y float64
}

type snake struct {
	lives         int32
	dim           int32 // dimension
	cellSize      int32
	score         int
	head, lastPos point
	tail          []point
	direction     string
	step          float64
	renderer      *sdl.Renderer
	texture       *sdl.Texture
	fg, bg        sdl.Color
	rect          sdl.Rect
}

func newSnake(dim, cellSize int32, rect sdl.Rect, renderer *sdl.Renderer) *snake {
	fg := sdl.Color{0, 0, 255, 255}
	bg := sdl.Color{0, 0, 64, 255}
	head := point{4, 5}
	tail := []point{head, {3, 5}, {2, 5}}
	return &snake{
		lives:     3,
		dim:       dim,
		cellSize:  cellSize,
		head:      head,
		lastPos:   head,
		tail:      tail,
		direction: right,
		step:      0.05,
		renderer:  renderer,
		fg:        fg,
		bg:        bg,
		texture:   newSnakeTex(head, tail, renderer, sdl.Rect{0, 0, cellSize, cellSize}, fg, bg),
		rect:      rect,
	}
}

func (s *snake) reset() {
	if s.lives > 0 {
		s.head = point{4.95, 5}
		s.tail = []point{s.head, {3, 5}, {2, 5}}
		s.direction = right
		s.step = 0.05
		s.lives--
	} else {
		s.lives = 3
		s.head = point{4.95, 5}
		s.tail = []point{s.head, {3, 5}, {2, 5}}
		s.direction = right
		s.step = 0.05
		s.score = 0
	}
}

func (s *snake) checkBite() bool {
	for i, point := range s.tail {
		if i > 0 {
			if int(point.x) == int(s.head.x) && int(point.y) == int(s.head.y) {
				return true
			}
		}
	}
	return false
}

func (s *snake) grow() {
	s.tail = append(s.tail, s.tail[len(s.tail)-1])
	if len(s.tail)%10 == 0 {
		s.step += 0.01
	}
}

func (s *snake) move() {
	if s.direction == right {
		s.head.x += s.step
	} else if s.direction == left {
		s.head.x -= s.step
	} else if s.direction == up {
		s.head.y -= s.step
	} else if s.direction == down {
		s.head.y += s.step
	}
	if int32(s.head.x) != int32(s.lastPos.x) || int32(s.head.y) != int32(s.lastPos.y) {
		tail := make([]point, len(s.tail))
		copy(tail, s.tail)
		for i := 1; i < len(s.tail); i++ {
			s.tail[i] = tail[i-1]
		}
		s.tail[0] = s.head
	}
	s.lastPos = s.head
}

func (s *snake) moveLeft() {
	if s.direction == right {
		s.direction = up
	} else if s.direction == up {
		s.direction = left
	} else if s.direction == left {
		s.direction = down
	} else if s.direction == down {
		s.direction = right
	}
}

func (s *snake) moveRight() {
	if s.direction == right {
		s.direction = down
	} else if s.direction == down {
		s.direction = left
	} else if s.direction == left {
		s.direction = up
	} else if s.direction == up {
		s.direction = right
	}
}

func (s *snake) Render(renderer *sdl.Renderer) {
	renderer.Copy(s.texture, nil, &sdl.Rect{int32(s.head.x)*s.cellSize + s.rect.X, int32(s.head.y)*s.cellSize + s.rect.Y, s.cellSize, s.cellSize})
	for _, tail := range s.tail {
		renderer.Copy(s.texture, nil, &sdl.Rect{int32(tail.x)*s.cellSize + s.rect.X, int32(tail.y)*s.cellSize + s.rect.Y, s.cellSize, s.cellSize})
	}
}

func (s *snake) Update() {
	s.move()
}

func (s *snake) Event(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_LEFT && t.State == sdl.RELEASED {
			s.moveLeft()
		}
		if t.Keysym.Sym == sdl.K_RIGHT && t.State == sdl.RELEASED {
			s.moveRight()
		}
	}
}

func (s *snake) Destroy() { s.texture.Destroy() }

func newSnakeTex(head point, tail []point, renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) *sdl.Texture {
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
