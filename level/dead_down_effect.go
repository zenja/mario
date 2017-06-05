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

func NewDeadDownEffect(res graphic.Resource, levelRect sdl.Rect, ticks uint32) *deadDownEffect {
	return &deadDownEffect{
		res:        res,
		levelRect:  levelRect,
		velocity:   vector.Vec2D{400, -1000},
		startTicks: ticks,
		lastTicks:  ticks,
		finished:   false,
	}
}

func (dde *deadDownEffect) UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {
	if ticks-dde.startTicks > 2000 {
		dde.finished = true
		return
	}

	gravity := vector.Vec2D{0, 50}
	dde.velocity.Add(gravity)
	velStep := CalcVelocityStep(dde.velocity, ticks, dde.lastTicks, nil)
	dde.levelRect.X += velStep.X
	dde.levelRect.Y += velStep.Y

	g.DrawResource(dde.res, dde.levelRect, camPos)

	dde.lastTicks = ticks
}

func (dde *deadDownEffect) Finished() bool {
	return dde.finished
}
