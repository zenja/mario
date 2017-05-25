package object

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
)

type hero struct {
	resource graphic.Resource
	levelPos *sdl.Rect
}

func NewHero(xStart, yStart int32, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
	resource, ok := resourceRegistry[graphic.RESOURCE_TYPE_HERO]
	if !ok {
		log.Fatalf("resource not fount in resource registry: %d", graphic.RESOURCE_TYPE_HERO)
	}
	return &hero{
		resource: resource,
		levelPos: &sdl.Rect{xStart, yStart, resource.GetW(), resource.GetH()},
	}
}

func (h *hero) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	drawResource(g, h.resource, h.levelPos, xCamStart, yCamStart)
}
