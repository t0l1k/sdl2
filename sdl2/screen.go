package sdl2

import (
	"fmt"

	"github.com/t0l1k/sdl2/sdl2/ui"

	"github.com/veandco/go-sdl2/sdl"
)

type Screen struct {
	title                  string
	window                 *sdl.Window
	renderer               *sdl.Renderer
	width, height          int32
	flags                  uint32
	running                bool
	fpsCountTime, fpsCount uint32
	sprites                []ui.Sprite
	appManager             *AppManager
}

func NewScreen(title string, window *sdl.Window, renderer *sdl.Renderer, width, height int32) *Screen {
	return &Screen{
		title:    title,
		window:   window,
		renderer: renderer,
		width:    width,
		height:   height,
	}
}

func (s *Screen) setup() {
	s.appManager = NewAppManager(s.renderer, s.title, s.width, s.height, s.quit)
	s.appManager.setup()
	s.sprites = append(s.sprites, s.appManager)
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
		if t.Keysym.Sym == sdl.K_F11 && t.State == sdl.RELEASED {
			s.setMode()
		}
	case *sdl.WindowEvent:
		switch t.Event {
		case sdl.WINDOWEVENT_RESIZED:
			s.width, s.height = t.Data1, t.Data2
			s.Destroy()
			s.setup()
			// fmt.Println("window resized", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			// fmt.Println("window focus gained", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_LOST:
			// fmt.Println("window focus lost", s.width, s.height)
		case sdl.WINDOW_MINIMIZED:
			// fmt.Println("window minimized", s.width, s.height)
			s.Destroy()
		case sdl.WINDOWEVENT_RESTORED:
			// fmt.Println("window restored", s.width, s.height)
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
}

func (s *Screen) quit() {
	s.running = false
}
