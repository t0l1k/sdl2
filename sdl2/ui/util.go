package ui

import "github.com/veandco/go-sdl2/sdl"

func SetColor(renderer *sdl.Renderer, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}
