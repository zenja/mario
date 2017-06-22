package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// showOnceEffect is an Effect
var _ Effect = &showOnceEffect{}

// showOnceEffect shows a resource for a while and then disappear
type showOnceEffect struct {
	res        graphic.Resource
	levelRect  sdl.Rect // TODO only need to keep start position, width and height should use res's
	startTicks uint32
	durationMs uint32
	finished   bool
}

func NewShowOnceEffect(res graphic.Resource, levelRect sdl.Rect, ticks uint32, durationMs uint32) *showOnceEffect {
	return &showOnceEffect{
		res:        res,
		levelRect:  levelRect,
		startTicks: ticks,
		durationMs: durationMs,
		finished:   false,
	}
}

func (soe *showOnceEffect) Update(ticks uint32) {
	if ticks-soe.startTicks > soe.durationMs {
		soe.finished = true
	}
}

func (soe *showOnceEffect) Draw(camPos vector.Pos, ticks uint32) {
	if !soe.Finished() {
		graphic.DrawResource(soe.res, soe.levelRect, camPos)
	}
}

func (soe *showOnceEffect) Finished() bool {
	return soe.finished
}

func (soe *showOnceEffect) OnFinished() {
	// Do nothing
}
