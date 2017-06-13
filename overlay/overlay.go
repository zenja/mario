package overlay

import (
	"fmt"

	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/vector"
)

type Overlay interface {
	Draw(h *level.Hero, ticks uint32)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// FPSOverlay
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FPSOverlay struct {
	currentTicks uint32
}

func (fo *FPSOverlay) Draw(h *level.Hero, ticks uint32) {
	var fps uint32
	pos := vector.Pos{50, 50}
	color := sdl.Color{255, 255, 255, 0}
	if fo.currentTicks == 0 {
		fo.currentTicks = ticks
		return
	}
	if ticks-fo.currentTicks <= 0 {
		log.Printf("FPSOverlay: strange, ticks (%d) <= fo.currentTicks (%d)", ticks, fo.currentTicks)
		graphic.DrawText(fmt.Sprint("FPS: NaN"), pos, color)
		return
	}
	fps = 1000 / (ticks - fo.currentTicks)
	graphic.DrawText(fmt.Sprintf("FPS: %d", fps), pos, color)
	fo.currentTicks = ticks
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HeroLiveOverlay
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type HeroLiveOverlay struct{}

func (hlo *HeroLiveOverlay) Draw(h *level.Hero, ticks uint32) {
	pos := vector.Pos{graphic.SCREEN_WIDTH - 150, 50}
	color := sdl.Color{255, 255, 255, 0}
	graphic.DrawText(fmt.Sprintf("Lives: %d", h.GetLives()), pos, color)
}
