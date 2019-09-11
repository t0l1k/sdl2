package snake

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/ttf"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	dim, cellSize      int32
	renderer           *sdl.Renderer
	rect               sdl.Rect
	texture            *sdl.Texture
	fg, bg             sdl.Color
	snake              *snake
	apple              *apple
	lblLives, lblScore *ui.Label
	btnLeft, btnRight  *ui.Button
	sprites            []ui.Sprite
	font               *ttf.Font
	show               bool
}

func NewGame(renderer *sdl.Renderer, r sdl.Rect) *Game {
	dim := int32(15)
	rect := r
	fg := sdl.Color{0, 128, 0, 255}
	bg := sdl.Color{0, 64, 0, 255}
	texture, cellSize := newBackground(dim, renderer, rect, fg, bg)
	return &Game{
		dim:      dim,
		renderer: renderer,
		rect:     rect,
		fg:       fg,
		bg:       bg,
		texture:  texture,
		cellSize: cellSize,
		show:     false,
	}
}

func (s *Game) Setup() {
	var (
		fontName = "assets/Roboto-Regular.ttf"
		err      error
	)
	s.font, err = ttf.OpenFont(fontName, int(float64(s.cellSize)*0.8))
	if err != nil {
		panic(err)
	}
	s.apple = newApple(s.dim, s.cellSize, s.rect, s.renderer)
	s.sprites = append(s.sprites, s.apple)
	s.snake = newSnake(s.dim, s.cellSize, s.rect, s.renderer)
	s.sprites = append(s.sprites, s.snake)
	s.lblLives = ui.NewLabel(" ", sdl.Point{s.rect.X + s.cellSize*2, s.rect.Y + s.cellSize}, s.fg, s.renderer, s.font)
	s.sprites = append(s.sprites, s.lblLives)
	s.lblScore = ui.NewLabel(" ", sdl.Point{s.dim*s.cellSize - s.cellSize*2, s.rect.Y + s.cellSize}, s.fg, s.renderer, s.font)
	s.sprites = append(s.sprites, s.lblScore)
	s.btnLeft = ui.NewButton(s.renderer, "<<<", sdl.Rect{
		s.rect.X,
		s.cellSize * (s.dim + 1),
		(s.cellSize * s.dim) / 2,
		s.cellSize,
	}, s.fg, s.bg, s.font, s.snake.moveLeft)
	s.sprites = append(s.sprites, s.btnLeft)
	s.btnRight = ui.NewButton(s.renderer, ">>>", sdl.Rect{
		s.rect.X + (s.cellSize*s.dim)/2,
		s.cellSize * (s.dim + 1),
		(s.cellSize * s.dim) / 2,
		s.cellSize,
	}, s.fg, s.bg, s.font, s.snake.moveRight)
	s.sprites = append(s.sprites, s.btnRight)
}

func (s *Game) Render(renderer *sdl.Renderer) {
	if s.show {
		if err := renderer.Copy(s.texture, nil, &sdl.Rect{s.rect.X, s.rect.Y, s.cellSize * s.dim, s.cellSize * s.dim}); err != nil {
			panic(err)
		}
		for _, sprite := range s.sprites {
			sprite.Render(s.renderer)
		}
	}
}

func (s *Game) SetShow(show bool) {
	s.show = show
}

func (s *Game) nextApple() (x, y int32) {
	rand.Seed(time.Now().UTC().UnixNano())
	for {
	again:
		x = int32(rand.Intn(int(s.dim)))*s.cellSize + s.rect.X
		y = int32(rand.Intn(int(s.dim)))*s.cellSize + s.rect.Y
		for _, point := range s.snake.tail {
			pX := int32(point.x)*s.cellSize + s.rect.X
			pY := int32(point.y)*s.cellSize + s.rect.Y
			if pX == x && pY == y {
				goto again
			}
		}
		break
	}
	return x, y
}

func (s *Game) checkSnakeOutOfBoard() bool {
	if s.snake.head.x <= 0 || s.snake.head.y <= 0 || int32(s.snake.head.x)*s.cellSize >= s.dim*s.cellSize || int32(s.snake.head.y)*s.cellSize >= s.dim*s.cellSize {
		return true
	}
	return false
}

func (s *Game) checkEatApple() {
	if s.apple.rect.X == int32(s.snake.head.x)*s.snake.cellSize+s.rect.X && s.apple.rect.Y == int32(s.snake.head.y)*s.snake.cellSize+s.rect.Y {
		s.snake.grow()
		s.apple.next(s.nextApple())
		s.snake.score++
	}
}

func (s *Game) showScoreAndLives() {
	var strLives string
	for i := 0; i < int(s.snake.lives); i++ {
		strLives += "@"
	}
	if len(strLives) == 0 {
		strLives = " "
	}
	s.lblLives.SetText(strLives)
	s.lblScore.SetText(fmt.Sprintf("%.3v", s.snake.score))
}

func (s *Game) Update() {
	if s.show {
		for _, sprite := range s.sprites {
			sprite.Update()
		}
		if s.snake.checkBite() {
			s.snake.reset()
		} else if s.checkSnakeOutOfBoard() {
			s.snake.reset()
		}
		s.checkEatApple()
		s.showScoreAndLives()
	}
}

func (s *Game) Event(event sdl.Event) {
	if s.show {
		for _, sprite := range s.sprites {
			sprite.Event(event)
		}
	}
}

func (s *Game) Destroy() {
	s.texture.Destroy()
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
}

func newBackground(dim int32, renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texture *sdl.Texture, cellSize int32) {
	var err error
	if rect.W > rect.H {
		cellSize = rect.H / (int32(dim) + 1)
	} else {
		cellSize = rect.W / (int32(dim) + 1)
	}
	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, cellSize*dim, cellSize*dim)
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texture)
	renderer.SetDrawColor(bg.R, bg.G, bg.B, bg.A)
	renderer.Clear()
	renderer.SetDrawColor(fg.R, fg.G, fg.B, fg.A)
	renderer.SetDrawColor(64, 0, 0, 255)
	renderer.DrawRect(&sdl.Rect{0, 0, cellSize * dim, cellSize * dim})
	for dy := 0; dy < int(dim); dy++ {
		for dx := 0; dx <= int(dim); dx++ {
			x := int32(dx) * cellSize
			y := int32(dy) * cellSize
			renderer.DrawLine(0, y, cellSize*dim, y)
			renderer.DrawLine(x, 0, x, cellSize*dim)
		}
	}
	renderer.SetRenderTarget(nil)
	return texture, cellSize
}
