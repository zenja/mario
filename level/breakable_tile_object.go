package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/audio"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
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

func (bto *breakableTileObject) Draw(camPos vector.Pos) {
	graphic.DrawResource(bto.mainRes, bto.levelRect, camPos)
}

func (bto *breakableTileObject) Update(ticks uint32, level *Level) {
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
	hitEnemiesOnTop(&bto.levelRect, level, ticks)

	// remove object and obstacle
	tid := GetTileID(vector.Pos{bto.levelRect.X, bto.levelRect.Y}, false, false)
	level.RemoveObstacleTileObject(tid)

	// show breaking effect
	level.AddEffect(NewBreakTileEffect(bto.pieceRes, tid, ticks))

	// play sound
	audio.PlaySound(audio.SOUND_BREAK_BRICK)
}
