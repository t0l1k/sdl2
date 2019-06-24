package ui

import "github.com/veandco/go-sdl2/sdl"

type Sprite interface {
	Render(*sdl.Renderer)
	Update()
	Event(sdl.Event)
	Destroy()
}
