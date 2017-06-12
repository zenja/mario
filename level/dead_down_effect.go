package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// deadDownEffect is an Effect
var _ Effect = &deadDownEffect{}

type deadDownEffect struct {
	res        graphic.Resource
	levelRect  sdl.Rect
	velocity   vector.Vec2D
	startTicks uint32
	lastTicks  uint32
	finished   bool
}

func NewDeadDownEffect(res graphic.Resource, toRight bool, levelRect sdl.Rect, ticks uint32) *deadDownEffect {
	var velocityX int32
	if toRight {
		velocityX = 400
	} else {
		velocityX = -400
	}
	return &deadDownEffect{
		res:        res,
		levelRect:  levelRect,
		velocity:   vector.Vec2D{velocityX, -1000},
		startTicks: ticks,
		lastTicks:  ticks,
		finished:   false,
	}
}

func NewStraightDeadDownEffect(res graphic.Resource, levelRect sdl.Rect, ticks uint32) *deadDownEffect {
	return &deadDownEffect{
		res:        res,
		levelRect:  levelRect,
		velocity:   vector.Vec2D{0, -1000},
		startTicks: ticks,
		lastTicks:  ticks,
		finished:   false,
	}
}

func (dde *deadDownEffect) Update(ticks uint32) {
	if ticks-dde.startTicks > 2000 {
		dde.finished = true
		return
	}

	gravity := vector.Vec2D{0, 50}
	dde.velocity.Add(gravity)
	velStep := CalcVelocityStep(dde.velocity, ticks, dde.lastTicks, nil)
	dde.levelRect.X += velStep.X
	dde.levelRect.Y += velStep.Y

	dde.lastTicks = ticks
}

func (dde *deadDownEffect) Draw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	if !dde.Finished() {
		g.DrawResource(dde.res, dde.levelRect, camPos)
	}
}

func (dde *deadDownEffect) Finished() bool {
	return dde.finished
}
