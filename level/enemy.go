package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Enemy interface {
	// Enemy is hittable by hero
	heroHittableObject

	GetLevelRect() sdl.Rect

	// if the enemy is dead, if so, don't need to update/draw
	IsDead() bool

	Update(ticks uint32, level *Level)

	Draw(g *graphic.Graphic, camPos vector.Pos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MushroomEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type mushroomEnemy struct {
	res0      graphic.Resource
	res1      graphic.Resource
	resHit    graphic.Resource
	currRes   graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
	isHit     bool
	isDead    bool
	hitTicks  uint32 // ticks when the enemy is hit
}

func NewMushroomEnemy(startPos vector.Pos, resourceRegistry map[graphic.ResourceID]graphic.Resource) Enemy {
	res0 := resourceRegistry[graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_0]
	return &mushroomEnemy{
		res0:      res0,
		res1:      resourceRegistry[graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_1],
		resHit:    resourceRegistry[graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_HIT],
		currRes:   res0,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res0.GetW(), res0.GetH()},
		velocity:  vector.Vec2D{100, 0},
	}
}

func (m *mushroomEnemy) GetLevelRect() sdl.Rect {
	return m.levelRect
}

func (m *mushroomEnemy) Update(ticks uint32, level *Level) {
	if m.lastTicks == 0 {
		m.lastTicks = ticks
		return
	}

	gravity := vector.Vec2D{0, 50}
	m.velocity.Add(gravity)

	maxVel := vector.Vec2D{int32(graphic.TILE_SIZE * 30 / 100), int32(graphic.TILE_SIZE * 30 / 100)}
	velocityStep := CalcVelocityStep(m.velocity, ticks, m.lastTicks, &maxVel)
	m.levelRect.X += velocityStep.X
	m.levelRect.Y += velocityStep.Y

	_, hitRight, hitBottom, hitLeft, _ := level.ObstMngr.SolveCollision(&m.levelRect)

	if hitRight {
		m.velocity.X = -100
	}
	if hitLeft {
		m.velocity.X = 100
	}

	// prevent too big down velocity
	if velocityStep.Y > 0 && hitBottom {
		m.velocity.Y = 0
	}

	m.updateResource(ticks)

	// check if should kill the enemy (wait for hit animation to finish)
	if m.hitTicks > 0 && ticks-m.hitTicks > 500 {
		m.isDead = true
	}

	m.lastTicks = ticks
}

func (m *mushroomEnemy) Draw(g *graphic.Graphic, camPos vector.Pos) {
	g.DrawResource(m.currRes, m.levelRect, camPos)
}

func (m *mushroomEnemy) IsDead() bool {
	return m.isDead
}

func (m *mushroomEnemy) updateResource(ticks uint32) {
	if m.hitTicks != 0 {
		m.currRes = m.resHit
	} else {
		if ticks%1000 < 500 {
			m.currRes = m.res0
		} else {
			m.currRes = m.res1
		}
	}
}

func (m *mushroomEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if m.hitTicks > 0 {
		// if already hit, it will lost interaction with hero
		return
	}

	if direction == HIT_FROM_TOP_W_INTENT {
		// ignore if already hit
		if m.hitTicks > 0 {
			return
		}

		// mark hit by setting hit time
		m.hitTicks = ticks

		// bounce the hero up
		h.velocity.Y = -1200
	} else {
		// hero is hurt
		h.Hurt()
	}
}
