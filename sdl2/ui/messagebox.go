package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MessageBox struct {
	title, message string
	renderer       *sdl.Renderer
	texture        *sdl.Texture
	rect           sdl.Rect
	font           *ttf.Font
	fg, bg         sdl.Color
	btnOk          *Button
	show           bool
}

func NewMessageBox(title, message string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, fnOk func()) *MessageBox {
	fnSize := fontSize(rect.H)
	font, err := ttf.OpenFont("assets/Roboto-Regular.ttf", fnSize)
	if err != nil {
		panic(err)
	}
	wFont, hFont, _ := font.SizeUTF8(message)
	w, h := int32(float64(wFont)*1.3), int32(float64(hFont)*5)
	x, y := rect.X+(rect.W-w)/2, rect.Y+(rect.H-h)/2
	rect = sdl.Rect{x, y, w, h}
	fmt.Println(rect)
	texture := newMessageBoxTexture(title, message, rect, sdl.Color{255, 0, 0, 255}, sdl.Color{0, 128, 0, 255}, renderer, font)
	btnOk := NewButton(renderer, "OK", sdl.Rect{x + (w-w/2)/2, y + (h - h/5), w / 2, h / 5}, fg, bg, font, fnOk)
	return &MessageBox{
		title:    title,
		message:  message,
		rect:     rect,
		font:     font,
		fg:       fg,
		bg:       bg,
		renderer: renderer,
		texture:  texture,
		show:     false,
		btnOk:    btnOk,
	}
}

func newMessageBoxTexture(title, message string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *sdl.Texture {
	lblTitle := NewLabel(title, sdl.Point{3, 0}, sdl.Color{0, 0, 255, 255}, renderer, font)
	defer lblTitle.Destroy()
	lblMessage := NewLabel(message, sdl.Point{0, 0}, sdl.Color{255, 0, 255, 255}, renderer, font)
	defer lblMessage.Destroy()
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texture)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)
	renderer.DrawRect(&sdl.Rect{0, 0, rect.W, rect.H})
	_, hFont, _ := font.SizeUTF8(title)
	renderer.DrawRect(&sdl.Rect{0, 0, rect.W, int32(hFont) + 3})
	lblTitle.Render(renderer)
	rectMessage := lblMessage.GetSize()
	x, y := (rect.W-rectMessage.W)/2, (rect.H-rectMessage.H)/2
	lblMessage.SetPos(sdl.Point{x, y})
	lblMessage.Render(renderer)
	renderer.SetRenderTarget(nil)
	return texture
}

func (s *MessageBox) GetTexture() *sdl.Texture {
	return s.texture
}

func (s *MessageBox) SetMessage(message string) {
	s.texture.Destroy()
	s.message = message
	s.texture = newMessageBoxTexture(s.title, s.message, s.rect, s.fg, s.bg, s.renderer, s.font)
	_, _, w, h, _ := s.texture.Query()
	s.SetSize(w, h)
}

func (s *MessageBox) GetShow() bool {
	return s.show
}

func (s *MessageBox) SetShow(show bool) {
	s.show = show
}

func (s *MessageBox) GetSize() sdl.Rect {
	return s.rect
}

func (s *MessageBox) SetSize(w, h int32) {
	s.rect.W = w
	s.rect.H = h
}

func (s *MessageBox) SetPos(pos sdl.Point) {
	s.rect.X = pos.X
	s.rect.Y = pos.Y
}

func (s *MessageBox) Render(renderer *sdl.Renderer) {
	if s.show {
		if err := renderer.Copy(s.texture, nil, &s.rect); err != nil {
			panic(err)
		}
		s.btnOk.Render(renderer)
	}
}

func (s *MessageBox) Update() {
	if s.show {
		s.btnOk.Update()
	}
}
func (s *MessageBox) Event(e sdl.Event) {
	if s.show {
		s.btnOk.Event(e)
	}
}
func (s *MessageBox) Destroy() {
	s.texture.Destroy()
	s.btnOk.Destroy()
	s.font.Close()
}

func fontSize(height int32) int {
	return int(float64(height) * 0.08) // Главная константа перерисовки экрана размер шрифта
}

func lineHeight(fontSize int) int32 {
	return int32(float64(fontSize) * 1.5) // Высота строки меню и строки статуса
}
