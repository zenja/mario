package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/vector"
)

type Game struct {
	// start position (left top) of camera
	camPos vector.Pos

	gra          *graphic.Graphic
	currentLevel *level.Level
}

func NewGame() *Game {
	return &Game{gra: graphic.New()}
}

func (game *Game) LoadLevel(filename string) {
	game.currentLevel = level.ParseLevelFromFile(filename, game.gra.ResourceRegistry)
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
			game.gra.ClearScreen()
			game.currentLevel.Draw(game.gra, game.camPos.X, game.camPos.Y)
			game.gra.ShowScreen()

			frameTime := sdl.GetTicks() - frameStart

			if frameTime < graphic.DELAY_TIME_MS {
				sdl.Delay(graphic.DELAY_TIME_MS - frameTime)
			}
		}
	}
	game.Quit()
}
