package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type bullet interface {
	Object

	GetDamage() int
	IsDead() bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// bounceAndBoomBullet
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// fireball is a bullet
var _ bullet = &bounceAndBoomBullet{}

type bounceAndBoomBullet struct {
	res0    graphic.Resource
	res1    graphic.Resource
	res2    graphic.Resource
	res3    graphic.Resource
	resBoom graphic.Resource
	currRes graphic.Resource

	maxDurationMS  uint32
	boomDurationMS uint32
	gravityY       int32

	startTicks uint32
	lastTicks  uint32
	levelRect  sdl.Rect
	velocity   vector.Vec2D
	damage     int
	isDead     bool
}

func NewBounceAndBoomBullet(
	res0, res1, res2, res3, resBoom graphic.Resource,
	heroRect sdl.Rect, toRight bool, upper bool,
	maxDurationMS uint32, boomDurationMS uint32, initVelX, initVelY, initVelYUpper int32, gravityY int32, damage int,
	ticks uint32) *bounceAndBoomBullet {

	var levelRect sdl.Rect
	if toRight {
		levelRect.X = heroRect.X + heroRect.W
	} else {
		levelRect.X = heroRect.X - res0.GetW()
	}
	levelRect.Y = heroRect.Y + heroRect.H/2 - res0.GetH()/2
	levelRect.W = res0.GetW()
	levelRect.H = res0.GetH()

	initVelocity := vector.Vec2D{initVelX, initVelY}
	if !toRight {
		initVelocity.X = -initVelocity.X
	}
	if upper {
		initVelocity.Y = initVelYUpper
	}

	return &bounceAndBoomBullet{
		res0:           res0,
		res1:           res1,
		res2:           res2,
		res3:           res3,
		resBoom:        resBoom,
		currRes:        res0,
		maxDurationMS:  maxDurationMS,
		boomDurationMS: boomDurationMS,
		gravityY:       gravityY,
		startTicks:     ticks,
		lastTicks:      ticks,
		levelRect:      levelRect,
		velocity:       initVelocity,
		damage:         damage,
	}
}

func (b *bounceAndBoomBullet) Draw(camPos vector.Pos) {
	graphic.DrawResource(b.currRes, b.levelRect, camPos)
}

func (b *bounceAndBoomBullet) Update(ticks uint32, level *Level) {
	// should not last long
	if ticks-b.startTicks > b.maxDurationMS {
		b.isDead = true
	}

	// apply gravity
	gravity := vector.Vec2D{0, b.gravityY}
	b.velocity.Add(gravity)

	maxVel := vector.Vec2D{400, 200}
	velStep := CalcVelocityStep(b.velocity, ticks, b.lastTicks, &maxVel)
	b.levelRect.X += velStep.X
	b.levelRect.Y += velStep.Y

	hitTop, hitRight, hitBottom, hitLeft, _ := level.ObstMngr.SolveCollision(&b.levelRect, SOLVE_COLLISION_NORMAL)

	// if hit top/right/left, dieDown, show boom effect
	if hitTop || hitRight || hitLeft {
		b.boom(level, ticks)
		return
	}

	// bounce if hit bottom
	if hitBottom {
		b.velocity.Y = -b.velocity.Y
	}

	// switch resources for animation
	r := ticks % 400
	switch {
	case r < 100:
		b.currRes = b.res0
	case r < 200:
		b.currRes = b.res1
	case r < 300:
		b.currRes = b.res2
	default:
		b.currRes = b.res3
	}

	// check if hit any enemies
	for _, emy := range level.Enemies {
		if emy.IsDead() {
			continue
		}

		emyRect := emy.GetRect()
		if b.levelRect.HasIntersection(&emyRect) {
			emy.hitByBullet(b, level, ticks)
			b.boom(level, ticks)
		}
	}

	// update last ticks
	b.lastTicks = ticks
}

func (b *bounceAndBoomBullet) GetRect() sdl.Rect {
	return b.levelRect
}

func (b *bounceAndBoomBullet) GetZIndex() int {
	return ZINDEX_4
}

func (b *bounceAndBoomBullet) GetDamage() int {
	return b.damage
}

func (b *bounceAndBoomBullet) IsDead() bool {
	return b.isDead
}

func (b *bounceAndBoomBullet) boom(level *Level, ticks uint32) {
	b.isDead = true
	boomStartPos := vector.Vec2D{
		X: b.levelRect.X,
		Y: b.levelRect.Y,
	}
	level.AddEffect(NewShowOnceEffect(b.resBoom, boomStartPos, ticks, b.boomDurationMS))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Fireball
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	fireballMaxDurationMS  = 2000
	fireballBoomDurationMS = 100
	fireballInitVelX       = 450
	fireballInitVelY       = 200
	fireballInitVelYUpper  = 50
	fireballGravityY       = 15
	fireballDamage         = 1
)

func NewFireball(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	ticks uint32) *bounceAndBoomBullet {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_0)
	res1 := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_1)
	res2 := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_2)
	res3 := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_3)
	resBoom := graphic.Res(graphic.RESOURCE_TYPE_FIREBALL_BOOM)

	return NewBounceAndBoomBullet(
		res0, res1, res2, res3, resBoom,
		heroRect, toRight, upper,
		fireballMaxDurationMS, fireballBoomDurationMS,
		fireballInitVelX, fireballInitVelY, fireballInitVelYUpper, fireballGravityY, fireballDamage, ticks)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Shit
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	shitMaxDurationMS  = 2000
	shitBoomDurationMS = 100
	shitInitVelX       = 450
	shitInitVelY       = 200
	shitInitVelYUpper  = 50
	shitGravityY       = 15
	shitDamage         = 8
)

func NewShit(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	ticks uint32) *bounceAndBoomBullet {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_0)
	res1 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_1)
	res2 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_2)
	res3 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_3)
	resBoom := graphic.Res(graphic.RESOURCE_TYPE_SHIT_BOOM)

	return NewBounceAndBoomBullet(
		res0, res1, res2, res3, resBoom,
		heroRect, toRight, upper,
		shitMaxDurationMS, shitBoomDurationMS,
		shitInitVelX, shitInitVelY, shitInitVelYUpper, shitGravityY, shitDamage, ticks)
}

func NewShitEx(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	initVelX, initVelY, gravityY int32,
	ticks uint32) *bounceAndBoomBullet {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_0)
	res1 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_1)
	res2 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_2)
	res3 := graphic.Res(graphic.RESOURCE_TYPE_SHIT_3)
	resBoom := graphic.Res(graphic.RESOURCE_TYPE_SHIT_BOOM)

	return NewBounceAndBoomBullet(
		res0, res1, res2, res3, resBoom,
		heroRect, toRight, upper,
		shitMaxDurationMS, shitBoomDurationMS,
		initVelX, initVelY, shitInitVelYUpper, gravityY, shitDamage, ticks)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Bug
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	bugMaxDurationMS  = 2000
	bugBoomDurationMS = 100
	bugInitVelX       = 450
	bugInitVelY       = 200
	bugInitVelYUpper  = 50
	bugGravityY       = 15
	bugDamage         = 4
)

func NewBug(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	ticks uint32) *bounceAndBoomBullet {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_BUG_0)
	res1 := graphic.Res(graphic.RESOURCE_TYPE_BUG_1)
	res2 := graphic.Res(graphic.RESOURCE_TYPE_BUG_2)
	res3 := graphic.Res(graphic.RESOURCE_TYPE_BUG_3)
	resBoom := graphic.Res(graphic.RESOURCE_TYPE_BUG_BOOM)

	return NewBounceAndBoomBullet(
		res0, res1, res2, res3, resBoom,
		heroRect, toRight, upper,
		bugMaxDurationMS, bugBoomDurationMS, bugInitVelX, bugInitVelY, bugInitVelYUpper, bugGravityY, bugDamage, ticks)
}

func NewBugEx(
	heroRect sdl.Rect,
	toRight bool,
	upper bool,
	initVelX, initVelY, gravityY int32,
	ticks uint32) *bounceAndBoomBullet {

	res0 := graphic.Res(graphic.RESOURCE_TYPE_BUG_0)
	res1 := graphic.Res(graphic.RESOURCE_TYPE_BUG_1)
	res2 := graphic.Res(graphic.RESOURCE_TYPE_BUG_2)
	res3 := graphic.Res(graphic.RESOURCE_TYPE_BUG_3)
	resBoom := graphic.Res(graphic.RESOURCE_TYPE_BUG_BOOM)

	return NewBounceAndBoomBullet(
		res0, res1, res2, res3, resBoom,
		heroRect, toRight, upper,
		bugMaxDurationMS, bugBoomDurationMS, initVelX, initVelY, bugInitVelYUpper, gravityY, bugDamage, ticks)
}
