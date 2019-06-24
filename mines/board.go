package mines

import (
	"fmt"
	"strconv"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MinesBoard struct {
	gameConfig *gameConfig
	field      *minesField
	renderer   *sdl.Renderer
	rect       sdl.Rect
	fg, bg     sdl.Color
	sprites    []ui.Sprite
	font       *ttf.Font
	show       bool
	lastTime   uint32
	status     *statusLine
	message    *ui.MessageBox
	fnMessage  func()
}

func NewMinesBoard(row, column, mines int, renderer *sdl.Renderer, rect sdl.Rect) *MinesBoard {
	gameConfig := newGameConfig(row, column, mines)
	fg, bg := sdl.Color{0, 0, 0, 255}, sdl.Color{96, 192, 32, 255}
	return &MinesBoard{
		gameConfig: gameConfig,
		renderer:   renderer,
		rect:       rect,
		fg:         fg,
		bg:         bg,
		show:       false,
	}
}

func (s *MinesBoard) Setup() {
	var (
		fontName = "assets/Roboto-Regular.ttf"
		cellSize int32
		fnSize   int
		err      error
	)
	cellSize = s.calcCellSize()
	fnSize = int(float64(cellSize) * 0.8)
	s.font, err = ttf.OpenFont(fontName, fnSize)
	if err != nil {
		panic(err)
	}
	s.field = newMinesField(s.gameConfig.row, s.gameConfig.column, s.gameConfig.mines)

	marginX := (s.rect.W - cellSize*int32(s.gameConfig.row)) / 2
	marginY := (s.rect.H - cellSize*int32(s.gameConfig.column+1)) / 2
	for i, cell := range s.field.getBoard() {
		var x, y int32
		x = s.rect.X + cellSize*int32(i%s.gameConfig.row) + marginX
		y = s.rect.Y + cellSize*int32(i/s.gameConfig.row) + marginY
		btn := ui.NewButton(s.renderer, cell.String(), sdl.Rect{x, y, cellSize, cellSize}, s.fg, s.bg, s.font, s.checkGameLogic)
		s.sprites = append(s.sprites, btn)
	}
	s.status = newStatusLine(sdl.Rect{s.rect.X, s.rect.Y + cellSize*int32(s.gameConfig.column), s.rect.W, int32(cellSize)}, s.fg, s.bg, s.renderer, s.font)
	s.fnMessage = s.newGame
	s.message = ui.NewMessageBox("Minesweeper", "test message", s.rect, s.fg, s.bg, s.renderer, s.fnMessage)
	s.SetShow(false)
}

func (s *MinesBoard) calcCellSize() (cellSize int32) {
	if s.rect.W > s.rect.H && s.gameConfig.row > s.gameConfig.column+1 {
		cellSize = s.rect.W / int32(s.gameConfig.row)
		if cellSize*int32(s.gameConfig.column+1) > s.rect.H {
			cellSize = s.rect.H / int32(s.gameConfig.column+1)
		}
	} else {
		cellSize = s.rect.H / int32(s.gameConfig.column+1)
		if cellSize*int32(s.gameConfig.row) > s.rect.W {
			cellSize = s.rect.W / int32(s.gameConfig.row)
		}
	}
	return cellSize
}

func (s *MinesBoard) newGame() {
	s.field.destroy()
	s.field = newMinesField(s.gameConfig.row, s.gameConfig.column, s.gameConfig.mines)
	s.gameConfig.stop()
	s.message.SetShow(false)
	s.paint()
}

func (s *MinesBoard) reset() {
	s.field.reset()
	s.gameConfig.stop()
	s.message.SetShow(false)
	s.paint()
}

func (s *MinesBoard) paint() {
	for idx, sprite := range s.sprites {
		switch sprite.(type) {
		case *ui.Button:
			if s.field.getBoard()[idx].getState() == opened {
				sprite.(*ui.Button).SetFGBG(s.bg, s.fg)
			} else {
				sprite.(*ui.Button).SetFGBG(s.fg, s.bg)
			}
			sprite.(*ui.Button).SetText(s.field.getBoard()[idx].String())
		}
	}
	str := fmt.Sprintf("%02v", s.field.getStatus())
	s.status.setFlags(str)
}

func (s *MinesBoard) checkGameLogic() {
	for i, sprite := range s.sprites {
		button := sprite.(*ui.Button)
		if button.IsPressed() {
			x, y := s.field.getPos(i)
			cell := s.field.getCell(x, y)
			switch s.field.getState() {
			case gameStart:
				if cell.getState() == closed {
					s.field.shuffle(x, y)
					s.gameConfig.start()
					s.lastTime = sdl.GetTicks()
					s.field.open(x, y)
				}
			case gamePlay:
				if button.IsPressedLeft() {
					if cell.getState() == closed {
						s.field.open(x, y)
					} else if cell.getState() == opened {
						s.field.autoMarkFlags(x, y)
					}
					if s.field.isWin() || s.field.isGameOver() {
						s.gameConfig.stop()
						if s.field.isWin() {
							s.message.SetMessage("Board Solved")
							s.fnMessage = s.newGame
							s.message.SetShow(true)
						}
						if s.field.isGameOver() {
							s.fnMessage = s.reset
							s.message.SetMessage("Game Over")
							s.message.SetShow(true)
						}
					}
				} else if button.IsPressedRight() {
					s.field.markFlag(x, y)
				}
			case gamePause:
			case gameWin:
			case gameOver:
			}
			s.paint()
		}
	}
}

func (s *MinesBoard) SetShow(show bool) {
	s.show = show
}

func (s *MinesBoard) Render(renderer *sdl.Renderer) {
	if s.show {
		for _, sprite := range s.sprites {
			sprite.Render(s.renderer)
		}
		s.status.Render(renderer)
		s.message.Render(renderer)
	}
}

func (s *MinesBoard) Update() {
	if s.show {
		if s.field.getState() == gamePlay && sdl.GetTicks()-s.lastTime > 999 {
			str := fmt.Sprintf("%03v", strconv.Itoa(s.gameConfig.getStopper()))
			s.status.setStopper(str)
			s.lastTime = sdl.GetTicks()
		}
		for _, sprite := range s.sprites {
			sprite.Update()
		}
	}
	s.message.Update()
}

func (s *MinesBoard) Event(event sdl.Event) {
	if s.show {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_SPACE && t.State == sdl.RELEASED {
				s.reset()
			} else if t.Keysym.Sym == sdl.K_RETURN && t.State == sdl.RELEASED {
				s.newGame()
			}
		}
		for _, sprite := range s.sprites {
			sprite.Event(event)
		}
	}
	s.message.Event(event)
}

func (s *MinesBoard) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
	s.status.Destroy()
	s.message.Destroy()
	s.font.Close()
}
