package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type hittableByHero interface {
	hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32)
}

type hittableByBrokenTile interface {
	hitByBrokenTile(level *Level, ticks uint32)
}

type hittableByFireball interface {
	hitByFireball(fb *fireball, level *Level, ticks uint32)
}

type Enemy interface {
	// Enemy is an object
	Object

	// Enemy is hittable by hero
	hittableByHero

	// Enemy is hittable by breaking tile (from enemy's bottom tile)
	hittableByBrokenTile

	// Enemy is hittable by fireball
	hittableByFireball

	// if the enemy is dead, if so, don't need to update/draw
	IsDead() bool

	// after Kill(), IsDead() should return true
	Kill()
}

type basicEnemy struct {
	isDead bool
}

func (be *basicEnemy) IsDead() bool {
	return be.isDead
}

func (be *basicEnemy) Kill() {
	be.isDead = true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MushroomEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type mushroomEnemy struct {
	basicEnemy

	res0      graphic.Resource
	res1      graphic.Resource
	resHit    graphic.Resource
	resDown   graphic.Resource
	currRes   graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
}

func NewMushroomEnemy(startPos vector.Pos) *mushroomEnemy {
	res0 := graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_0)
	return &mushroomEnemy{
		res0:      res0,
		res1:      graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_1),
		resHit:    graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_HIT),
		resDown:   graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_DOWN),
		currRes:   res0,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res0.GetW(), res0.GetH()},
		velocity:  vector.Vec2D{100, 0},
	}
}

func (m *mushroomEnemy) GetRect() sdl.Rect {
	return m.levelRect
}

func (m *mushroomEnemy) GetZIndex() int {
	return ZINDEX_4
}

func (m *mushroomEnemy) Update(ticks uint32, level *Level) {
	if m.lastTicks == 0 {
		m.lastTicks = ticks
		return
	}

	enemySimpleMove(ticks, m.lastTicks, &m.velocity, &m.levelRect, level)

	m.updateResource(ticks)

	m.lastTicks = ticks
}

func (m *mushroomEnemy) Draw(camPos vector.Pos) {
	graphic.DrawResource(m.currRes, m.levelRect, camPos)
}

func (m *mushroomEnemy) updateResource(ticks uint32) {
	if ticks%1000 < 500 {
		m.currRes = m.res0
	} else {
		m.currRes = m.res1
	}
}

func (m *mushroomEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if direction == HIT_FROM_TOP_W_INTENT {
		// dead immediately!
		m.isDead = true

		// bounce the hero up
		h.velocity.Y = -1200

		// add dead effect
		level.AddEffect(NewShowOnceEffect(m.resHit, m.levelRect, ticks, 500))
	} else {
		// hero is hurt
		h.Hurt()
	}
}

func (m *mushroomEnemy) hitByBrokenTile(level *Level, ticks uint32) {
	m.dieDown(true, level, ticks)
}

func (m *mushroomEnemy) hitByFireball(fb *fireball, level *Level, ticks uint32) {
	var dieToRight bool
	if fb.levelRect.X < m.levelRect.X {
		dieToRight = true
	}
	m.dieDown(dieToRight, level, ticks)
}

func (m *mushroomEnemy) dieDown(toRight bool, level *Level, ticks uint32) {
	m.isDead = true
	level.AddEffect(NewDeadDownEffect(m.resDown, toRight, m.levelRect, ticks))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TortoiseEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	tortoiseInitVelocityX         = 50
	tortoiseBumpingVelocityXRight = 800
)

type tortoiseEnemy struct {
	basicEnemy

	resLeft0      graphic.Resource
	resLeft1      graphic.Resource
	resRight0     graphic.Resource
	resRight1     graphic.Resource
	resSemiInside graphic.Resource
	resInside     graphic.Resource
	currRes       graphic.Resource

	isFacingRight bool
	levelRect     sdl.Rect
	velocity      vector.Vec2D
	lastTicks     uint32

	insideStartTicks uint32 // when tortoise go inside
	bumpStartTicks   uint32 // when tortoise start bumping
}

func NewTortoiseEnemy(startPos vector.Pos) Enemy {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_LEFT_0)
	return &tortoiseEnemy{
		resLeft0:      resLeft0,
		resLeft1:      graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_LEFT_1),
		resRight0:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RIGHT_0),
		resRight1:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RIGHT_1),
		resSemiInside: graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_SEMI_INSIDE),
		resInside:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_INSIDE),
		currRes:       resLeft0,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{-100, 0},
	}
}

func (t *tortoiseEnemy) GetRect() sdl.Rect {
	return t.levelRect
}

func (t *tortoiseEnemy) GetZIndex() int {
	return ZINDEX_4
}

func (t *tortoiseEnemy) Update(ticks uint32, level *Level) {
	if t.lastTicks == 0 {
		t.lastTicks = ticks
		return
	}

	onHitLeft := func() { t.isFacingRight = true }
	onHitRight := func() { t.isFacingRight = false }
	enemySimpleMoveEx(ticks, t.lastTicks, &t.velocity, &t.levelRect, level, onHitLeft, onHitRight)

	t.updateResource(ticks)

	t.lastTicks = ticks
}

func (t *tortoiseEnemy) Draw(camPos vector.Pos) {
	graphic.DrawResource(t.currRes, t.levelRect, camPos)
}

func (t *tortoiseEnemy) updateResource(ticks uint32) {
	if t.insideStartTicks > 0 || t.bumpStartTicks > 0 {
		t.currRes = t.resInside
		return
	}

	if ticks%1000 < 500 {
		if t.isFacingRight {
			t.currRes = t.resRight0
		} else {
			t.currRes = t.resLeft0
		}
	} else {
		if t.isFacingRight {
			t.currRes = t.resRight1
		} else {
			t.currRes = t.resLeft1
		}
	}
}

func (t *tortoiseEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if t.insideStartTicks > 0 && t.bumpStartTicks > 0 {
		log.Fatal("bug! insideStartTicks and bumpStartTicks cannot be positive at the same time!")
	}

	switch direction {
	case HIT_FROM_TOP_W_INTENT:
		// bounce the hero up
		h.velocity.Y = -1200

		switch {
		// case 1: normal state => go inside, don't move in X
		case t.insideStartTicks == 0 && t.bumpStartTicks == 0:
			t.toInsideState(ticks)

		// case 2: inside state => start bumping
		case t.bumpStartTicks == 0:
			// decide move right or left
			heroMidX := h.levelRect.X + h.levelRect.W/2
			tortoiseMidX := t.levelRect.X + t.levelRect.W/2
			if heroMidX < tortoiseMidX {
				t.toBumpingState(ticks, true)
			} else {
				t.toBumpingState(ticks, false)
			}

		// case 3: bumping state => stop bumping, turn to inside state, don't move in X
		default:
			t.toInsideState(ticks)
		}

	case HIT_FROM_LEFT_W_INTENT:
		if t.insideStartTicks > 0 {
			t.toBumpingState(ticks, true)
		} else {
			h.Hurt()
		}

	case HIT_FROM_RIGHT_W_INTENT:
		if t.insideStartTicks > 0 {
			t.toBumpingState(ticks, false)
		} else {
			h.Hurt()
		}

	default:
		// hero is hurt
		h.Hurt()
	}
}

func (t *tortoiseEnemy) hitByBrokenTile(level *Level, ticks uint32) {
	t.dieDown(true, level, ticks)
}

func (t *tortoiseEnemy) hitByFireball(fb *fireball, level *Level, ticks uint32) {
	var dieToRight bool
	if fb.levelRect.X < t.levelRect.X {
		dieToRight = true
	}
	t.dieDown(dieToRight, level, ticks)
}

func (t *tortoiseEnemy) dieDown(toRight bool, level *Level, ticks uint32) {
	t.isDead = true
	level.AddEffect(NewDeadDownEffect(t.resInside, toRight, t.levelRect, ticks))
}

func (t *tortoiseEnemy) toInsideState(ticks uint32) {
	// change state
	t.insideStartTicks = ticks
	t.bumpStartTicks = 0

	t.velocity.X = 0
}

func (t *tortoiseEnemy) toBumpingState(ticks uint32, toRight bool) {
	// change state
	t.insideStartTicks = 0
	t.bumpStartTicks = ticks

	if toRight {
		// move right
		t.velocity.X = tortoiseBumpingVelocityXRight
	} else {
		// move left
		t.velocity.X = -tortoiseBumpingVelocityXRight
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GoodMushroom
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type goodMushroom struct {
	basicEnemy

	res       graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
}

func NewGoodMushroom(startPos vector.Pos) *goodMushroom {
	res := graphic.Res(graphic.RESOURCE_TYPE_GOOD_MUSHROOM)
	return &goodMushroom{
		res:       res,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res.GetW(), res.GetH()},
		velocity:  vector.Vec2D{-100, -500},
	}
}

func (gm *goodMushroom) GetRect() sdl.Rect {
	return gm.levelRect
}

func (gm *goodMushroom) GetZIndex() int {
	return ZINDEX_4
}

func (gm *goodMushroom) Update(ticks uint32, level *Level) {
	if gm.lastTicks == 0 {
		gm.lastTicks = ticks
		return
	}

	enemySimpleMove(ticks, gm.lastTicks, &gm.velocity, &gm.levelRect, level)

	gm.lastTicks = ticks
}

func (gm *goodMushroom) Draw(camPos vector.Pos) {
	graphic.DrawResource(gm.res, gm.levelRect, camPos)
}

func (gm *goodMushroom) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	gm.isDead = true
	// upgrade hero
	h.upgrade(level)
}

func (gm *goodMushroom) hitByBrokenTile(level *Level, ticks uint32) {
	// bounce up
	gm.velocity.Y -= 500
}

func (gm *goodMushroom) hitByFireball(fb *fireball, level *Level, ticks uint32) {
	// No interaction with fireball; Do nothing
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func enemySimpleMove(
	ticks uint32,
	lastTicks uint32,
	vel *vector.Vec2D,
	levelRect *sdl.Rect,
	level *Level) {

	enemySimpleMoveEx(ticks, lastTicks, vel, levelRect, level, nil, nil)
}

func enemySimpleMoveEx(
	ticks uint32,
	lastTicks uint32,
	vel *vector.Vec2D,
	levelRect *sdl.Rect,
	level *Level,
	onHitLeft func(),
	onHitRight func()) {

	gravity := vector.Vec2D{0, 50}
	vel.Add(gravity)

	maxVel := vector.Vec2D{int32(graphic.TILE_SIZE * 30 / 100), int32(graphic.TILE_SIZE * 30 / 100)}
	velocityStep := CalcVelocityStep(*vel, ticks, lastTicks, &maxVel)
	levelRect.X += velocityStep.X
	levelRect.Y += velocityStep.Y

	_, hitRight, hitBottom, hitLeft, _ := level.EnemyObstMngr.SolveCollision(levelRect)

	if hitRight {
		vel.X = -vel.X
		if onHitRight != nil {
			onHitRight()
		}
	}
	if hitLeft {
		vel.X = -vel.X
		if onHitLeft != nil {
			onHitLeft()
		}
	}

	// prevent too big down velocity
	if velocityStep.Y > 0 && hitBottom {
		vel.Y = 0
	}
}
