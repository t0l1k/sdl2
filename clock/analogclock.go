package clock

import (
	"math"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
)

type GetTime func() (int, int, int, int)

type AnalogClock struct {
	renderer                                   *sdl.Renderer
	rect                                       sdl.Rect
	texFace                                    *sdl.Texture
	fg, bg, secHandColor, tweentyPointColor    sdl.Color
	hourHand, minuteHand, secondHand, mSecHand *clockHand
	drawMsec                                   bool
	tweentyBlinkTimer                          *blinkTimer
	tipTweentyPoint                            sdl.Point
	fnAnalog                                   GetTime
	show                                       bool
}

func NewAnalogClock(renderer *sdl.Renderer, rect sdl.Rect, fg, secHandColor, tweentyPointColor, bg sdl.Color, fnAnalog GetTime) *AnalogClock {
	texFace := newClockFace(renderer, rect, fg, bg)
	rectWidth, rectHeight := int32(float64(rect.H)*0.470), int32(float64(rect.H)*0.02)
	mSecHand := newSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1), rectHeight / 2}, sdl.Point{int32(float64(rectHeight) * 0.2), rectHeight / 4}, secHandColor, bg)
	secondHand := newSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1.13), rectHeight / 2}, sdl.Point{int32(float64(rectWidth) * 0.2), rectHeight / 4}, secHandColor, bg)
	minuteHand := newBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.9), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	hourHand := newBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.7), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	tipTweentyPoint := getTip(sdl.Point{rect.W / 2, rect.H / 2}, 0/60, float64(rect.H/2-(rect.H/90)*3), 0, 0)
	blinkTimer := newBlinkTimer(1000 / 2)
	go blinkTimer.run()
	return &AnalogClock{
		renderer:          renderer,
		rect:              rect,
		texFace:           texFace,
		fg:                fg,
		bg:                bg,
		tweentyPointColor: tweentyPointColor,
		hourHand:          hourHand,
		minuteHand:        minuteHand,
		secondHand:        secondHand,
		mSecHand:          mSecHand,
		tweentyBlinkTimer: blinkTimer,
		tipTweentyPoint:   tipTweentyPoint,
		fnAnalog:          fnAnalog,
	}
}

func (s *AnalogClock) SetShow(show bool) {
	s.show = show
}

func (s *AnalogClock) Render(renderer *sdl.Renderer) {
	if s.show {
		if err := renderer.Copy(s.texFace, nil, &s.rect); err != nil {
			panic(err)
		}
		if s.tweentyBlinkTimer.isOn() {
			ui.FillCircle(s.renderer, s.rect.X+s.tipTweentyPoint.X, s.rect.Y+s.tipTweentyPoint.Y, s.rect.H/200, s.tweentyPointColor)
		}
		s.hourHand.Render(s.renderer)
		s.minuteHand.Render(s.renderer)
		s.secondHand.Render(s.renderer)
		if s.drawMsec {
			s.mSecHand.Render(s.renderer)
		}
	}
}

func (s *AnalogClock) Update() {
	if s.show {
		mSec, second, minute, hour := s.fnAnalog()
		s.mSecHand.Update(float64(mSec) / 1000.0)
		s.secondHand.Update((float64(second) + s.mSecHand.getFraction()) / 60.0)
		s.minuteHand.Update((float64(minute) + s.secondHand.getFraction()) / 60.0)
		s.hourHand.Update((float64(hour) + s.minuteHand.getFraction()) / 12.0)
	}
}

func (s *AnalogClock) Event(sdl.Event) {
	if s.show {
	}
}
func (s *AnalogClock) Destroy() {
	s.texFace.Destroy()
	s.hourHand.Destroy()
	s.minuteHand.Destroy()
	s.secondHand.Destroy()
	s.mSecHand.Destroy()
}

func newClockFace(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texClockFace *sdl.Texture) {
	var err error
	if texClockFace, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H); err != nil {
		panic(err)
	}
	center := sdl.Point{rect.W / 2, rect.H / 2}
	margin := rect.H / 90
	renderer.SetRenderTarget(texClockFace)
	texClockFace.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, sdl.Color{128, 255, 255, 255})
	var x, y int32
	for y = 0; y < rect.H; y += 5 {
		for x = 0; x < rect.W; x += 5 {
			renderer.DrawLine(x, 0, x, rect.H)
			renderer.DrawLine(0, y, rect.W, y)
		}
	}
	for i := 0; i < 60; i++ {
		var (
			tip    sdl.Point
			radius int32
		)
		if i%5 == 0 {
			radius = margin * 2
		} else {
			radius = margin
		}
		tip = getTip(center, float64(i)/60.0, float64(center.Y-margin*3), 0, 0)
		ui.FillCircle(renderer, tip.X, tip.Y, radius, bg)
		ui.FillCircle(renderer, tip.X, tip.Y, radius/2, fg)
	}
	renderer.SetRenderTarget(nil)
	return texClockFace
}

type clockHand struct {
	renderer        *sdl.Renderer
	texture         *sdl.Texture
	rect, paintRect sdl.Rect
	handCenter      sdl.Point
	width, height   int32
	fg, bg          sdl.Color
	angle, fraction float64
}

func (s *clockHand) getFraction() float64 { return s.fraction }

func (s *clockHand) Render(renderer *sdl.Renderer) {
	if err := renderer.CopyEx(s.texture, nil, &s.paintRect, s.angle, &s.handCenter, sdl.FLIP_NONE); err != nil {
		panic(err)
	}
}

func (s *clockHand) Update(percent float64) {
	s.fraction = percent
	s.angle = getAngle(percent)
}
func (s *clockHand) Event(sdl.Event) {}
func (s *clockHand) Destroy()        { s.texture.Destroy() }

func newSmallHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *clockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, fg)
	renderer.FillRect(&sdl.Rect{0, 0, center.X - center.X/4, rect.H})
	renderer.FillRect(&sdl.Rect{0, rect.H / 3, rect.W, rect.H - rect.H/3*2})
	ui.FillCircle(renderer, center.X, center.Y, rect.H/5, bg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &clockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

func newBigHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *clockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, fg)
	renderer.DrawLine(0, 0, rect.W-rect.H/2, rect.H/4)
	renderer.DrawLine(rect.W-rect.H/2, rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(rect.W-rect.H/2, (rect.H-1)-rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(0, rect.H-1, rect.W-rect.H/2, (rect.H-1)-rect.H/4)
	renderer.DrawLine(0, 0, 0, rect.H)
	ui.FillCircle(renderer, center.X, center.Y, rect.H/5, fg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &clockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

type blinkTimer struct {
	delay   uint32
	blinkOn bool
	running bool
}

func newBlinkTimer(delay uint32) *blinkTimer {
	return &blinkTimer{
		delay:   delay,
		running: true,
	}
}

func (s *blinkTimer) isOn() bool { return s.blinkOn }

func (s *blinkTimer) switchOn() { s.blinkOn = !s.blinkOn }

func (s *blinkTimer) run() {
	for s.running {
		sdl.Delay(s.delay)
		s.switchOn()
	}
}

func (s *blinkTimer) stop() {
	s.running = false
}

func getTip(center sdl.Point, percent, lenght, width, height float64) (tip sdl.Point) {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(radians)
	cosine := math.Cos(radians)
	tip.X = center.X + int32(lenght*sine-width)
	tip.Y = center.Y + int32(lenght*cosine-height)
	return tip
}

func getAngle(percent float64) float64 {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	angle := (radians * -180 / math.Pi) + 90
	return angle
}
