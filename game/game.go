package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type Game struct {
	// start position (left top) of camera
	camPos       vector.Pos
	gra          *graphic.Graphic
	currentLevel *level.Level
	running      bool
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
	game.running = true
	for game.running {
		frameStart := sdl.GetTicks()

		events := game.handleEvent()

		// update
		for _, o := range game.currentLevel.Objects {
			o.Update(events, sdl.GetTicks())
		}

		// render
		game.gra.ClearScreen()
		game.currentLevel.Draw(game.gra, game.camPos)
		game.gra.ShowScreen()

		frameTime := sdl.GetTicks() - frameStart

		// Fixed frame rate
		if frameTime < graphic.DELAY_TIME_MS {
			sdl.Delay(graphic.DELAY_TIME_MS - frameTime)
		}
	}
	game.Quit()
}

func (game *Game) handleEvent() *intsets.Sparse {
	var events intsets.Sparse
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			game.running = true
			return nil
		case *sdl.KeyDownEvent:
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_LEFT:
				events.Insert(event.EVENT_KEYDOWN_LEFT)
			case sdl.SCANCODE_RIGHT:
				events.Insert(event.EVENT_KEYDOWN_RIGHT)
			case sdl.SCANCODE_SPACE:
				events.Insert(event.EVENT_KEYDOWN_SPACE)
			}
		}
	}
	return &events
}
