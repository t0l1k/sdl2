package sdl2

import (
	"github.com/t0l1k/sdl2/clock"
	"github.com/t0l1k/sdl2/fifteen"
	"github.com/t0l1k/sdl2/life"
	"github.com/t0l1k/sdl2/mines"
	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AppManager struct {
	screenIndex    int
	title          string
	titles         []string
	renderer       *sdl.Renderer
	rect           sdl.Rect
	sprites        []ui.Sprite
	font           *ttf.Font
	fg, bg, lineBg sdl.Color
	menuLine       *MenuLine
	analogClock    *clock.AnalogClock
	life           *life.Life
	fifteen        *fifteen.Game
	mines          *mines.MinesBoard
	fnAnalog       clock.GetTime
	fnQuit         func()
}

func NewAppManager(renderer *sdl.Renderer, title string, width, height int32, quit func()) *AppManager {
	titles := make([]string, 0)
	titles = append(titles, title, "Clock", "Conway's Life", "Game Fifteen", "Game Minesweeper")
	return &AppManager{
		screenIndex: 0,
		title:       title,
		titles:      titles,
		renderer:    renderer,
		rect:        sdl.Rect{0, 0, width, height},
		fg:          sdl.Color{0, 0, 0, 255},
		bg:          sdl.Color{192, 192, 192, 0},
		lineBg:      sdl.Color{128, 128, 128, 0},
		fnQuit:      quit,
	}
}

func (s *AppManager) setup() {
	var err error
	fontSize := int(float64(s.rect.H) * 0.03) // Главная константа перерисовки экрана
	s.font, err = ttf.OpenFont("assets/Roboto-Regular.ttf", fontSize)
	if err != nil {
		panic(err)
	}
	lineHeight := int32(float64(fontSize) * 1.5)
	s.menuLine = NewMenuLine(
		s.title,
		sdl.Rect{0, 0, s.rect.W, lineHeight},
		s.fg,
		s.lineBg,
		s.renderer,
		s.font,
		func() { s.quit() })
	s.sprites = append(s.sprites, s.menuLine)

	s.analogClock = clock.NewAnalogClock(
		s.renderer,
		sdl.Rect{0, lineHeight, s.rect.W, s.rect.H - lineHeight},
		s.fg,
		sdl.Color{255, 0, 0, 255},
		s.bg, s.bg,
		s.fnAnalog)
	s.sprites = append(s.sprites, s.analogClock)
	s.life = life.NewLife(
		128, 255,
		s.renderer,
		sdl.Rect{0, lineHeight, s.rect.W, s.rect.H - lineHeight},
		s.fg, s.bg)
	s.sprites = append(s.sprites, s.life)
	s.fifteen = fifteen.NewGame(
		s.renderer,
		sdl.Rect{0, lineHeight, s.rect.W, s.rect.H - lineHeight})
	s.fifteen.Setup()
	s.sprites = append(s.sprites, s.fifteen)
	s.mines = mines.NewMinesBoard(
		s.renderer,
		sdl.Rect{0, lineHeight, s.rect.W, s.rect.H - lineHeight})
	s.mines.Setup()
	s.sprites = append(s.sprites, s.mines)
	fnArr := []func(){s.selectClock, s.selectLife, s.selectFifteen, s.selectMines}
	for i := 0; i < len(fnArr); i++ {
		btn := ui.NewButton(s.renderer, s.titles[i+1], sdl.Rect{0, lineHeight + lineHeight*int32(i), lineHeight * 10, lineHeight}, s.fg, s.bg, s.font, fnArr[i])
		s.sprites = append(s.sprites, btn)
	}
}

func (s *AppManager) selectMines() {
	s.screenIndex = 4
	s.title = s.titles[s.screenIndex]
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.life.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(true)
	s.sprites = s.sprites[:len(s.sprites)-len(s.titles)+1]
}

func (s *AppManager) selectFifteen() {
	s.screenIndex = 3
	s.title = s.titles[s.screenIndex]
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.life.SetShow(false)
	s.mines.SetShow(false)
	s.fifteen.SetShow(true)
	s.sprites = s.sprites[:len(s.sprites)-len(s.titles)+1]
}

func (s *AppManager) selectLife() {
	s.screenIndex = 2
	s.title = s.titles[s.screenIndex]
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(false)
	s.life.SetShow(true)
	s.sprites = s.sprites[:len(s.sprites)-len(s.titles)+1]
}

func (s *AppManager) selectClock() {
	s.fnAnalog = clock.Get
	s.screenIndex = 1
	s.title = s.titles[s.screenIndex]
	s.Destroy()
	s.setup()
	s.life.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(false)
	s.analogClock.SetShow(true)
	s.sprites = s.sprites[:len(s.sprites)-len(s.titles)+1]
}

func (s *AppManager) selectMain() {
	s.screenIndex = 0
	s.title = s.titles[s.screenIndex]
	s.Destroy()
	s.setup()
	s.life.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(false)
	s.analogClock.SetShow(false)
}

func (s *AppManager) quit() {
	switch s.screenIndex {
	case 0:
		s.fnQuit()
	default:
		s.selectMain()
	}

}

func (s *AppManager) Render(renderer *sdl.Renderer) {
	ui.SetColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	for _, sprite := range s.sprites {
		sprite.Render(renderer)
	}
}

func (s *AppManager) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

func (s *AppManager) Event(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_ESCAPE && t.State == sdl.RELEASED {
			s.quit()
		}
	}
	for _, sprite := range s.sprites {
		sprite.Event(e)
	}
}

func (s *AppManager) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.font.Close()
	s.sprites = s.sprites[:0]
}
