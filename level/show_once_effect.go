package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// showOnceEffect shows a resource for a while and then disappear
type showOnceEffect struct {
	res        graphic.Resource
	levelRect  sdl.Rect
	startTicks uint32
	durationMs uint32
	finished   bool
}

func NewShowOnceEffect(res graphic.Resource, levelRect sdl.Rect, ticks uint32, durationMs uint32) Effect {
	return &showOnceEffect{
		res:        res,
		levelRect:  levelRect,
		startTicks: ticks,
		durationMs: durationMs,
		finished:   false,
	}
}

func (soe *showOnceEffect) UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	if ticks-soe.startTicks > soe.durationMs {
		soe.finished = true
		return
	}

	g.DrawResource(soe.res, soe.levelRect, camPos)
}

func (soe *showOnceEffect) Finished() bool {
	return soe.finished
}
