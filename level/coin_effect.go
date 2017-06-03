package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type coinEffect struct {
	coinRes    graphic.Resource
	levelRect  sdl.Rect
	startTicks uint32
	finished   bool
}

func NewCoinEffect(tid vector.TileID, resourceRegistry map[graphic.ResourceID]graphic.Resource, ticks uint32) Effect {
	coinRes, _ := resourceRegistry[graphic.RESOURCE_TYPE_COIN]
	return &coinEffect{
		coinRes:    coinRes,
		levelRect:  *GetTileRect(tid),
		startTicks: ticks,
		finished:   false,
	}
}

func (ci *coinEffect) UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	if ticks-ci.startTicks > 700 {
		ci.finished = true
	} else {
		g.DrawResource(ci.coinRes, ci.levelRect, camPos)
	}
}

func (ci *coinEffect) Finished() bool {
	return ci.finished
}
