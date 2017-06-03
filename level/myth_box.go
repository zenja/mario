package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type mythBox struct {
	// resources
	resNormal graphic.Resource
	resEmpty  graphic.Resource // empty, no coins
	currRes   graphic.Resource

	// myth box has both a tile rect and current level rect,
	// because we allow myth box to move a little bit after being hit
	tileRect  sdl.Rect
	levelRect sdl.Rect

	numCoinsLeft int
}

func NewMythBox(startPos vector.Pos, numCoins int, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
	resNormal, _ := resourceRegistry[graphic.RESOURCE_TYPE_MYTH_BOX_NORMAL]
	resEmpty, _ := resourceRegistry[graphic.RESOURCE_TYPE_MYTH_BOX_EMPTY]
	tileRect := sdl.Rect{startPos.X, startPos.Y, graphic.TILE_SIZE, graphic.TILE_SIZE}
	return &mythBox{
		resNormal:    resNormal,
		resEmpty:     resEmpty,
		currRes:      resNormal,
		tileRect:     tileRect,
		levelRect:    tileRect,
		numCoinsLeft: numCoins,
	}
}

func (mb *mythBox) Draw(g *graphic.Graphic, camPos vector.Pos) {
	drawResource(g, mb.currRes, &mb.levelRect, camPos)
}

func (mb *mythBox) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	// TODO
}

func (mb *mythBox) GetRect() sdl.Rect {
	return mb.levelRect
}

func (mb *mythBox) GetZIndex() int {
	return ZINDEX_1
}
