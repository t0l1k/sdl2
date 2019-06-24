package ui

import "github.com/veandco/go-sdl2/sdl"

func SetColor(renderer *sdl.Renderer, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

func FillCircle(renderer *sdl.Renderer, x0, y0, radius int32, color sdl.Color) {
	SetColor(renderer, color)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(x0+x, y0+y)
			}
		}
	}
}
