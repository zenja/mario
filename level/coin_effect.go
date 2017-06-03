package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type coinEffect struct {
	coinRes    graphic.Resource
	tileRect   sdl.Rect
	levelRect  sdl.Rect
	velocity   vector.Vec2D
	startTicks uint32
	lastTicks  uint32
	finished   bool
}

func NewCoinEffect(tid vector.TileID, resourceRegistry map[graphic.ResourceID]graphic.Resource, ticks uint32) Effect {
	coinRes, _ := resourceRegistry[graphic.RESOURCE_TYPE_COIN]
	tileRect := *GetTileRect(tid)
	return &coinEffect{
		coinRes:    coinRes,
		tileRect:   tileRect,
		levelRect:  tileRect,
		velocity:   vector.Vec2D{0, -450},
		startTicks: ticks,
		lastTicks:  ticks,
		finished:   false,
	}
}

func (ci *coinEffect) UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	// speed up
	ci.velocity.Y -= 50

	velocityStep := CalcVelocityStep(ci.velocity, ticks, ci.lastTicks, nil)
	ci.tileRect.X += velocityStep.X
	ci.tileRect.Y += velocityStep.Y

	if ticks-ci.startTicks > 100 {
		ci.finished = true
	} else {
		g.DrawResource(ci.coinRes, ci.tileRect, camPos)
	}

	ci.lastTicks = ticks
}

func (ci *coinEffect) Finished() bool {
	return ci.finished
}
