package object

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type HeroState int

const (
	HERO_STATE_STAND HeroState = iota
	HERO_STATE_WALKING
)

type hero struct {
	resStand   graphic.Resource
	resWalking graphic.Resource

	// current state
	currState HeroState

	// current resource
	currRes graphic.Resource

	// current rect in level
	currLevelRect *sdl.Rect

	velocity vector.Vec2D
}

func NewHero(startPos vector.Pos, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
	resStand, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_STAND]
	resWalking, _ := resourceRegistry[graphic.RESOURCE_TYPE_HERO_WALKING]
	return &hero{
		resStand:   resStand,
		resWalking: resWalking,
		// init stat is standing
		currState:     HERO_STATE_STAND,
		currRes:       resStand,
		currLevelRect: &sdl.Rect{startPos.X, startPos.Y, resStand.GetW(), resStand.GetH()},
		// init velocity is zero
		velocity: vector.Vec2D{},
	}
}

func (h *hero) Draw(g *graphic.Graphic, camPos vector.Pos) {
	drawResource(g, h.currRes, h.currLevelRect, camPos)
}

func (h *hero) Update(events *intsets.Sparse, ticks uint32) {
	// handle movement
	switch {
	case events.Has(int(event.EVENT_KEYDOWN_LEFT)):
		h.currState = HERO_STATE_WALKING
		h.velocity.X = -10
	case events.Has(int(event.EVENT_KEYDOWN_RIGHT)):
		h.currState = HERO_STATE_WALKING
		h.velocity.X = 10
	case events.Has(int(event.EVENT_KEYDOWN_SPACE)):
		h.currLevelRect.Y -= 1
	default:
		// FIXME
		//h.currState = HERO_STATE_STAND
		h.velocity.X = 0
		h.velocity.Y = 0
	}

	// apply velocity
	h.currLevelRect.X += h.velocity.X
	h.currLevelRect.Y += h.velocity.Y

	if h.currState == HERO_STATE_WALKING {
		if ticks%400 < 200 {
			h.currRes = h.resStand
		} else {
			h.currRes = h.resWalking
		}
	}

	if h.currState == HERO_STATE_STAND {
		h.currRes = h.resStand
	}

}

func (h *hero) GetRect() sdl.Rect {
	return *h.currLevelRect
}

func (h *hero) GetZIndex() int {
	return ZINDEX_4
}
