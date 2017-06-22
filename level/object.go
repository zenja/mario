package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// Note that Z-Index must start from 0 and increase one-by-one
const ZINDEX_NUM = 5
const (
	ZINDEX_0 = iota
	ZINDEX_1
	ZINDEX_2
	ZINDEX_3
	ZINDEX_4
)

type Object interface {
	Draw(camPos vector.Pos)
	Update(ticks uint32, level *Level)
	// object hit box (not render box)
	GetRect() sdl.Rect
	GetZIndex() int
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Important interfaces
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type hitDirection int

const (
	HIT_FROM_TOP_W_INTENT hitDirection = iota
	HIT_FROM_RIGHT_W_INTENT
	HIT_FROM_BOTTOM_W_INTENT
	HIT_FROM_LEFT_W_INTENT
	HIT_WITH_NO_INTENTION
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Single tile object
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type singleTileObject struct {
	resource graphic.Resource
	// the position rect on the level
	levelRect sdl.Rect
	zIndex    int
}

func NewSingleTileObject(resource graphic.Resource, startPos vector.Pos, zIndex int) Object {
	return &singleTileObject{
		resource:  resource,
		levelRect: sdl.Rect{startPos.X, startPos.Y, resource.GetW(), resource.GetH()},
		zIndex:    zIndex,
	}
}

func (sto *singleTileObject) Draw(camPos vector.Pos) {
	graphic.DrawResource(sto.resource, sto.levelRect, camPos)
}

func (sto *singleTileObject) Update(ticks uint32, level *Level) {
	// Do nothing
}

// object hit box (not render box)
func (sto *singleTileObject) GetRect() sdl.Rect {
	return sto.levelRect
}

func (sto *singleTileObject) GetZIndex() int {
	return sto.zIndex
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Invisible tile object
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type invisibleTileObject struct {
	tid vector.TileID
}

func NewInvisibleTileObject(tid vector.TileID) Object {
	return &invisibleTileObject{tid}
}

func (ito *invisibleTileObject) Draw(camPos vector.Pos) {
	// Do nothing
}

func (ito *invisibleTileObject) Update(ticks uint32, level *Level) {
	// Do nothing
}

func (ito *invisibleTileObject) GetRect() sdl.Rect {
	return GetTileRect(ito.tid)
}

func (ito *invisibleTileObject) GetZIndex() int {
	return ZINDEX_0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
