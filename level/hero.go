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
	resStandRight   graphic.Resource
	resWalkingRight graphic.Resource
	resStandLeft    graphic.Resource
	resWalkingLeft  graphic.Resource

	// current resource
	currRes graphic.Resource

	// current rect in level, it is a hit box, not a render box
	// for hero its hit box and render box are different, hit box is smaller
	levelRect sdl.Rect

	// render box
	renderBoxXDelta int32
	renderBoxYDelta int32
	renderBoxW      int32
	renderBoxH      int32

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
	hitBoxXDelta, hitBoxYDelta, hitBoxWDelta, hitBoxHDelta int32,
	resourceRegistry map[graphic.ResourceID]graphic.Resource) *Hero {

	resStandRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_STAND_RIGHT]
	resWalkingRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_WALKING_RIGHT]
	resStandLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_STAND_LEFT]
	resWalkingLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_WALKING_LEFT]

	hitBox := sdl.Rect{
		renderBoxStartPos.X + hitBoxXDelta,
		renderBoxStartPos.Y + hitBoxYDelta,
		resStandLeft.GetW() + hitBoxWDelta,
		resStandLeft.GetH() + hitBoxHDelta,
	}

	h := &Hero{
		resStandRight:   resStandRight,
		resWalkingRight: resWalkingRight,
		resStandLeft:    resStandLeft,
		resWalkingLeft:  resWalkingLeft,
		currRes:         resStandRight,
		levelRect:       hitBox,
		renderBoxXDelta: -hitBoxXDelta,
		renderBoxYDelta: -hitBoxYDelta,
		renderBoxW:      resStandLeft.GetW(),
		renderBoxH:      resStandLeft.GetH(),
		velocity:        vector.Vec2D{0, 0},
		isOnGround:      false,
		isFacingRight:   true,
		lives:           3,
	}
	return h
}

func (h *Hero) HandleEvents(events *intsets.Sparse) {
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
			h.currRes = h.resStandRight
		} else {
			h.currRes = h.resStandLeft
		}
	} else if h.velocity.X != 0 && h.lastTicks%600 >= 300 {
		if h.isFacingRight {
			h.currRes = h.resWalkingRight
		} else {
			h.currRes = h.resWalkingLeft
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
	return sdl.Rect{
		h.levelRect.X + h.renderBoxXDelta,
		h.levelRect.Y + h.renderBoxYDelta,
		h.renderBoxW,
		h.renderBoxH,
	}
}
