package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Object = &animationTileObj{}

type animationTileObj struct {
	reses      []graphic.Resource
	currResIdx int
	levelRect  sdl.Rect
	frameMs    uint32 // how long each res shows
	zIndex     int
}

func NewAnimationObject(startPos vector.Vec2D, resIDs []graphic.ResourceID, frameMs uint32, zIndex int) *animationTileObj {
	var reses []graphic.Resource
	for _, id := range resIDs {
		reses = append(reses, graphic.Res(id))
	}

	return &animationTileObj{
		reses:     reses,
		levelRect: sdl.Rect{startPos.X, startPos.Y, reses[0].GetW(), reses[0].GetH()},
		frameMs:   frameMs,
		zIndex:    zIndex,
	}
}

func NewAnimationObjectTID(tid vector.TileID, resIDs []graphic.ResourceID, frameMs uint32, zIndex int) *animationTileObj {
	return NewAnimationObject(GetTileStartPos(tid), resIDs, frameMs, zIndex)
}

func NewWaterSurfaceAnimationObject(tid vector.TileID) *animationTileObj {
	resIDs := []graphic.ResourceID{
		graphic.RESOURCE_TYPE_WATER_0,
		graphic.RESOURCE_TYPE_WATER_1,
		graphic.RESOURCE_TYPE_WATER_2,
		graphic.RESOURCE_TYPE_WATER_3,
		graphic.RESOURCE_TYPE_WATER_4,
		graphic.RESOURCE_TYPE_WATER_5,
		graphic.RESOURCE_TYPE_WATER_6,
	}
	return NewAnimationObject(GetTileStartPos(tid), resIDs, 250, ZINDEX_1)
}

func (ato *animationTileObj) Draw(camPos vector.Pos) {
	graphic.DrawResource(ato.reses[ato.currResIdx], ato.GetRect(), camPos)
}

func (ato *animationTileObj) Update(ticks uint32, level *Level) {
	ato.currResIdx = int((ticks / ato.frameMs) % uint32(len(ato.reses)))
}

func (ato *animationTileObj) GetRect() sdl.Rect {
	return ato.levelRect
}

func (ato *animationTileObj) GetZIndex() int {
	return ato.zIndex
}
