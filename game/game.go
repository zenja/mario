package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
)

type Game struct {
	// start position (left top) of camera
	XCam, YCam int32

	gra          *graphic.Graphic
	currentLevel *level.Level
}

func NewGame() *Game {
	return &Game{gra: graphic.New()}
}

func (game *Game) LoadLevel(filename string) {
	game.currentLevel = level.ParseLevelFromFile(filename)
}

func (game *Game) Quit() {
	game.gra.DestroyAndQuit()
}

func (game *Game) StartGameLoop() {
	var event sdl.Event
	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			frameStart := sdl.GetTicks()

			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_LEFT:
				case sdl.SCANCODE_RIGHT:
				case sdl.SCANCODE_SPACE:
				}
			}

			// render
			game.currentLevel.Draw(game.gra, game.XCam, game.YCam)

			frameTime := sdl.GetTicks() - frameStart

			if frameTime < graphic.DELAY_TIME {
				sdl.Delay(graphic.DELAY_TIME - frameTime)
			}
		}

	}
}
