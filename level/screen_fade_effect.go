package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Effect = &screenFadeEffect{}

type screenFadeEffect struct {
	res        graphic.Resource
	fadeIn     bool
	durationMs uint32
	startTicks uint32
	finished   bool
}

func NewScreenFadeEffect(resourceRegistry map[graphic.ResourceID]graphic.Resource, fadeIn bool, durationMs uint32, ticks uint32) *screenFadeEffect {
	return &screenFadeEffect{
		res:        resourceRegistry[graphic.RESOURCE_TYPE_BLACK_SCREEN],
		fadeIn:     fadeIn,
		durationMs: durationMs,
		startTicks: ticks,
	}
}

func (sfe *screenFadeEffect) Update(ticks uint32) {
	if ticks-sfe.startTicks >= sfe.durationMs {
		sfe.finished = true
		return
	}

	ratio := float64(ticks-sfe.startTicks) / float64(sfe.durationMs)
	if sfe.fadeIn {
		sfe.res.SetResourceAlpha(sdl.ALPHA_OPAQUE - uint8(float64(sdl.ALPHA_OPAQUE)*ratio))
	} else {
		sfe.res.SetResourceAlpha(uint8(float64(sdl.ALPHA_OPAQUE) * ratio))
	}
}

func (sfe *screenFadeEffect) Draw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	g.DrawResource(sfe.res, sdl.Rect{0, 0, graphic.SCREEN_WIDTH, graphic.SCREEN_HEIGHT}, vector.Pos{0, 0})
}

func (sfe *screenFadeEffect) Finished() bool {
	return sfe.finished
}
