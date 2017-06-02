package game

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/overlay"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type Game struct {
	// start position (left top) of camera
	camPos       vector.Pos
	gra          *graphic.Graphic
	currentLevel *level.Level
	running      bool
	overlays     []overlay.Overlay
}

func NewGame() *Game {
	gra := graphic.New()

	// register overlays
	var overlays []overlay.Overlay
	overlays = append(overlays, &overlay.FPSOverlay{})

	return &Game{
		gra:      gra,
		overlays: overlays,
	}
}

func (game *Game) LoadLevel(filename string) {
	game.currentLevel = level.ParseLevelFromFile(filename, game.gra.ResourceRegistry)
}

func (game *Game) Quit() {
	game.gra.DestroyAndQuit()
}

func (game *Game) StartGameLoop() {
	// this may prevent window not responding
	runtime.LockOSThread()

	game.running = true
	for game.running {
		frameStart := sdl.GetTicks()

		events := game.handleEvents()

		// update
		for _, o := range game.currentLevel.Objects {
			o.Update(events, sdl.GetTicks(), game.currentLevel)
		}

		// update camera position
		game.updateCamPos()

		// start render
		game.gra.ClearScreen()
		game.currentLevel.Draw(game.gra, game.camPos)

		// render overlays
		for _, ol := range game.overlays {
			ol.Draw(game.gra, sdl.GetTicks())
		}

		// show screen
		game.gra.ShowScreen()

		frameTime := sdl.GetTicks() - frameStart

		// Fixed frame rate
		if frameTime < graphic.DELAY_TIME_MS {
			sdl.Delay(graphic.DELAY_TIME_MS - frameTime)
		}
	}
	game.Quit()
}

func (game *Game) handleEvents() *intsets.Sparse {
	var events intsets.Sparse
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch e.(type) {
		case *sdl.QuitEvent:
			game.running = false
			return nil
		}
	}
	kbState := sdl.GetKeyboardState()
	if kbState[int(sdl.SCANCODE_LEFT)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_LEFT))
	}
	if kbState[int(sdl.SCANCODE_RIGHT)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_RIGHT))
	}
	if kbState[int(sdl.SCANCODE_SPACE)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_SPACE))
	}
	return &events
}

// updateCamPos update the position of camera based on hero's position
// It tries to put hero center in vertical, 2/3 camera height from top,
// but when that exceeds level boundary, it will respect level boundary
func (game *Game) updateCamPos() {
	heroRect := game.currentLevel.Hero.GetRect()
	perfectX := heroRect.X - (graphic.SCREEN_WIDTH-heroRect.W)/2
	perfectY := heroRect.Y - (graphic.SCREEN_HEIGHT-heroRect.H)*2/3
	game.camPos.X = perfectX
	game.camPos.Y = perfectY
	// check left
	if perfectX < 0 {
		game.camPos.X = 0
	}
	// check top
	if perfectY < 0 {
		game.camPos.Y = 0
	}
	// check right
	if perfectX+graphic.SCREEN_WIDTH > game.currentLevel.GetLevelWidth() {
		game.camPos.X = game.currentLevel.GetLevelWidth() - graphic.SCREEN_WIDTH
	}
	// check bottom
	if perfectY+graphic.SCREEN_HEIGHT > game.currentLevel.GetLevelHeight() {
		game.camPos.Y = game.currentLevel.GetLevelHeight() - graphic.SCREEN_HEIGHT
	}
}
