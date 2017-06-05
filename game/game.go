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
	Gra *graphic.Graphic

	// start position (left top) of camera
	camPos       vector.Pos
	currentLevel *level.Level
	running      bool
	overlays     []overlay.Overlay
}

func NewGame() *Game {
	gra := graphic.New()

	// register overlays
	var overlays []overlay.Overlay
	overlays = append(overlays, &overlay.FPSOverlay{})
	overlays = append(overlays, &overlay.HeroLiveOverlay{})

	return &Game{
		Gra:      gra,
		overlays: overlays,
	}
}

func (game *Game) LoadLevel(filename string) {
	game.currentLevel = level.ParseLevelFromFile(filename, game.Gra.ResourceRegistry)
}

func (game *Game) Quit() {
	game.Gra.DestroyAndQuit()
}

func (game *Game) StartGameLoop() {
	// this may prevent window not responding
	runtime.LockOSThread()

	game.running = true
	for game.running {
		frameStart := sdl.GetTicks()

		events := game.handleEvents()

		// update tile objects
		for i := 0; i < int(game.currentLevel.NumTiles.X); i++ {
			for j := 0; j < int(game.currentLevel.NumTiles.Y); j++ {
				o := game.currentLevel.TileObjects[i][j]
				if o == nil {
					continue
				}
				o.Update(events, sdl.GetTicks(), game.currentLevel)
			}
		}

		// update non-tile objects
		game.currentLevel.TheHero.Update(events, sdl.GetTicks(), game.currentLevel)

		// update camera position
		game.updateCamPos()

		// start render
		game.Gra.ClearScreenWithColor(game.currentLevel.BGColor)
		game.currentLevel.UpdateAndDraw(game.Gra, game.camPos)

		// render overlays
		for _, ol := range game.overlays {
			ol.Draw(game.Gra, game.currentLevel.TheHero, sdl.GetTicks())
		}

		// show screen
		game.Gra.ShowScreen()

		frameTime := sdl.GetTicks() - frameStart

		// Fixed frame rate
		if frameTime < graphic.DELAY_TIME_MS {
			sdl.Delay(graphic.DELAY_TIME_MS - frameTime)
		}
	}
	game.Quit()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
	if kbState[int(sdl.SCANCODE_UP)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_UP))
	}
	if kbState[int(sdl.SCANCODE_SPACE)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_SPACE))
	}
	if kbState[int(sdl.SCANCODE_F)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_F))
	}
	return &events
}

// updateCamPos update the position of camera based on hero's position
// It tries to put hero center in vertical, 2/3 camera height from top,
// but when that exceeds level boundary, it will respect level boundary
func (game *Game) updateCamPos() {
	heroRect := game.currentLevel.TheHero.GetRect()
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
