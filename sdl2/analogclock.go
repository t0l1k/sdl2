package sdl2

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type AnalogClock struct {
	renderer                                   *sdl.Renderer
	rect                                       sdl.Rect
	texFace                                    *sdl.Texture
	fg, bg, secHandColor, tweentyPointColor    sdl.Color
	hourHand, minuteHand, secondHand, mSecHand *ClockHand
	drawMsec                                   bool
	tweentyBlinkTimer                          *BlinkTimer
	tipTweentyPoint                            sdl.Point
	fnAnalog                                   GetTime
	show                                       bool
}

func NewAnalogClock(renderer *sdl.Renderer, rect sdl.Rect, fg, secHandColor, tweentyPointColor, bg sdl.Color, fnAnalog GetTime) *AnalogClock {
	texFace := NewClockFace(renderer, rect, fg, bg)
	rectWidth, rectHeight := int32(float64(rect.H)*0.470), int32(float64(rect.H)*0.02)
	mSecHand := NewSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1), rectHeight / 2}, sdl.Point{int32(float64(rectHeight) * 0.2), rectHeight / 4}, secHandColor, bg)
	secondHand := NewSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1.13), rectHeight / 2}, sdl.Point{int32(float64(rectWidth) * 0.2), rectHeight / 4}, secHandColor, bg)
	minuteHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.9), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	hourHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.7), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	tipTweentyPoint := getTip(sdl.Point{rect.W / 2, rect.H / 2}, 0/60, float64(rect.H/2-(rect.H/90)*3), 0, 0)
	blinkTimer := NewBlinkTimer(1000 / 2)
	go blinkTimer.Run()
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

func (s *AnalogClock) GetShow() bool {
	return s.show
}

func (s *AnalogClock) SetShow(show bool) {
	s.show = show
}

func (s *AnalogClock) Render(renderer *sdl.Renderer) {
	if s.show {
		if err := renderer.Copy(s.texFace, nil, &s.rect); err != nil {
			panic(err)
		}
		if s.tweentyBlinkTimer.IsOn() {
			FillCircle(s.renderer, s.rect.X+s.tipTweentyPoint.X, s.rect.Y+s.tipTweentyPoint.Y, s.rect.H/200, s.tweentyPointColor)
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
		s.secondHand.Update((float64(second) + s.mSecHand.GetFraction()) / 60.0)
		s.minuteHand.Update((float64(minute) + s.secondHand.GetFraction()) / 60.0)
		s.hourHand.Update((float64(hour) + s.minuteHand.GetFraction()) / 12.0)
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

func NewClockFace(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texClockFace *sdl.Texture) {
	var err error
	if texClockFace, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H); err != nil {
		panic(err)
	}
	center := sdl.Point{rect.W / 2, rect.H / 2}
	margin := rect.H / 90
	renderer.SetRenderTarget(texClockFace)
	texClockFace.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, sdl.Color{128, 255, 255, 255})
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
		FillCircle(renderer, tip.X, tip.Y, radius, bg)
		FillCircle(renderer, tip.X, tip.Y, radius/2, fg)
	}
	renderer.SetRenderTarget(nil)
	return texClockFace
}

type ClockHand struct {
	renderer        *sdl.Renderer
	texture         *sdl.Texture
	rect, paintRect sdl.Rect
	handCenter      sdl.Point
	width, height   int32
	fg, bg          sdl.Color
	angle, fraction float64
}

func (s *ClockHand) GetFraction() float64 { return s.fraction }

func (s *ClockHand) Render(renderer *sdl.Renderer) {
	if err := renderer.CopyEx(s.texture, nil, &s.paintRect, s.angle, &s.handCenter, sdl.FLIP_NONE); err != nil {
		panic(err)
	}
}

func (s *ClockHand) Update(percent float64) {
	s.fraction = percent
	s.angle = getAngle(percent)
}
func (s *ClockHand) Event(sdl.Event) {}
func (s *ClockHand) Destroy()        { s.texture.Destroy() }

func NewSmallHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)
	renderer.FillRect(&sdl.Rect{0, 0, center.X - center.X/4, rect.H})
	renderer.FillRect(&sdl.Rect{0, rect.H / 3, rect.W, rect.H - rect.H/3*2})
	FillCircle(renderer, center.X, center.Y, rect.H/5, bg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &ClockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

func NewBigHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)
	renderer.DrawLine(0, 0, rect.W-rect.H/2, rect.H/4)
	renderer.DrawLine(rect.W-rect.H/2, rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(rect.W-rect.H/2, (rect.H-1)-rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(0, rect.H-1, rect.W-rect.H/2, (rect.H-1)-rect.H/4)
	renderer.DrawLine(0, 0, 0, rect.H)
	FillCircle(renderer, center.X, center.Y, rect.H/5, fg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &ClockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

type BlinkTimer struct {
	delay   uint32
	blinkOn bool
	running bool
}

func NewBlinkTimer(delay uint32) *BlinkTimer {
	return &BlinkTimer{
		delay:   delay,
		running: true,
	}
}

func (s *BlinkTimer) IsOn() bool { return s.blinkOn }

func (s *BlinkTimer) switchOn() { s.blinkOn = !s.blinkOn }

func (s *BlinkTimer) Run() {
	for s.running {
		sdl.Delay(s.delay)
		s.switchOn()
	}
}

func (s *BlinkTimer) Stop() {
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

func FillCircle(renderer *sdl.Renderer, x0, y0, radius int32, color sdl.Color) {
	setColor(renderer, color)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(x0+x, y0+y)
			}
		}
	}
}
