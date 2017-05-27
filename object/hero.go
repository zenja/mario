package object

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"golang.org/x/tools/container/intsets"
	"github.com/zenja/mario/vector"
)

type hero struct {
	resource  graphic.Resource
	levelRect *sdl.Rect
}

func NewHero(startPos vector.Pos, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
	resource, ok := resourceRegistry[graphic.RESOURCE_TYPE_HERO]
	if !ok {
		log.Fatalf("resource not found in resource registry: %d", graphic.RESOURCE_TYPE_HERO)
	}
	return &hero{
		resource:  resource,
		levelRect: &sdl.Rect{startPos.X, startPos.Y, resource.GetW(), resource.GetH()},
	}
}

func (h *hero) Draw(g *graphic.Graphic, camPos vector.Pos) {
	drawResource(g, h.resource, h.levelRect, camPos)
}

func (h *hero) Update(events *intsets.Sparse, ticks uint32) {
	// handle movement
	switch {
	case events.Has(event.EVENT_KEYDOWN_LEFT):
		h.levelRect.X -= 1
	case events.Has(event.EVENT_KEYDOWN_RIGHT):
		h.levelRect.X += 1
	case events.Has(event.EVENT_KEYDOWN_SPACE):
		h.levelRect.Y -= 1
	}
}

func (h *hero) GetRect() sdl.Rect {
	return *h.levelRect
}

func (h *hero) GetZIndex() int {
	return ZINDEX_4
}
