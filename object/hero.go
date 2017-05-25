package object

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
)

type hero struct {
	tileUp   graphic.TileID
	tileDown graphic.TileID

	upLevelPos   *sdl.Rect
	downLevelPos *sdl.Rect
}

func NewHero(xStart, yStart int32) Object {
	return &hero{
		tileUp:       graphic.TILE_TYPE_HERO,
		tileDown:     graphic.TILE_TYPE_HERO,
		upLevelPos:   &sdl.Rect{xStart, yStart, graphic.TILE_SIZE, graphic.TILE_SIZE},
		downLevelPos: &sdl.Rect{xStart, yStart + graphic.TILE_SIZE, graphic.TILE_SIZE, graphic.TILE_SIZE},
	}
}

func (h *hero) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	drawTile(g, h.tileUp, h.upLevelPos, xCamStart, yCamStart)
	drawTile(g, h.tileDown, h.downLevelPos, xCamStart, yCamStart)
}
