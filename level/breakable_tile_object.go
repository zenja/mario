package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

// assert that breakableTileObject is hit-able by hero
var _ hittableByHero = &breakableTileObject{}

type breakableTileObject struct {
	mainRes   graphic.Resource
	pieceRes  graphic.Resource
	levelRect sdl.Rect
	zIndex    int
}

func NewBreakableTileObject(mainRes graphic.Resource, pieceRes graphic.Resource, startPos vector.Pos, zIndex int) Object {
	return &breakableTileObject{
		mainRes:   mainRes,
		pieceRes:  pieceRes,
		levelRect: sdl.Rect{startPos.X, startPos.Y, graphic.TILE_SIZE, graphic.TILE_SIZE},
	}
}

func (bto *breakableTileObject) Draw(g *graphic.Graphic, camPos vector.Pos) {
	g.DrawResource(bto.mainRes, bto.levelRect, camPos)
}

func (bto *breakableTileObject) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	// Do nothing
}

func (bto *breakableTileObject) GetRect() sdl.Rect {
	return bto.levelRect
}

func (bto *breakableTileObject) GetZIndex() int {
	return bto.zIndex
}

func (bto *breakableTileObject) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	// can only be hit from bottom
	// TODO in the future the direction may need to be decided from input
	if direction != HIT_FROM_BOTTOM_W_INTENT {
		return
	}

	// check if any enemy stand on this tile, hit them
	for _, e := range level.Enemies {
		emyRectLower := sdl.Rect{
			X: e.GetRect().X,
			Y: e.GetRect().Y + 1,
			W: e.GetRect().W,
			H: e.GetRect().H,
		}
		if emyRectLower.HasIntersection(&bto.levelRect) {
			if !e.IsDead() {
				e.hitByBrokenTile(level, ticks)
			}
		}
	}

	// remove object and obstacle
	tid := GetTileID(vector.Pos{bto.levelRect.X, bto.levelRect.Y}, false, false)
	level.RemoveObstacleTileObject(tid)

	// show breaking effect
	level.AddEffect(NewBreakTileEffect(bto.pieceRes, tid, ticks))
}
