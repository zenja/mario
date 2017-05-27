package overlay

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Overlay interface {
	Draw(g *graphic.Graphic, ticks uint32)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// FPSOverlay
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FPSOverlay struct {
	currentTicks uint32
}

func (fo *FPSOverlay) Draw(g *graphic.Graphic, ticks uint32) {
	var fps uint32
	if fo.currentTicks == 0 {
		fo.currentTicks = ticks
		return
	}
	fps = 1000 / (ticks - fo.currentTicks)
	g.DrawText(fmt.Sprintf("FPS: %d", fps), vector.Pos{50, 50}, sdl.Color{255, 255, 255, 0})
	fo.currentTicks = ticks
}
