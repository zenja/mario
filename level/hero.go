package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	mutils "github.com/zenja/mario/math_utils"
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
	drawResource(g, h.currRes, h.levelRect, camPos)
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

	// calculate velocity step
	velocityStep := h.velocity
	velocityStep.Multiply(int32(ticks - h.lastTicks))
	velocityStep.Divide(1000)

	// limit max velocity step
	maxVel := int32(graphic.TILE_SIZE * 30 / 100)
	if mutils.Abs(velocityStep.X) > maxVel {
		log.Printf("warning: velocity step's |X| is %d > %d, limited", velocityStep.X, maxVel)
		if velocityStep.X > 0 {
			velocityStep.X = maxVel
		} else {
			velocityStep.X = -maxVel
		}
	}
	if mutils.Abs(velocityStep.Y) > maxVel {
		log.Printf("warning: velocity step's |Y| is %d > %d, limited", velocityStep.Y, maxVel)
		if velocityStep.Y > 0 {
			velocityStep.Y = maxVel
		} else {
			velocityStep.Y = -maxVel
		}
	}

	// apply velocity step
	log.Printf("applying velocity step: %v\n", velocityStep)
	h.levelRect.X += velocityStep.X
	h.levelRect.Y += velocityStep.Y

	// solve collision
	log.Printf("desired rect: %v\n", h.levelRect)
	hitTop, hitRight, hitBottom, hitLeft := level.ObstMngr.SolveCollision(&h.levelRect)
	log.Printf("solved rect: %v\n", h.levelRect)

	// is on ground
	h.isOnGround = hitBottom

	log.Println("---")

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
