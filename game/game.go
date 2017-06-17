package game

import (
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/overlay"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

const (
	first_level_name = "level-0"
	level_dir        = "assets/levels"
)

type Game struct {
	// start position (left top) of camera
	camPos       vector.Pos
	levelSpecs   map[string]*level.LevelSpec
	currentLevel *level.Level
	running      bool
	overlays     []overlay.Overlay
}

func NewGame() *Game {
	// register overlays
	var overlays []overlay.Overlay
	overlays = append(overlays, &overlay.FPSOverlay{})
	overlays = append(overlays, &overlay.HeroLiveOverlay{})

	return &Game{
		levelSpecs: make(map[string]*level.LevelSpec),
		overlays:   overlays,
	}
}

func (game *Game) Init() {
	game.loadLevels()
	game.currentLevel.Init()
}

func (game *Game) Quit() {
	graphic.DestroyAndQuit()
}

func (game *Game) StartGameLoop() {
	game.running = true
	for game.running {
		frameStart := sdl.GetTicks()

		events := game.gatherEvents()

		// game event handling
		game.handleGlobalEvents(events)

		// level event handling
		game.currentLevel.HandleEvents(events)

		// update current level
		game.currentLevel.Update(events, sdl.GetTicks())

		// update camera position
		game.updateCamPos()

		// start render
		graphic.ClearScreenWithColor(game.currentLevel.BGColor)

		// render current level
		game.currentLevel.Draw(game.camPos, sdl.GetTicks())

		// render overlays
		for _, ol := range game.overlays {
			ol.Draw(game.currentLevel.TheHero, sdl.GetTicks())
		}

		// show screen
		graphic.ShowScreen()

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

func (game *Game) gatherEvents() *intsets.Sparse {
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
	if kbState[int(sdl.SCANCODE_F1)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_F1))
	}
	if kbState[int(sdl.SCANCODE_F2)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_F2))
	}
	if kbState[int(sdl.SCANCODE_F3)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_F3))
	}
	if kbState[int(sdl.SCANCODE_F4)] == 1 {
		events.Insert(int(event.EVENT_KEYDOWN_F4))
	}
	return &events
}

// updateCamPos update the position of camera based on hero's position
// It tries to put hero center in vertical & top,
// but when that exceeds level boundary, it will respect level boundary
func (game *Game) updateCamPos() {
	heroRect := game.currentLevel.TheHero.GetRect()
	perfectX := heroRect.X - (graphic.SCREEN_WIDTH-heroRect.W)/2
	perfectY := heroRect.Y - (graphic.SCREEN_HEIGHT-heroRect.H)/2
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

func (game *Game) handleGlobalEvents(events *intsets.Sparse) {
	if events.Has(int(event.EVENT_KEYDOWN_F1)) {
		game.currentLevel.Restart()
	}
}

func (game *Game) loadLevels() {
	fileInfos, err := ioutil.ReadDir(level_dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if info.IsDir() {
			continue
		}
		spec := level.ParseLevelSpec(level_dir + "/" + info.Name())
		game.levelSpecs[spec.Name] = spec
	}

	// check if the next level of each actually exists
	for name, spec := range game.levelSpecs {
		_, ok := game.levelSpecs[spec.NextLevelName]
		if !ok {
			log.Fatalf("%s's next level is %s, but not found", name, spec.NextLevelName)
		}
	}

	firstLevel, ok := game.levelSpecs[first_level_name]
	if !ok {
		log.Fatalf("level not found: %s", first_level_name)
	}
	game.currentLevel = level.BuildLevel(firstLevel)
}
