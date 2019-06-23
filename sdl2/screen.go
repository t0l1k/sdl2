package sdl2

import (
	"fmt"

	"github.com/t0l1k/sdl2/clock"
	"github.com/t0l1k/sdl2/fifteen"
	"github.com/t0l1k/sdl2/life"
	"github.com/t0l1k/sdl2/mines"
	"github.com/t0l1k/sdl2/sdl2/ui"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Screen struct {
	title                  string
	window                 *sdl.Window
	renderer               *sdl.Renderer
	width, height          int32
	flags                  uint32
	running                bool
	bg, fg, lineBg         sdl.Color
	fpsCountTime, fpsCount uint32
	sprites                []ui.Sprite
	font                   *ttf.Font
	menuLine               *MenuLine
	statusLine             *StatusLine
	analogClock            *clock.AnalogClock
	life                   *life.Life
	fifteen                *fifteen.Game
	mines                  *mines.MinesBoard
	fnAnalog               clock.GetTime
}

func NewScreen(title string, window *sdl.Window, renderer *sdl.Renderer, width, height int32) *Screen {
	return &Screen{
		title:    title,
		window:   window,
		renderer: renderer,
		width:    width,
		height:   height,
		bg:       sdl.Color{192, 192, 192, 0},
		fg:       sdl.Color{0, 0, 0, 255},
		lineBg:   sdl.Color{128, 128, 128, 0},
		fnAnalog: clock.Get,
	}
}

func (s *Screen) setup() {
	var err error
	fontSize := int(float64(s.height) * 0.03) // Главная константа перерисовки экрана
	s.font, err = ttf.OpenFont("assets/Roboto-Regular.ttf", fontSize)
	if err != nil {
		panic(err)
	}
	lineHeight := int32(float64(fontSize) * 1.5)
	s.menuLine = NewMenuLine(
		s.title,
		sdl.Rect{0, 0, s.width, lineHeight},
		s.fg,
		s.lineBg,
		s.renderer,
		s.font,
		func() { s.quit() })
	s.sprites = append(s.sprites, s.menuLine)
	s.statusLine = NewStatusLine(
		sdl.Rect{0, s.height - lineHeight, s.width, lineHeight},
		s.fg,
		s.lineBg,
		s.renderer,
		s.font,
		s.selectClock,
		s.selectLife,
		s.selectFifteen,
		s.selectMines)
	s.sprites = append(s.sprites, s.statusLine)
	s.analogClock = clock.NewAnalogClock(
		s.renderer,
		sdl.Rect{0, lineHeight, s.width, s.height - lineHeight*2},
		s.fg,
		sdl.Color{255, 0, 0, 255},
		s.bg,
		s.bg,
		s.fnAnalog)
	s.sprites = append(s.sprites, s.analogClock)
	s.life = life.NewLife(128, 250, s.renderer, sdl.Rect{0, lineHeight, s.width, s.height - lineHeight*2}, s.fg, s.bg)
	s.sprites = append(s.sprites, s.life)
	s.fifteen = fifteen.NewGame(s.renderer, sdl.Rect{0, lineHeight, s.width, s.height - lineHeight*2})
	s.fifteen.Setup()
	s.sprites = append(s.sprites, s.fifteen)
	s.mines = mines.NewMinesBoard(8, 8, 12, s.renderer, sdl.Rect{0, lineHeight, s.width, s.height - lineHeight*2})
	s.mines.Setup()
	s.sprites = append(s.sprites, s.mines)
}

func (s *Screen) selectMines() {
	s.title = "Game Minesweeper"
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.life.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(true)
}

func (s *Screen) selectFifteen() {
	s.title = "Game Fifteen"
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.life.SetShow(false)
	s.mines.SetShow(false)
	s.fifteen.SetShow(true)
}

func (s *Screen) selectLife() {
	s.title = "Conway's Life"
	s.Destroy()
	s.setup()
	s.analogClock.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(false)
	s.life.SetShow(true)
}

func (s *Screen) selectClock() {
	s.fnAnalog = clock.Get
	s.title = "Clock"
	s.Destroy()
	s.setup()
	s.life.SetShow(false)
	s.fifteen.SetShow(false)
	s.mines.SetShow(false)
	s.analogClock.SetShow(true)
}

func (s *Screen) setMode() {
	if s.flags == 0 {
		s.flags = sdl.WINDOW_FULLSCREEN_DESKTOP
		mode, err := sdl.GetCurrentDisplayMode(0)
		if err != nil {
			panic(err)
		}
		s.width, s.height = mode.W, mode.H
	} else {
		s.flags = 0
		s.width, s.height = 800, 600
	}
	s.window.SetFullscreen(s.flags)
	s.window.SetSize(s.width, s.height)
	s.Destroy()
	s.setup()
}

func (s *Screen) Event() {
	event := sdl.WaitEventTimeout(3)
	switch t := event.(type) {
	case *sdl.QuitEvent:
		s.quit()
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_ESCAPE && t.State == sdl.RELEASED {
			s.quit()
		}
		if t.Keysym.Sym == sdl.K_F11 && t.State == sdl.RELEASED {
			s.setMode()
		}
	case *sdl.WindowEvent:
		switch t.Event {
		case sdl.WINDOWEVENT_RESIZED:
			s.width, s.height = t.Data1, t.Data2
			s.Destroy()
			s.setup()
			fmt.Println("window resized", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			fmt.Println("window focus gained", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_LOST:
			fmt.Println("window focus lost", s.width, s.height)
		case sdl.WINDOW_MINIMIZED:
			fmt.Println("window minimized", s.width, s.height)
			s.Destroy()
		case sdl.WINDOWEVENT_RESTORED:
			fmt.Println("window restored", s.width, s.height)
			s.setup()
		}
	}
	for _, sprite := range s.sprites {
		sprite.Event(event)
	}
}

func (s *Screen) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
	if sdl.GetTicks()-s.fpsCountTime > 999 {
		s.window.SetTitle(fmt.Sprintf("%s fps:%v", s.title, s.fpsCount))
		s.fpsCount = 0
		s.fpsCountTime = sdl.GetTicks()
	}
}

func (s *Screen) Render() {
	setColor(s.renderer, s.bg)
	s.renderer.Clear()
	for _, sprite := range s.sprites {
		sprite.Render(s.renderer)
	}
	s.renderer.Present()
	s.fpsCount++
}

func (s *Screen) Run() {
	s.setup()
	frameRate := uint32(1000 / 60)
	lastTime := sdl.GetTicks()
	s.running = true
	for s.running {
		now := sdl.GetTicks()
		if now >= lastTime {
			i := 0
			for {
				s.Event()
				s.Update()
				lastTime += frameRate
				now = sdl.GetTicks()
				if lastTime > now {
					break
				}
				i++
				if i >= 3 {
					lastTime = now + frameRate
					break
				}
			}
			s.Render()
		} else {
			sdl.Delay(lastTime - now)
		}
	}
}

func (s *Screen) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
	s.font.Close()
}

func (s *Screen) quit() {
	s.running = false
}

func setColor(renderer *sdl.Renderer, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}
