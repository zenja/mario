package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

const hurtAnimationMS = 2000

// assert &Hero is an Object
var _ Object = &Hero{}

type Hero struct {
	// hero 0 res
	res0StandRight   graphic.Resource
	res0WalkingRight graphic.Resource
	res0StandLeft    graphic.Resource
	res0WalkingLeft  graphic.Resource

	// hero 1 res
	res1StandRight   graphic.Resource
	res1WalkingRight graphic.Resource
	res1StandLeft    graphic.Resource
	res1WalkingLeft  graphic.Resource

	// hero 2 res
	res2StandRight   graphic.Resource
	res2WalkingRight graphic.Resource
	res2StandLeft    graphic.Resource
	res2WalkingLeft  graphic.Resource

	// current set of resource
	currResStandRight   graphic.Resource
	currResWalkingRight graphic.Resource
	currResStandLeft    graphic.Resource
	currResWalkingLeft  graphic.Resource

	// current resource
	currRes graphic.Resource

	// hero's level: 0, 1, 2
	grade int

	// current rect in level, it is a hit box, not a render box
	// for hero its hit box and render box are different, hit box is smaller
	levelRect sdl.Rect

	// render box
	// for X direction, shrink area includes both left and right
	renderBoxWExpandRatio float64
	// for Y direction, shrink area includes top but not bottom, to make sure correct standing
	renderBoxHExpandRatio float64
	renderBoxW            int32
	renderBoxH            int32

	// event state
	upPressed bool
	fPressed  bool

	// current velocity, unit is pixels per second
	velocity vector.Vec2D

	lastTicks uint32

	lastFireTicks uint32

	isOnGround bool

	isFacingRight bool

	lives int

	// a non-zero hurtStartTicks means it has been hurt before and the "hurt" (or "super") status is still in effect
	// when hero got hurt, it will be set to current ticks
	// will be reset after a while
	hurtStartTicks uint32
}

func NewHero(
	renderBoxStartPos vector.Pos,
	renderBoxWExpandRatio, renderBoxHExpandRatio float64,
	resourceRegistry map[graphic.ResourceID]graphic.Resource) *Hero {

	if renderBoxWExpandRatio <= -1 || renderBoxWExpandRatio >= 1 {
		log.Fatalf("render box X expand ratio should be (-1, 1) but was %f", renderBoxWExpandRatio)
	}
	if renderBoxHExpandRatio <= -1 || renderBoxHExpandRatio >= 1 {
		log.Fatalf("render box Y expand ratio should be (-1, 1) but was %f", renderBoxHExpandRatio)
	}

	res0StandRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_0_STAND_RIGHT]
	res0WalkingRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_0_WALKING_RIGHT]
	res0StandLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_0_STAND_LEFT]
	res0WalkingLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_0_WALKING_LEFT]

	res1StandRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_1_STAND_RIGHT]
	res1WalkingRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_1_WALKING_RIGHT]
	res1StandLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_1_STAND_LEFT]
	res1WalkingLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_1_WALKING_LEFT]

	res2StandRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_2_STAND_RIGHT]
	res2WalkingRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_2_WALKING_RIGHT]
	res2StandLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_2_STAND_LEFT]
	res2WalkingLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_2_WALKING_LEFT]

	resX := renderBoxStartPos.X
	resY := renderBoxStartPos.Y
	resW := res0StandLeft.GetW()
	resH := res0StandLeft.GetH()

	hitBox := sdl.Rect{
		resX + int32(float64(resW)*renderBoxWExpandRatio/2),
		resY + int32(float64(resH)*renderBoxHExpandRatio),
		resW - int32(float64(resW)*renderBoxWExpandRatio),
		resH - int32(float64(resH)*renderBoxHExpandRatio),
	}

	h := &Hero{
		res0StandRight:   res0StandRight,
		res0WalkingRight: res0WalkingRight,
		res0StandLeft:    res0StandLeft,
		res0WalkingLeft:  res0WalkingLeft,
		res1StandRight:   res1StandRight,
		res1WalkingRight: res1WalkingRight,
		res1StandLeft:    res1StandLeft,
		res1WalkingLeft:  res1WalkingLeft,
		res2StandRight:   res2StandRight,
		res2WalkingRight: res2WalkingRight,
		res2StandLeft:    res2StandLeft,
		res2WalkingLeft:  res2WalkingLeft,

		currResStandRight:   res0StandRight,
		currResWalkingRight: res0WalkingRight,
		currResStandLeft:    res0StandLeft,
		currResWalkingLeft:  res0WalkingLeft,
		currRes:             res0StandRight,

		grade:                 0,
		levelRect:             hitBox,
		renderBoxWExpandRatio: renderBoxWExpandRatio,
		renderBoxHExpandRatio: renderBoxHExpandRatio,
		renderBoxW:            res0StandLeft.GetW(),
		renderBoxH:            res0StandLeft.GetH(),
		velocity:              vector.Vec2D{0, 0},
		isOnGround:            false,
		isFacingRight:         true,
		lives:                 3,
	}
	return h
}

func (h *Hero) HandleEvents(events *intsets.Sparse, level *Level) {
	// standing on ground will absorb all X-velocity
	if h.isOnGround {
		h.velocity.X = 0
	}

	if events.Has(int(event.EVENT_KEYDOWN_LEFT)) {
		h.isFacingRight = false
		h.velocity.X = -350
	} else if events.Has(int(event.EVENT_KEYDOWN_RIGHT)) {
		h.isFacingRight = true
		h.velocity.X = 350
	}
	if events.Has(int(event.EVENT_KEYDOWN_SPACE)) {
		if h.isOnGround {
			h.velocity.Y = -1000
		}
	}
	if events.Has(int(event.EVENT_KEYDOWN_F)) {
		h.fPressed = true
	} else {
		h.fPressed = false
	}
	if events.Has(int(event.EVENT_KEYDOWN_UP)) {
		h.upPressed = true
	} else {
		h.upPressed = false
	}
	if events.Has(int(event.EVENT_KEYDOWN_F2)) {
		h.upgrade(level)
	}
	if events.Has(int(event.EVENT_KEYDOWN_F3)) {
		h.downgrade(level)
	}
}

func (h *Hero) Draw(g *graphic.Graphic, camPos vector.Pos) {
	// if hurt, blink for a while, otherwise just draw the hero
	ticks := sdl.GetTicks()
	if h.hurtStartTicks > 0 && ticks-h.hurtStartTicks < hurtAnimationMS {
		if (ticks-h.hurtStartTicks)%200 > 100 {
			g.DrawResource(h.currRes, h.getRenderRect(), camPos)
		} else {
			// Draw nothing to create a blink effect
		}
	} else {
		g.DrawResource(h.currRes, h.getRenderRect(), camPos)
	}

	// FIXME for debug only
	//g.DrawRect(h.getRenderRect(), camPos)
	//g.DrawRect(h.levelRect, camPos)
}

func (h *Hero) Update(ticks uint32, level *Level) {
	// skip first update due to lack of ticks
	if h.lastTicks == 0 {
		h.lastTicks = ticks
		return
	}

	// gravity: unit is pixels per second
	gravity := vector.Vec2D{0, 50}
	h.velocity.Add(gravity)

	maxVel := vector.Vec2D{int32(graphic.TILE_SIZE * 30 / 100), int32(graphic.TILE_SIZE * 30 / 100)}
	velocityStep := CalcVelocityStep(h.velocity, ticks, h.lastTicks, &maxVel)

	// apply velocity step
	h.levelRect.X += velocityStep.X
	h.levelRect.Y += velocityStep.Y

	// solve collision
	hitTop, hitRight, hitBottom, hitLeft, tilesHit := level.ObstMngr.SolveCollision(&h.levelRect)

	// update tiles hit
	h.notifyTilesHit(tilesHit, h.levelRect, velocityStep, level, ticks)

	// check if hit any live enemies
	for _, emy := range level.Enemies {
		if emy.IsDead() {
			continue
		}

		hit, direction := isHitEnemy(velocityStep, h.levelRect, emy.GetRect())
		if !hit {
			continue
		}

		emy.hitByHero(h, direction, level, ticks)
	}

	// is on ground
	h.isOnGround = hitBottom

	// reset velocity according to collision and direction
	if velocityStep.X > 0 && hitRight {
		h.velocity.X = 0
	}
	if velocityStep.X < 0 && hitLeft {
		h.velocity.X = 0
	}
	if velocityStep.Y > 0 && hitBottom {
		h.velocity.Y = 0
	}
	if velocityStep.Y < 0 && hitTop {
		h.velocity.Y = 0
	}

	// fire if needed and not too frequent
	if h.fPressed && ticks-h.lastFireTicks > 400 {
		level.AddVolatileObject(NewFireball(h.levelRect, h.isFacingRight, h.upPressed, ticks, level.ResourceRegistry))
		h.lastFireTicks = ticks
	}

	// check if need to reset hurt status
	if ticks-h.hurtStartTicks > hurtAnimationMS {
		h.hurtStartTicks = 0
	}

	// update resource
	h.updateRes()

	// update ticks
	h.lastTicks = ticks
}

func (h *Hero) GetRect() sdl.Rect {
	return h.levelRect
}

func (h *Hero) GetZIndex() int {
	return ZINDEX_4
}

func (h *Hero) Hurt() {
	// cannot hurt during super time
	if h.hurtStartTicks == 0 {
		h.lives--
		h.hurtStartTicks = sdl.GetTicks()
	}
}

func (h *Hero) GetLives() int {
	return h.lives
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private helpers
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *Hero) updateRes() {
	if h.velocity.X == 0 || h.lastTicks%600 < 300 {
		if h.isFacingRight {
			h.currRes = h.currResStandRight
		} else {
			h.currRes = h.currResStandLeft
		}
	} else if h.velocity.X != 0 && h.lastTicks%600 >= 300 {
		if h.isFacingRight {
			h.currRes = h.currResWalkingRight
		} else {
			h.currRes = h.currResWalkingLeft
		}
	}
}

func (h *Hero) notifyTilesHit(
	tilesHit []vector.TileID,
	resolvedRect sdl.Rect,
	heroVelStep vector.Vec2D,
	level *Level, ticks uint32) {

	for _, tid := range tilesHit {
		o := level.TileObjects[tid.X][tid.Y]
		if o == nil {
			log.Fatal("bug! notify hit tile object which is a nil object")
		}
		switch o.(type) {
		case hittableByHero:
			o.(hittableByHero).hitByHero(h, calcHitDirection(heroVelStep, resolvedRect, o.GetRect()), level, ticks)
		}
	}
}

// calcHitDirection decides from which direction was the tile being hit by hero
// NOTE:
// 1. It assumed that the hero and tile was intersected and then collision has been resolved
// 2. After resolution the two rects should have to be non-intersected
func calcHitDirection(heroVelStep vector.Vec2D, resolvedHeroRect sdl.Rect, tileRect sdl.Rect) hitDirection {
	if _, intersected := tileRect.Intersect(&resolvedHeroRect); intersected {
		log.Fatalf("calcHitDirection: hero %v and tile %v are intersected but should not", resolvedHeroRect, tileRect)
	}

	if resolvedHeroRect.Y >= tileRect.Y+tileRect.H && heroVelStep.Y < 0 {
		return HIT_FROM_BOTTOM_W_INTENT
	}

	if resolvedHeroRect.Y+resolvedHeroRect.H <= tileRect.Y && heroVelStep.Y > 0 {
		return HIT_FROM_TOP_W_INTENT
	}

	if resolvedHeroRect.X+resolvedHeroRect.W <= tileRect.X && heroVelStep.X > 0 {
		return HIT_FROM_LEFT_W_INTENT
	}

	if resolvedHeroRect.X >= tileRect.X+tileRect.W && heroVelStep.X < 0 {
		return HIT_FROM_RIGHT_W_INTENT
	}

	log.Println("calcHitDirection: no hit intention found")

	return HIT_WITH_NO_INTENTION
}

func isHitEnemy(heroVelStep vector.Vec2D, heroRect sdl.Rect, enemyRect sdl.Rect) (hit bool, hd hitDirection) {
	interRect, intersected := heroRect.Intersect(&enemyRect)
	if !intersected {
		return
	}

	hit = true

	// note that we also need to check velocity direction here,
	// because enemy is not an obstacle so is no constantly being collision resolved with hero
	// so a hero can move from an position where him already collides with the enemy
	if interRect.Y == enemyRect.Y && interRect.W > interRect.H && heroVelStep.Y > 0 {
		hd = HIT_FROM_TOP_W_INTENT
	} else if interRect.X == enemyRect.X && interRect.W < interRect.H && heroVelStep.X > 0 {
		hd = HIT_FROM_LEFT_W_INTENT
	} else if interRect.X+interRect.W == enemyRect.X+enemyRect.W && interRect.W < interRect.H && heroVelStep.X < 0 {
		hd = HIT_FROM_RIGHT_W_INTENT
	} else if interRect.Y+interRect.H == enemyRect.Y+enemyRect.H && interRect.W > interRect.H && heroVelStep.Y < 0 {
		hd = HIT_FROM_BOTTOM_W_INTENT
	} else {
		hd = HIT_WITH_NO_INTENTION
	}

	return
}

func (h *Hero) getRenderRect() sdl.Rect {
	resW := h.currResStandLeft.GetW()
	resH := h.currResStandLeft.GetH()
	return sdl.Rect{
		h.levelRect.X - int32(float64(resW)*h.renderBoxWExpandRatio/2),
		h.levelRect.Y - int32(float64(resH)*h.renderBoxHExpandRatio),
		resW,
		resH,
	}
}

func (h *Hero) upgrade(level *Level) {
	switch h.grade {
	case 0:
		h.grade = 1
		h.switchResSet(1)
	case 1:
		h.grade = 2
		h.switchResSet(2)
	}

	h.reCalcLevelRectSize()

	// show shine effects
	level.AddEffect(NewShineEffect(level.ResourceRegistry, h, sdl.GetTicks()))
}

func (h *Hero) downgrade(level *Level) {
	switch h.grade {
	case 1:
		h.grade = 0
		h.switchResSet(0)
	case 2:
		h.grade = 1
		h.switchResSet(1)
	}

	h.reCalcLevelRectSize()
}

func (h *Hero) switchResSet(grade int) {
	if grade < 0 || grade > 2 {
		log.Fatalf("cannot switch resource set: grade (%d) should be 0, 1 or 2", grade)
	}
	switch grade {
	case 0:
		h.currResStandRight = h.res0StandRight
		h.currResWalkingRight = h.res0WalkingRight
		h.currResStandLeft = h.res0StandLeft
		h.currResWalkingLeft = h.res0WalkingLeft
	case 1:
		h.currResStandRight = h.res1StandRight
		h.currResWalkingRight = h.res1WalkingRight
		h.currResStandLeft = h.res1StandLeft
		h.currResWalkingLeft = h.res1WalkingLeft
	case 2:
		h.currResStandRight = h.res2StandRight
		h.currResWalkingRight = h.res2WalkingRight
		h.currResStandLeft = h.res2StandLeft
		h.currResWalkingLeft = h.res2WalkingLeft
	}
}

// reCalcLevelRectSize reset the width and height of hit box to match current resource set
func (h *Hero) reCalcLevelRectSize() {
	preLevelRect := h.levelRect

	renderRect := h.getRenderRect()
	resX := renderRect.X
	resY := renderRect.Y
	resW := renderRect.W
	resH := renderRect.H
	h.levelRect = sdl.Rect{
		resX + int32(float64(resW)*h.renderBoxWExpandRatio/2),
		resY + int32(float64(resH)*h.renderBoxHExpandRatio),
		resW - int32(float64(resW)*h.renderBoxWExpandRatio),
		resH - int32(float64(resH)*h.renderBoxHExpandRatio),
	}

	// make sure the new level rect has same bottom position as the old one
	h.levelRect.Y += preLevelRect.Y + preLevelRect.H - (h.levelRect.Y + h.levelRect.H)
}
