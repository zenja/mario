package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type hero struct {
	resStandRight   graphic.Resource
	resWalkingRight graphic.Resource
	resStandLeft    graphic.Resource
	resWalkingLeft  graphic.Resource

	// current resource
	currRes graphic.Resource

	// current rect in level
	levelRect sdl.Rect

	// current velocity, unit is pixels per second
	velocity vector.Vec2D

	lastTicks uint32

	isOnGround bool

	isFacingRight bool
}

func NewHero(startPos vector.Pos, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
	resStandRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_STAND_RIGHT]
	resWalkingRight, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_WALKING_RIGHT]
	resStandLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_STAND_LEFT]
	resWalkingLeft, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_WALKING_LEFT]
	h := &hero{
		resStandRight:   resStandRight,
		resWalkingRight: resWalkingRight,
		resStandLeft:    resStandLeft,
		resWalkingLeft:  resWalkingLeft,
		currRes:         resStandRight,
		levelRect:       sdl.Rect{startPos.X, startPos.Y, resStandLeft.GetW(), resStandLeft.GetH()},
		velocity:        vector.Vec2D{0, 0},
		isOnGround:      false,
		isFacingRight:   true,
	}
	return h
}

func (h *hero) Draw(g *graphic.Graphic, camPos vector.Pos) {
	g.DrawResource(h.currRes, h.levelRect, camPos)
}

func (h *hero) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	// skip first update due to lack of ticks
	if h.lastTicks == 0 {
		h.lastTicks = ticks
		return
	}

	// standing on ground will absorb all X-velocity
	if h.isOnGround {
		h.velocity.X = 0
	}

	// ---------------------------------------
	// Handle events
	// ---------------------------------------
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

	// gravity: unit is pixels per second
	gravity := vector.Vec2D{0, 50}
	h.velocity.Add(gravity)

	maxVel := vector.Vec2D{int32(graphic.TILE_SIZE * 30 / 100), int32(graphic.TILE_SIZE * 30 / 100)}
	velocityStep := CalcVelocityStep(h.velocity, ticks, h.lastTicks, &maxVel)

	// apply velocity step
	//log.Printf("applying velocity step: %v\n", velocityStep)
	h.levelRect.X += velocityStep.X
	h.levelRect.Y += velocityStep.Y

	// solve collision
	//log.Printf("desired rect: %v\n", h.levelRect)
	hitTop, hitRight, hitBottom, hitLeft, tilesHit := level.ObstMngr.SolveCollision(&h.levelRect)
	//log.Printf("solved rect: %v\n", h.levelRect)

	// update tiles hit
	h.notifyTilesHit(tilesHit, h.levelRect, level, ticks)

	// check if hit any live enemies
	for _, emy := range level.Enemies {
		if emy.IsDead() {
			continue
		}

		hit, hitEmyTop := isHitEnemy(h.levelRect, emy.GetLevelRect())
		if !hit {
			continue
		}

		if hitEmyTop {
			emy.hitByHero(h, HIT_FROM_TOP, level, ticks)
		} else {
			emy.hitByHero(h, HIT_NOT_FROM_TOP, level, ticks)
		}
	}

	// is on ground
	h.isOnGround = hitBottom

	//log.Println("---")

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

	// update resource
	h.updateRes()

	// update ticks
	h.lastTicks = ticks
}

func (h *hero) GetRect() sdl.Rect {
	return h.levelRect
}

func (h *hero) GetZIndex() int {
	return ZINDEX_4
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private helpers
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *hero) updateRes() {
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

func (h *hero) notifyTilesHit(tilesHit []vector.TileID, resolvedRect sdl.Rect, level *Level, ticks uint32) {
	for _, tid := range tilesHit {
		o := level.TileObjects[tid.X][tid.Y]
		if o == nil {
			log.Fatal("bug! notify hit tile object which is a nil object")
		}
		switch o.(type) {
		case heroHittableObject:
			o.(heroHittableObject).hitByHero(h, calcHitDirection(resolvedRect, o.GetRect()), level, ticks)
		}
	}
}

// calcHitDirection decides from which direction was the tile being hit by hero
// It assumed that the hero and tile was intersected and then collision has been resolved
// note that after resolution the two rects should have to be non-intersected
func calcHitDirection(resolvedHeroRect sdl.Rect, tileRect sdl.Rect) hitDirection {
	if _, intersected := tileRect.Intersect(&resolvedHeroRect); intersected {
		log.Fatalf("calcHitDirection: hero %v and tile %v are intersected but should not", resolvedHeroRect, tileRect)
	}

	if resolvedHeroRect.Y >= tileRect.Y+tileRect.H {
		return HIT_FROM_BOTTOM
	}

	if resolvedHeroRect.Y+resolvedHeroRect.H <= tileRect.Y {
		return HIT_FROM_TOP
	}

	if resolvedHeroRect.X+resolvedHeroRect.W <= tileRect.X {
		return HIT_FROM_LEFT
	}

	if resolvedHeroRect.X >= tileRect.X+tileRect.W {
		return HIT_FROM_RIGHT
	}

	log.Fatal("bug! should already covered all possible cases")

	return HIT_FROM_TOP
}

func isHitEnemy(heroRect sdl.Rect, enemyRect sdl.Rect) (hit bool, hitEnemyTop bool) {
	interRect, intersected := heroRect.Intersect(&enemyRect)
	if !intersected {
		return
	}

	hit = true

	if interRect.Y == enemyRect.Y && interRect.W > interRect.H {
		hitEnemyTop = true
	}

	return
}
