package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type volatileObject interface {
	Object

	IsDead() bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Fireball
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// fireball is a volatileObject
var _ volatileObject = &fireball{}

const (
	fireballMaxDurationMS  = 2000
	fireballBoomDurationMS = 100
	fireballInitVelX       = 450
	fireballInitVelY       = 200
	fireballInitVelYUpper  = 50
	fireballGravityY       = 15
)

type fireball struct {
	res0    graphic.Resource
	res1    graphic.Resource
	res2    graphic.Resource
	res3    graphic.Resource
	resBoom graphic.Resource
	currRes graphic.Resource

	startTicks uint32
	lastTicks  uint32
	levelRect  sdl.Rect
	velocity   vector.Vec2D
	isDead     bool
}

func NewFireball(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	ticks uint32) *fireball {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_0)

	var levelRect sdl.Rect
	if toRight {
		levelRect.X = heroRect.X + heroRect.W
	} else {
		levelRect.X = heroRect.X - res0.GetW()
	}
	levelRect.Y = heroRect.Y + heroRect.H/2 - res0.GetH()/2
	levelRect.W = res0.GetW()
	levelRect.H = res0.GetH()

	initVelocity := vector.Vec2D{fireballInitVelX, fireballInitVelY}
	if !toRight {
		initVelocity.X = -initVelocity.X
	}
	if upper {
		initVelocity.Y = fireballInitVelYUpper
	}

	return &fireball{
		res0:       res0,
		res1:       graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_1),
		res2:       graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_2),
		res3:       graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_3),
		resBoom:    graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_BOOM),
		currRes:    res0,
		startTicks: ticks,
		lastTicks:  ticks,
		levelRect:  levelRect,
		velocity:   initVelocity,
	}
}

func (f *fireball) Draw(camPos vector.Pos) {
	graphic.DrawResource(f.currRes, f.levelRect, camPos)
}

func (f *fireball) Update(ticks uint32, level *Level) {
	// a fireball should not last long
	if ticks-f.startTicks > fireballMaxDurationMS {
		f.isDead = true
	}

	// apply gravity
	gravity := vector.Vec2D{0, fireballGravityY}
	f.velocity.Add(gravity)

	maxVel := vector.Vec2D{400, 200}
	velStep := CalcVelocityStep(f.velocity, ticks, f.lastTicks, &maxVel)
	f.levelRect.X += velStep.X
	f.levelRect.Y += velStep.Y

	hitTop, hitRight, hitBottom, hitLeft, _ := level.ObstMngr.SolveCollision(&f.levelRect)

	// if hit top/right/left, dieDown, show boom effect
	if hitTop || hitRight || hitLeft {
		f.boom(level, ticks)
		return
	}

	// bounce if hit bottom
	if hitBottom {
		f.velocity.Y = -f.velocity.Y
	}

	// switch resources for animation
	r := ticks % 200
	switch {
	case r < 50:
		f.currRes = f.res0
	case r < 100:
		f.currRes = f.res1
	case r < 150:
		f.currRes = f.res2
	default:
		f.currRes = f.res3
	}

	// check if hit any enemies
	for _, emy := range level.Enemies {
		if emy.IsDead() {
			continue
		}

		emyRect := emy.GetRect()
		if f.levelRect.HasIntersection(&emyRect) {
			emy.hitByFireball(f, level, ticks)
			f.boom(level, ticks)
		}
	}

	// update last ticks
	f.lastTicks = ticks
}

func (f *fireball) GetRect() sdl.Rect {
	return f.levelRect
}

func (f *fireball) GetZIndex() int {
	return ZINDEX_4
}

func (f *fireball) IsDead() bool {
	return f.isDead
}

func (f *fireball) boom(level *Level, ticks uint32) {
	f.isDead = true
	boomStartPos := vector.Vec2D{
		X: f.levelRect.X,
		Y: f.levelRect.Y,
	}
	level.AddEffect(NewShowOnceEffect(f.resBoom, boomStartPos, ticks, fireballBoomDurationMS))
}
